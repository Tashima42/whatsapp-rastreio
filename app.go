package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/tashima42/shared-expenses-manager-backend/helpers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/tashima42/shared-expenses-manager-backend/handlers"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
	Cache  *cache.Cache
	Logger *helpers.Logger
}

var apikey string

func (a *App) Initialize(
	user string,
	password string,
	dbname string,
	whatsappPhoneNumberId string,
	whatsappUserAccessToken string,
	whatsappBaseUrl string,
	secret string,
	env string,
) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode='disable'", user, password, dbname)

	apikey = secret

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.ensureTableExists()
	a.Cache = cache.New(24*time.Hour, 48*time.Hour)
	a.Logger = &helpers.Logger{DB: a.DB, Env: env}

	whatsappProvider := helpers.WhatsappProvider{
		PhoneNumberId:   whatsappPhoneNumberId,
		UserAccessToken: whatsappUserAccessToken,
		BaseUrl:         whatsappBaseUrl,
	}

	bucketHandler := handlers.BucketHandler{DB: a.DB}
	userHandler := handlers.UserHandler{DB: a.DB}
	whatsappHandler := handlers.WhatsappHandler{
		DB:               a.DB,
		WhatsappProvider: whatsappProvider,
		Cache:            a.Cache,
		Logger:           a.Logger,
	}

	a.Router = mux.NewRouter()
	a.Router.Use(a.loggingMiddleware)
	a.Router.HandleFunc("/whatsapp/webhook", whatsappHandler.WebhookVerify).Methods(http.MethodGet)
	authRouter := a.Router.PathPrefix("/").Subrouter()
	authRouter.Use(a.loggingMiddleware)
	authRouter.Use(authorizeMiddleware)
	authRouter.HandleFunc("/bucket", bucketHandler.CreateBucket).Methods(http.MethodPost)
	authRouter.HandleFunc("/whatsapp/webhook", whatsappHandler.Webhook).Methods(http.MethodPost)
	authRouter.HandleFunc("/whatsapp/message", whatsappHandler.SendMessage).Methods(http.MethodPost)
	authRouter.HandleFunc("/user", userHandler.CreateUser).Methods(http.MethodPost)

}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		body := r.Body
		b, _ := io.ReadAll(body)
		originalBody := io.NopCloser(bytes.NewBuffer(b))
		type logRequest struct {
			Uri    string `json:"uri"`
			Body   string `json:"body"`
			Method string `json:"method"`
		}

		logReq := &logRequest{Uri: r.RequestURI, Body: string(b), Method: r.Method}
		logJson, err := json.Marshal(logReq)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(logJson))

		a.Logger.Log(string(logJson), "request")

		r.Body = originalBody
		next.ServeHTTP(w, r)
	})
}
func authorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		fmt.Println("test")
		fmt.Println(apikey, authorizationHeader)
		if authorizationHeader != apikey {
			helpers.RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED-INVALID-APIKEY", "Invalid apikey")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) ensureTableExists() {
	b, _ := ioutil.ReadFile("./schema.sql")
	tableCreationQuery := string(b)
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

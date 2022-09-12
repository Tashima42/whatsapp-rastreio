package main

import (
	"database/sql"
	"fmt"
	"github.com/patrickmn/go-cache"
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
) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode='disable'", user, password, dbname)

	apikey = secret

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Cache = cache.New(24*time.Hour, 48*time.Hour)

	whatsappProvider := helpers.WhatsappProvider{
		PhoneNumberId:   whatsappPhoneNumberId,
		UserAccessToken: whatsappUserAccessToken,
		BaseUrl:         whatsappBaseUrl,
	}

	bucketHandler := handlers.BucketHandler{DB: a.DB}
	whatsappHandler := handlers.WhatsappHandler{
		DB:               a.DB,
		WhatsappProvider: whatsappProvider,
		Cache:            a.Cache,
	}

	a.Router = mux.NewRouter()
	a.Router.Use(loggingMiddleware)
	a.Router.HandleFunc("/whatsapp/webhook", whatsappHandler.WebhookVerify).Methods(http.MethodGet)
	a.Router.HandleFunc("/bucket", bucketHandler.CreateBucket).Methods(http.MethodPost)
	a.Router.HandleFunc("/whatsapp/webhook", whatsappHandler.Webhook).Methods(http.MethodPost)
	a.Router.HandleFunc("/whatsapp/message", whatsappHandler.SendMessage).Methods(http.MethodPost)
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
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

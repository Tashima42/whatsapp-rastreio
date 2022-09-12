package main

import (
	"database/sql"
	"fmt"
	"github.com/tashima42/shared-expenses-manager-backend/helpers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/tashima42/shared-expenses-manager-backend/handlers"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
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
	whatsappProvider := helpers.WhatsappProvider{
		PhoneNumberId:   whatsappPhoneNumberId,
		UserAccessToken: whatsappUserAccessToken,
		BaseUrl:         whatsappBaseUrl,
	}

	bucketHandler := handlers.BucketHandler{DB: a.DB}
	whatsappHandler := handlers.WhatsappHandler{DB: a.DB, WhatsappProvider: whatsappProvider}

	a.Router = mux.NewRouter()
	a.Router.Use(loggingMiddleware)
	a.Router.Use(authorizeMiddleware)
	a.Router.HandleFunc("/bucket", bucketHandler.CreateBucket)
	a.Router.HandleFunc("/whatsapp/webhook", whatsappHandler.Webhook)
	a.Router.HandleFunc("/whatsapp/message", whatsappHandler.SendMessage)
	a.Router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		helpers.RespondWithJSON(w, http.StatusOK, "ok")
	})
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

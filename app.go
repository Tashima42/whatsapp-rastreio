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

func (a *App) Initialize(user string, password string, dbname string, whatsappPhoneNumberId string, whatsappUserAccessToken string, whatsappBaseUrl string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode='disable'", user, password, dbname)

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
	a.Router.HandleFunc("/bucket", bucketHandler.CreateBucket)
	a.Router.HandleFunc("/whatsapp/webhook", whatsappHandler.Webhook)
	a.Router.HandleFunc("/whatsapp/message", whatsappHandler.SendMessage)
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

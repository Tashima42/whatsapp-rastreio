package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/tashima42/shared-expenses-manager-backend/helpers"
	"io"
	"net/http"
	"os"
)

type WhatsappHandler struct {
	DB               *sql.DB
	WhatsappProvider helpers.WhatsappProvider
	Cache            *cache.Cache
}

func (wh *WhatsappHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	type requestDTO struct {
		To string `json:"to"`
	}
	var request requestDTO
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Print(err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, "WHATSAPP-SEND-MESSAGE-INVALID-BODY", "Unable to parse request body")
		return
	}

	b, _ := wh.WhatsappProvider.SendMessage(request.To)

	type responseDTO struct {
		Success bool   `json:"success"`
		body    string `json:"body"`
	}

	b2, _ := io.ReadAll(b.Body)
	fmt.Println(string(b2))

	helpers.RespondWithJSON(w, http.StatusOK, responseDTO{Success: true, body: string(b2)})
}

func (wh *WhatsappHandler) WebhookVerify(w http.ResponseWriter, r *http.Request) {
	hubMode := r.URL.Query().Get("hub.mode")
	hubVerifyToken := r.URL.Query().Get("hub.verify_token")
	hubChallenge := r.URL.Query().Get("hub.challenge")

	if hubMode == "subscribe" && hubVerifyToken == os.Getenv("SECRET") {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(hubChallenge))
		return
	}

	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("Forbidden"))
}

type responseDTO struct {
	Success bool `json:"success"`
}

func (wh *WhatsappHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var message helpers.WhatsAppReceivedMessageObject
	err := decoder.Decode(&message)
	if err != nil {
		fmt.Print(err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, "WHATSAPP-WEBHOOK-INVALID-BODY", "Unable to parse request body")
		return
	}

	messages := message.Entry[0].Changes[0].Value.Messages
	if messages == nil {
		helpers.RespondWithJSON(w, http.StatusOK, responseDTO{Success: true})
		return
	}

	wasProcessed, _ := wh.Cache.Get(messages[0].Id)

	if wasProcessed == true {
		helpers.RespondWithJSON(w, http.StatusOK, responseDTO{Success: true})
		return
	}

	wh.Cache.SetDefault(messages[0].Id, true)

	messageId := messages[0].Id
	buttonPayload := messages[0].Button.Payload
	fromPhoneNumber := messages[0].From
	wh.WhatsappProvider.AckMessage(messageId)

	if buttonPayload == "Informações de pagamento" {
		b, _ := wh.WhatsappProvider.ReplyWithPix(messageId, fromPhoneNumber)
		b2, _ := io.ReadAll(b.Body)
		fmt.Println(string(b2))
	}

	helpers.RespondWithJSON(w, http.StatusOK, responseDTO{Success: true})
}

/*
func (wh *WhatsappHandler) UploadImage(w http.ResponseWriter, r *http.Request) {

	whatsappProvider.UploadImage()
}
*/

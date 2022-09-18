package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/tashima42/whatsapp-rastreio/data"
	"github.com/tashima42/whatsapp-rastreio/helpers"
	whatsappTemplates "github.com/tashima42/whatsapp-rastreio/helpers/whatsapp-templates"
	"github.com/tashima42/whatsapp-rastreio/providers"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type WhatsappHandler struct {
	DB               *sql.DB
	WhatsappProvider providers.WhatsappProvider
	Cache            *cache.Cache
	Logger           *helpers.Logger
	CorreiosProvider *providers.CorreiosProvider
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

	err = wh.registerPackage(message)
	if err != nil {
		fmt.Println(err.Error())
		invalidCodeTemplate := whatsappTemplates.CorreiosInvalidCode{}
		wh.WhatsappProvider.SendMessageTemplate(invalidCodeTemplate.GetTemplate(), messages[0].From)
	}

	messageId := messages[0].Id
	wh.WhatsappProvider.AckMessage(messageId)

	helpers.RespondWithJSON(w, http.StatusOK, responseDTO{Success: true})
}

func (wh *WhatsappHandler) registerPackage(whatsappMessage helpers.WhatsAppReceivedMessageObject) error {
	if whatsappMessage.Object != "whatsapp_business_account" {
		return fmt.Errorf("invalid object type")
	}
	if len(whatsappMessage.Entry) <= 0 {
		return fmt.Errorf("entry must have at least one member")
	}
	if len(whatsappMessage.Entry[0].Changes) <= 0 {
		return fmt.Errorf("changes must have at least one member")
	}
	changes := whatsappMessage.Entry[0].Changes[0]
	if changes.Field != "messages" {
		return fmt.Errorf("field must be messages")
	}
	if changes.Value.MessagingProduct != "whatsapp" {
		return fmt.Errorf("messaging Product must be whatsapp")
	}
	if len(changes.Value.Messages) <= 0 {
		return fmt.Errorf("messages must have at least one member")
	}
	messages := changes.Value.Messages[0]
	if messages.Type != "text" {
		return fmt.Errorf("type must be text")
	}
	codeRegex, err := regexp.Compile("^[A-Z]{2}[0-9]{9}[A-Z]{2}$")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return fmt.Errorf("internal error failed to compile regex")
	}
	textBody := messages.Text.Body
	splitBody := strings.Split(textBody, " ")
	code := splitBody[0]
	if codeRegex.MatchString(code) != true {
		return fmt.Errorf("invalid code format")
	}
	splitBody = splitBody[1:]
	name := strings.Join(splitBody, " ")
	user := data.User{Number: messages.From}
	err = user.GetByNumber(wh.DB)
	if err != nil {
		fmt.Printf("Error: %s", err)
		if err.Error() == "sql: no rows in result set" {
			err = user.CreateUser(wh.DB)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("failed to get user")
		}
	}
	valid := wh.validateObjectCode(code)
	if valid != true {
		return fmt.Errorf("invalid code")
	}
	object := data.Object{
		Code: code,
		Name: name,
	}
	err = object.GetByCode(wh.DB)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		if err.Error() == "sql: no rows in result set" {
			wh.CorreiosProvider.RegisterPackage(user, object)
		} else {
			err = user.AddObject(wh.DB, object)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (wh *WhatsappHandler) validateObjectCode(code string) bool {
	codes := []string{code}
	_, err := wh.CorreiosProvider.GetCorreiosObjects(codes)
	if err != nil {
		return false
	}
	return true
}

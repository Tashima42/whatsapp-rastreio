package routines

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/tashima42/whatsapp-rastreio/data"
	whatsappTemplates "github.com/tashima42/whatsapp-rastreio/helpers/whatsapp-templates"
	"github.com/tashima42/whatsapp-rastreio/providers"
	"net/http"
)

type WhatsappRoutines struct {
	DB               *sql.DB
	WhatsappProvider providers.WhatsappProvider
}

func (wr *WhatsappRoutines) SendUserUpdates() {
	users, err := data.GetUsers(wr.DB)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}
	for _, user := range users {
		go wr.sendUserObjectsUpdates(user)
	}
}

func (wr *WhatsappRoutines) sendUserObjectsUpdates(user data.User) {
	objects, err := user.GetUserObjects(wr.DB)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	for _, object := range objects {
		event, err := object.GetObjectEvent(wr.DB)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		if event.Hash != object.LastSentHash {
			correiosObjectEvent := providers.CorreiosObjectEvent{}
			decoder := json.NewDecoder(bytes.NewReader([]byte(event.Body)))
			err = decoder.Decode(&correiosObjectEvent)
			header := object.Code
			if object.Name != "" {
				header = fmt.Sprintf("%s - %s", object.Name, object.Code)
			}
			correiosUpdateTemplate := whatsappTemplates.CorreiosEventUpdate{
				Header:      header,
				Description: correiosObjectEvent.Description,
				Date:        fmt.Sprintf("%s - %s", correiosObjectEvent.Date, correiosObjectEvent.Hour),
				Local:       correiosObjectEvent.Unity.Local,
			}
			res, err := wr.WhatsappProvider.SendMessageTemplate(correiosUpdateTemplate.GetTemplate(), user.Number)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return
			}
			if res.StatusCode == http.StatusOK {
				object.LastSentHash = event.Hash
				err = object.UpdateLastSentHash(wr.DB)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
			}
		}
	}
}

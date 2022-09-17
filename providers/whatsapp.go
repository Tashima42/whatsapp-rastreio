package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WhatsappProvider struct {
	PhoneNumberId   string
	UserAccessToken string
	BaseUrl         string
}

func (wp *WhatsappProvider) SendMessageTemplate(template any, to string) (*http.Response, error) {
	tpl := struct {
		MessagingProduct string `json:"messaging_product"`
		To               string `json:"to"`
		Type             string `json:"type"`
		Template         any    `json:"template"`
	}{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template:         template,
	}

	fmt.Println("sending")
	body, err := json.Marshal(tpl)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	return wp.whatsappRequest(http.MethodPost, "messages", "application/json", bytes.NewReader(body))
}

func (wp *WhatsappProvider) AckMessage(messageId string) {
	type ackMessageBody struct {
		MessagingProduct string `json:"messaging_product"`
		Status           string `json:"status"`
		MessageId        string `json:"message_id"`
	}
	body, err := json.Marshal(ackMessageBody{MessagingProduct: "whatsapp", Status: "read", MessageId: messageId})
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	wp.whatsappRequest(http.MethodPut, "messages", "application/json", bytes.NewReader(body))
}

func (wp *WhatsappProvider) whatsappRequest(method string, endpoint string, contentType string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/%s", wp.BaseUrl, wp.PhoneNumberId, endpoint)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", wp.UserAccessToken))
	client := &http.Client{}
	return client.Do(req)
}

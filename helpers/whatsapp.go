package helpers

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

func (wp *WhatsappProvider) ReplyWithPix(messageId string, toNumber string) (*http.Response, error) {
	type text struct {
		PreviewUrl bool   `json:"preview_url"`
		Body       string `json:"body"`
	}
	type image struct {
		Id      string `json:"id"`
		Caption string `json:"caption"`
	}
	type context struct {
		MessageId string `json:"message_id"`
	}
	type message struct {
		MessagingProduct string  `json:"messaging_product"`
		RecipientType    string  `json:"recipient_type"`
		To               string  `json:"to"`
		Type             string  `json:"type"`
		Context          context `json:"context"`
		Image            image   `json:"image"`
	}

	sendMessage := message{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               toNumber,
		Type:             "image",
		Image: image{
			Id:      "2913711642258741",
			Caption: "00020126770014BR.GOV.BCB.PIX01362b071dc0-461c-4698-be90-486be9a352b70215Youtube Premium52040000530398654045.305802BR5913Pedro Tashima6008Londrina62410503***50300017BR.GOV.BCB.BRCODE01051.0.063046F88",
		},
		Context: context{MessageId: messageId},
	}

	body, err := json.Marshal(sendMessage)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	return wp.whatsappRequest(http.MethodPost, "messages", "application/json", bytes.NewReader(body))
}

func (wp *WhatsappProvider) SendMessage(toNumber string) (*http.Response, error) {
	type language struct {
		Code string `json:"code"`
	}
	type parameter struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	type component struct {
		Type       string      `json:"type"`
		Parameters []parameter `json:"parameters"`
	}
	type expense struct {
		Name       string      `json:"name"`
		Language   language    `json:"language"`
		Components []component `json:"components"`
	}
	type expenseReminderMessageTemplate struct {
		MessagingProduct string  `json:"messaging_product"`
		To               string  `json:"to"`
		Type             string  `json:"type"`
		Template         expense `json:"template"`
	}

	expenseReminderMessage := expenseReminderMessageTemplate{
		MessagingProduct: "whatsapp",
		To:               toNumber,
		Type:             "template",
		Template: expense{
			Name:       "expense_information_accept_message",
			Language:   language{Code: "pt_BR"},
			Components: []component{{Type: "body", Parameters: []parameter{{Type: "Text", Text: "Youtube Premium"}}}},
		},
	}
	body, err := json.Marshal(expenseReminderMessage)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	return wp.whatsappRequest(http.MethodPost, "messages", "application/json", bytes.NewReader(body))
}

/*
func (wp *WhatsappProvider) UploadImage() {
		options := pix.Options{
			Name:          "",
			Key:           "",
			City:          "",
			Amount:        23.69, // optional
			Description:   "",    // optional
			TransactionID: "***", // optional
		}

		copyPaste, _ := pix.Pix(options)
		qrCode, _ := pix.QRCode(pix.QRCodeOptions{Size: 256, Content: copyPaste})
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("messaging_product", "whatsapp")
	file, errFile2 := os.Open("")
	defer file.Close()
	part2,
		errFile2 := writer.CreateFormFile("file", filepath.Base(""))
	_, errFile2 = io.Copy(part2, file)
	if errFile2 != nil {
		fmt.Println(errFile2)
		return
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := wp.whatsappRequest(http.MethodPost, "media", string(writer.FormDataContentType()), payload)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	fmt.Printf("Response: %s", string(body))
}
*/

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
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", wp.UserAccessToken))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return client.Do(req)
}

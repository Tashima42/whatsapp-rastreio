package helpers

type WhatsAppReceivedMessageObject struct {
	Object string                         `json:"object"`
	Entry  []whatsAppReceivedMessageEntry `json:"entry"`
}
type whatsAppReceivedMessageEntry struct {
	Id      string                          `json:"id"`
	Changes []whatsAppReceivedMessageChange `json:"changes"`
}
type whatsAppReceivedMessageChange struct {
	Value whatsAppReceivedMessageChangeValue `json:"value"`
	Field string                             `json:"field"`
}
type whatsAppReceivedMessageChangeValue struct {
	MessagingProduct string                                      `json:"messaging_product"`
	Metadata         whatsAppReceivedMessageChangeValueMetadata  `json:"metadata"`
	Contacts         []whatsAppReceivedMessageChangeValueContact `json:"contacts"`
	Messages         []whatsAppReceivedMessageChangeValueMessage `json:"messages"`
}
type whatsAppReceivedMessageChangeValueMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberId      string `json:"phone_number_id"`
}
type whatsAppReceivedMessageChangeValueContact struct {
	Profile struct {
		Name string `json:"name"`
	} `json:"profile"`
	WaId string `json:"wa_id"`
}

type whatsAppReceivedMessageChangeValueMessage struct {
	Context   whatsAppReceivedMessageChangeValueMessageContext `json:"context"`
	From      string                                           `json:"from"`
	Id        string                                           `json:"id"`
	Timestamp string                                           `json:"timestamp"`
	Type      string                                           `json:"type"`
	Button    whatsAppReceivedMessageChangeValueMessageButton  `json:"button"`
	Text      whatsAppReceivedMessageChangeValueMessageText    `json:"text"`
}
type whatsAppReceivedMessageChangeValueMessageContext struct {
	From string `json:"from"`
	Id   string `json:"id"`
}
type whatsAppReceivedMessageChangeValueMessageButton struct {
	Payload string `json:"payload"`
	Text    string `json:"text"`
}
type whatsAppReceivedMessageChangeValueMessageText struct {
	Body string `json:"body"`
}

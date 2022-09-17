package whatsappTemplates

type Language struct {
	Code string `json:"code"`
}
type Parameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type Component struct {
	Type       string      `json:"type"`
	Parameters []Parameter `json:"parameters"`
}

type Template struct {
	Name       string      `json:"name"`
	Language   Language    `json:"language"`
	Components []Component `json:"components"`
}

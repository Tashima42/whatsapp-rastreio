package whatsappTemplates

type CorreiosInvalidCode struct{}

func (cc *CorreiosInvalidCode) GetTemplate() Template {
	language := Language{Code: "pt_BR"}
	var components []Component
	template := Template{
		Name:       "correios_invalid_code2",
		Language:   language,
		Components: components,
	}
	return template
}

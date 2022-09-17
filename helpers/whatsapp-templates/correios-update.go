package whatsappTemplates

type CorreiosEventUpdate struct {
	Header      string
	Description string
	Date        string
	Local       string
}

func (cu *CorreiosEventUpdate) GetTemplate() Template {
	language := Language{Code: "pt_BR"}
	var components []Component
	var headerParameters []Parameter
	var bodyParameters []Parameter
	headerParameters = append(headerParameters, Parameter{
		Type: "text",
		Text: cu.Header,
	})
	components = append(components, Component{
		Type:       "header",
		Parameters: headerParameters,
	})
	bodyParameters = append(bodyParameters, Parameter{
		Type: "text",
		Text: cu.Description,
	})
	bodyParameters = append(bodyParameters, Parameter{
		Type: "text",
		Text: cu.Date,
	})
	bodyParameters = append(bodyParameters, Parameter{
		Type: "text",
		Text: cu.Local,
	})
	components = append(components, Component{
		Type:       "body",
		Parameters: bodyParameters,
	})
	template := Template{
		Name:       "correios_event_update",
		Language:   language,
		Components: components,
	}
	return template
}

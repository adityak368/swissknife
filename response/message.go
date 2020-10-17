package response

import "github.com/adityak368/swissknife/localization"

// JSONMessage defines a translated JSONMessage
type JSONMessage struct {
	Message string `json:"message"`
}

// Message defines a translatable response message for the caller
type Message struct {
	MessageID   string        `json:"messagId"`
	MessageArgs []interface{} `json:"messageArgs"`
}

// AsTranslatedJSONMessage translates the message and returns a json message which can be sent to user
func (m *Message) AsTranslatedJSONMessage(translator localization.Translator) JSONMessage {
	return JSONMessage{
		Message: translator.Tr(m.MessageID, m.MessageArgs...),
	}
}

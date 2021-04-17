package response

// Message defines a translatable response message for the caller
type Message struct {
	MessageID         string        `json:"-"`
	MessageArgs       []interface{} `json:"-"`
	TranslatedMessage string        `json:"message"`
}

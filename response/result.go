package response

// Result Defines a internal function result
type Result interface {
	Data() interface{}
	ToMessage() *Message
}

// ExecResult defines an internal function result
type ExecResult struct {
	Result      interface{}
	MessageID   string
	MessageArgs []interface{}
}

// Data implements the Result interface
func (e *ExecResult) Data() interface{} {
	return e.Result
}

// Message returns a message which can be translated
func (e *ExecResult) ToMessage() *Message {
	return &Message{
		MessageID:         e.MessageID,
		MessageArgs:       e.MessageArgs,
		TranslatedMessage: e.MessageID,
	}
}

package response

// Error defines a internal function error
type Error struct {
	Code        int
	MessageID   string
	MessageArgs []interface{}
}

// Error returns the error description
func (e *Error) Error() string {
	return e.MessageID
}

// Message returns a message which can be translated
func (e *Error) Message() *Message {
	return &Message{
		MessageID:   e.MessageID,
		MessageArgs: e.MessageArgs,
	}
}

// NewError Creates a new internal function error
func NewError(code int, messageID string, messageArgs ...interface{}) error {
	return &Error{Code: code, MessageID: messageID, MessageArgs: messageArgs}
}

package response

import "github.com/pkg/errors"

// Error defines a internal function error
type Error struct {
	Code        int
	MessageID   string
	MessageArgs []interface{}
	NativeError error
}

// Error returns the error description
func (e *Error) Error() string {
	return e.MessageID
}

// Message returns a message which can be translated
func (e *Error) ToMessage() *Message {
	return &Message{
		MessageID:         e.MessageID,
		MessageArgs:       e.MessageArgs,
		TranslatedMessage: e.MessageID,
	}
}

// NewError Creates a new internal error
func NewError(code int, messageID string, messageArgs ...interface{}) error {
	return &Error{Code: code, MessageID: messageID, MessageArgs: messageArgs}
}

// NewErrorWithDetails Creates a new internal error with the details of the error
func NewErrorWithDetails(err error, code int, messageID string, messageArgs ...interface{}) error {
	wrappedError := errors.WithStack(errors.Wrap(err, ""))
	return &Error{Code: code, MessageID: messageID, NativeError: wrappedError, MessageArgs: messageArgs}
}

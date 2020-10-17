package email

// Mailer defines the emailer interface
type Mailer interface {
	SendMail(from, to, subject, body string) error
	StartDaemon() error
	StopDaemon() error
}

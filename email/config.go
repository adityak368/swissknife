package email

// MailerConfig Configuration to setup the mailer
type MailerConfig struct {
	Host              string
	Port              int
	Username          string
	Password          string
	MaxEmailQueueSize int
}

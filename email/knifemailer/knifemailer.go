package knifemailer

import (
	"errors"

	"github.com/adityak368/swissknife/email"
	"github.com/adityak368/swissknife/logger"
	"gopkg.in/gomail.v2"
)

// knifeMailer implements the Mailer Interface
type knifeMailer struct {
	sendMailChannel chan *gomail.Message
	gomailHandle    gomail.SendCloser
	config          email.MailerConfig
}

// StopDaemon Closes the mail daemon.
func (m *knifeMailer) StopDaemon() error {

	if m.sendMailChannel != nil {
		close(m.sendMailChannel)
	}

	if m.gomailHandle != nil {
		return m.gomailHandle.Close()
	}

	return nil
}

// SendMail sends an email. This is thread safe
func (m *knifeMailer) SendMail(from, to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	//m.Attach("/home/Alex/lolcat.jpg")

	// Send to the channel in a select as the thread might be blocked if buffer is full
	select {
	case m.sendMailChannel <- msg:
	default:
		return errors.New("Email MaxQueueSize reached. Could not send email")
	}
	return nil
}

// StartDaemon Starts the Email daemon
func (m *knifeMailer) StartDaemon() error {
	dialer := gomail.NewDialer(m.config.Host, m.config.Port, m.config.Username, m.config.Password)
	gomailHandle, err := dialer.Dial()
	if err != nil {
		return err
	}
	m.gomailHandle = gomailHandle
	go func() {
		for {
			select {
			case message, ok := <-m.sendMailChannel:
				if !ok {
					logger.Info("Closed Email Daemon")
					return
				}
				if m.gomailHandle != nil {
					if err := gomail.Send(m.gomailHandle, message); err != nil {
						logger.Trace(err)
					}
				}
			}
		}
	}()
	logger.Info("Started Email Daemon")
	return nil
}

// New Returns a new mailer from the given config
func New(config email.MailerConfig) email.Mailer {
	return &knifeMailer{
		config:          config,
		sendMailChannel: make(chan *gomail.Message, config.MaxEmailQueueSize),
	}
}

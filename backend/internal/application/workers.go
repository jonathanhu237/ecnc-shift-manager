package application

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wneessen/go-mail"
)

type MailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (app *Application) StartMailSender(ctx context.Context, ch *amqp.Channel) error {
	// establish mail client
	mailClient, err := mail.NewClient(
		app.config.MailClient.SMTPHost,
		mail.WithPort(465),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithSSL(),
		mail.WithUsername(app.config.MailClient.Sender),
		mail.WithPassword(app.config.MailClient.Password),
	)
	if err != nil {
		return err
	}

	diaCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := mailClient.DialWithContext(diaCtx); err != nil {
		return err
	}

	// consume messages
	msgs, err := ch.Consume("mail_queue", "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d := <-msgs:
				// parse the message
				var mailPayload MailPayload
				if err := json.Unmarshal(d.Body, &mailPayload); err != nil {
					app.logger.Error("failed to unmarshal mail", slog.String("error", err.Error()))
					if err := d.Nack(false, false); err != nil {
						app.logger.Error("failed to nack message", slog.String("error", err.Error()))
					}
					continue
				}

				// send the mail
				message := mail.NewMsg()
				message.Subject(mailPayload.Subject)
				message.SetBodyString(mail.TypeTextPlain, mailPayload.Body)
				if err := message.From(app.config.MailClient.Sender); err != nil {
					app.logger.Error("failed to set mail sender", slog.String("error", err.Error()))
					if err := d.Nack(false, false); err != nil {
						app.logger.Error("failed to nack message", slog.String("error", err.Error()))
					}
					continue
				}
				if err := message.To(mailPayload.To); err != nil {
					app.logger.Error("failed to set mail to", slog.String("error", err.Error()))
					if err := d.Nack(false, false); err != nil {
						app.logger.Error("failed to nack message", slog.String("error", err.Error()))
					}
					continue
				}

				if err := mailClient.Send(message); err != nil {
					app.logger.Error("failed to send mail", slog.String("error", err.Error()))
					if err := d.Nack(false, true); err != nil {
						app.logger.Error("failed to nack message", slog.String("error", err.Error()))
					}
					continue
				}

				if err := d.Ack(false); err != nil {
					app.logger.Error("failed to ack message", slog.String("error", err.Error()))
				}
			}
		}
	}()

	return nil
}

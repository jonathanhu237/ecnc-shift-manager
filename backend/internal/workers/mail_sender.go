package workers

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wneessen/go-mail"
)

type MailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type MailSender struct {
	config *config.Config
	logger *slog.Logger
	ch     *amqp.Channel
}

func NewMailSender(config *config.Config, logger *slog.Logger, ch *amqp.Channel) *MailSender {
	return &MailSender{
		config: config,
		logger: logger,
		ch:     ch,
	}
}

func (ms *MailSender) Run(ctx context.Context) error {
	// establish mail client
	mailClient, err := mail.NewClient(
		ms.config.MailClient.SMTPHost,
		mail.WithPort(465),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithSSL(),
		mail.WithUsername(ms.config.MailClient.Sender),
		mail.WithPassword(ms.config.MailClient.Password),
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
	msgs, err := ms.ch.Consume("mail_queue", "", false, false, false, false, nil)
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
					ms.logger.Error("failed to unmarshal mail", slog.String("error", err.Error()))
					if err := d.Nack(false, false); err != nil {
						ms.logger.Error("failed to nack message", slog.String("error", err.Error()))
					}
					continue
				}

				// send the mail
				message := mail.NewMsg()
				message.Subject(mailPayload.Subject)
				message.SetBodyString(mail.TypeTextPlain, mailPayload.Body)

				if err := message.From(ms.config.MailClient.Sender); err != nil {
					ms.logger.Error("failed to set mail sender", slog.String("error", err.Error()))
					if err := d.Nack(false, false); err != nil {
						ms.logger.Error("failed to nack message", slog.String("error", err.Error()))
					}
					continue
				}

				if err := message.To(mailPayload.To); err != nil {
					ms.logger.Error("failed to set mail to", slog.String("error", err.Error()))
					if err := d.Nack(false, false); err != nil {
						ms.logger.Error("failed to nack message", slog.String("error", err.Error()))
					}
					continue
				}

				if err := mailClient.Send(message); err != nil {
					ms.logger.Error("failed to send mail", slog.String("error", err.Error()))
					if err := d.Nack(false, true); err != nil {
						ms.logger.Error("failed to nack message", slog.String("error", err.Error()))
					}
					continue
				}

				if err := d.Ack(false); err != nil {
					ms.logger.Error("failed to ack message", slog.String("error", err.Error()))
				}
			}
		}
	}()

	return nil
}

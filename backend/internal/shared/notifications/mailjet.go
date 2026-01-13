package notifications

import (
	"alloy/internal/shared/config"
	"errors"

	"github.com/mailjet/mailjet-apiv3-go/v4"
	"go.uber.org/zap"
)

type MailjetService struct {
	publicKey     string
	privateKey    string
	emailFrom     string
	emailFromName string
	client        *mailjet.Client
	logger        *zap.Logger
}

func NewMailjetService(cfg *config.Config, logger *zap.Logger) (*MailjetService, error) {
	if cfg.MailjetPublicKey == "" || cfg.MailjetPrivateKey == "" || cfg.EmailFrom == "" || cfg.EmailFromName == "" {
		return nil, errors.New("mailjet public key, private key, email from or email from name is not set")
	}

	client := mailjet.NewMailjetClient(cfg.MailjetPublicKey, cfg.MailjetPrivateKey)

	return &MailjetService{
		publicKey:     cfg.MailjetPublicKey,
		privateKey:    cfg.MailjetPrivateKey,
		emailFrom:     cfg.EmailFrom,
		emailFromName: cfg.EmailFromName,
		client:        client,
		logger:        logger,
	}, nil
}

func (m *MailjetService) Send(payload *NotificationPayload) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: m.emailFrom,
				Name:  m.emailFromName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: payload.To,
					Name:  payload.To,
				},
			},
			Subject:  payload.Subject,
			TextPart: payload.Body,
			HTMLPart: payload.HTML,
		},
	}

	messages := mailjet.MessagesV31{
		Info: messagesInfo,
	}

	m.logger.Info("Sending email", zap.String("to", payload.To), zap.String("subject", payload.Subject))
	_, err := m.client.SendMailV31(&messages)
	if err != nil {
		m.logger.Error("Failed to send email", zap.Error(err))
		return err
	}

	m.logger.Info("Email sent successfully", zap.String("to", payload.To), zap.String("subject", payload.Subject))
	return nil
}

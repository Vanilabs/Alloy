package notifications

import (
	"alloy/internal/shared/config"
	"alloy/internal/shared/constants"
	internalTemplate "alloy/internal/shared/template"
	"errors"

	"go.uber.org/zap"
)

type Email struct {
	service        NotificationService
	TemplateParser *internalTemplate.TemplateParser
}

func NewEmail(cfg *config.Config, logger *zap.Logger, service string) (*Email, error) {
	var choiceService NotificationService
	var err error

	switch service {
	case constants.MAILJET_SERVICE:
		choiceService, err = NewMailjetService(cfg, logger)
		if err != nil {
			logger.Error("Failed to create mailjet service", zap.Error(err))
			return nil, err
		}
	default:
		return nil, errors.New("invalid service")
	}
	return &Email{
		service:        choiceService,
		TemplateParser: internalTemplate.NewTemplateParser("email_templates"),
	}, nil
}

func (e *Email) Send(payload *NotificationPayload) error {
	return e.service.Send(payload)
}

package auth

import (
	"alloy/internal/modules/users"
	"alloy/internal/shared/notifications"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Module struct {
	Repository Repository
	Service    Service
	Handler    *Handler
}

func NewModule(db *gorm.DB, logger *zap.Logger, notification *notifications.Notification, userRepository users.Repository) *Module {
	repository := NewRepository(db)
	service := NewService(repository, logger, notification, userRepository)
	handler := NewHandler(service)

	return &Module{
		Repository: repository,
		Service:    service,
		Handler:    handler,
	}
}

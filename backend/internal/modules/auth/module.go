package auth

import (
	"alloy/internal/modules/users"
	"alloy/internal/shared/constants"
	"alloy/internal/shared/notifications"

	"go.uber.org/zap"
)

type Module struct {
	Repository Repository
	Service    Service
	Handler    *Handler
}

func NewModule(store *constants.DataStores, logger *zap.Logger, notification *notifications.Notification, userRepository users.Repository) *Module {

	repository := NewRepository(store, logger)
	service := NewService(repository, logger, notification, userRepository)
	handler := NewHandler(service)

	return &Module{
		Repository: repository,
		Service:    service,
		Handler:    handler,
	}
}

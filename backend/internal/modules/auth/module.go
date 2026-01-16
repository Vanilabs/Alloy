package auth

import (
	"alloy/internal/modules/users"
	"alloy/internal/shared/constants"
	"alloy/internal/shared/notifications"
	"alloy/internal/shared/utils"

	"go.uber.org/zap"
)

type Module struct {
	Repository Repository
	Service    Service
	Handler    *Handler
}

func NewModule(store *constants.DataStores, logger *zap.Logger, notification *notifications.Notification, userRepository users.Repository, jwtManager *utils.JWTManager) *Module {

	repository := NewRepository(store, logger)
	service := NewService(repository, logger, notification, userRepository, jwtManager)
	handler := NewHandler(service)

	return &Module{
		Repository: repository,
		Service:    service,
		Handler:    handler,
	}
}

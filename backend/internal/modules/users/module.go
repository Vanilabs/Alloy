package users

import (
	"alloy/internal/shared/constants"
	"alloy/internal/shared/notifications"

	"go.uber.org/zap"
)

// Module represents the users module with all its dependencies
type Module struct {
	Repository Repository
	Service    Service
	Handler    *Handler
}

// NewModule creates and initializes the users module
func NewModule(store *constants.DataStores, logger *zap.Logger, notification *notifications.Notification) *Module {
	repository := NewRepository(store.PostGres)
	service := NewService(repository, logger, notification)
	handler := NewHandler(service)

	return &Module{
		Repository: repository,
		Service:    service,
		Handler:    handler,
	}
}

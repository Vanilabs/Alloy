package messaging

import (
	"alloy/internal/shared/socket"

	"alloy/internal/modules/users"

	"alloy/internal/shared/router"
)

// Module represents the messaging module with all its dependencies
type Module struct {
	Handler    *Handler
	Service  Service
	Repository Repository
}

// NewModule creates and initializes the messaging module
func NewModule(socketManager *socket.ConnectionManager, env *router.Environment) *Module {
	ms_repository := NewRepository(env.Stores.Cassandra, env.Stores.PostGres)
	user_repository := users.NewRepository(env.Stores.PostGres)
	service := NewService(ms_repository, user_repository, socketManager, env)
	handler := NewHandler(service, socketManager)

	return &Module{
		Handler:    handler,
		Service: service,
		Repository: ms_repository,
	}
}

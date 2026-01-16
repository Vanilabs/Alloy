package app

import (
	"alloy/internal/modules/auth"
	"alloy/internal/modules/messaging"
	"alloy/internal/modules/users"
	"alloy/internal/shared/config"
	"alloy/internal/shared/constants"
	"alloy/internal/shared/notifications"
	"alloy/internal/shared/router"

	"alloy/internal/shared/socket"

	"go.uber.org/zap"
)

// this holds all initialized modules...
type Modules struct {
	Users     *users.Module
	Messaging *messaging.Module
	Auth      *auth.Module
}

func InitModules(env *router.Environment) (*Modules, error) {
	env.Logger.Info("Initializing modules...")

	socketTracker := socket.NewSocketTracker(env.Stores.Redis)
	socketManager := socket.NewManager(socketTracker)

	notification, err := RegisterNotifications(env.Config, env.Logger)
	if err != nil {
		return nil, err
	}

	usersModule := users.NewModule(env.Stores, env.Logger, notification)
	messagingModule := messaging.NewModule(socketManager, env)
	authModule := auth.NewModule(env.Stores, env.Logger, notification, usersModule.Repository)

	return &Modules{
		Users:     usersModule,
		Messaging: messagingModule,
		Auth:      authModule,
	}, nil

}
func RegisterNotifications(cfg *config.Config, logger *zap.Logger) (*notifications.Notification, error) {

	email, err := notifications.NewEmail(cfg, logger, constants.MAILJET_SERVICE)
	if err != nil {
		logger.Error("Failed to create email", zap.Error(err))
		return nil, err
	}

	return notifications.NewNotification(email), nil
}

func RegisterServices(modules *Modules, services *router.ModuleServices) {
	services.Register(constants.USERS_MODULE_NAME, modules.Users.Service)
	services.Register(constants.MESSAGING_MODULE_NAME, modules.Messaging.Service)
	services.Register(constants.AUTH_MODULE_NAME, modules.Auth.Service)
}

func GetHandlers(modules *Modules) []router.IHandler {
	return []router.IHandler{
		modules.Users.Handler,
		modules.Messaging.Handler,
		modules.Auth.Handler,
	}
}

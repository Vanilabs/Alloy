package main

import (
	"alloy/internal/app"
	"alloy/internal/shared/config"
	"alloy/internal/shared/database"
	"alloy/internal/shared/cache"
	"alloy/internal/shared/router"
	"alloy/internal/shared/constants"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {

	cfg := config.LoadAppConfig()
	fmt.Println(cfg)

	zapLogger := config.InitLogger()
	defer zapLogger.Close()

	zapLogger.Logger.Info("Welcome to Alloy Backend")

	fiberApp := router.InitRouterWithConfig(cfg, zapLogger.Logger)

	rds, err :=  cache.GetRedisClient(cfg, zapLogger.Logger, 0)
	if err != nil{
		zapLogger.Logger.Error("Failed to connect to Redis", zap.Error(err))
		os.Exit(1)
	}

	pgdb, err := database.ConnectDB(cfg, zapLogger.Logger)
	if err != nil {
		zapLogger.Logger.Error("Failed to connect to postgres database", zap.Error(err))
		os.Exit(1)
	}

	cassandradb, err := database.NewCassandraSession(cfg)
	if err != nil {
		zapLogger.Logger.Error("Failed to connect to cassandra database", zap.Error(err))
		os.Exit(1)
	}

	defer cassandradb.Close()

	err = database.CreateCassandraEntities(cassandradb)
	if err != nil {
		zapLogger.Logger.Error("Failed to create entites in cassandra", zap.Error(err))
		os.Exit(1)
	}

	stores := &constants.DataStores{
		Redis: rds,
		PostGres: pgdb,
		Cassandra: cassandradb,
	}

	services := router.NewModuleServices()
		
	// create environment...
	env := router.NewEnvironment(cfg, fiberApp, zapLogger.Logger, stores, services)

	// initialize all modules...
	modules, err := app.InitModules(env)
	if err != nil {
		zapLogger.Logger.Error("Failed to initialize modules", zap.Error(err))
		os.Exit(1)
	}

	// create module services registry...
	app.RegisterServices(modules, services)

	// initialize handlers...
	router.InitHandlers(env, app.GetHandlers(modules))

	go func() {
		router.RunWithGracefulShutdown(fiberApp, cfg.PORT, env.Logger)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

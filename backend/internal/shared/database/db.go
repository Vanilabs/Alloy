package database

import (
	"fmt"
	"log"

	"alloy/internal/shared/config"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"

	gocqlastra "github.com/datastax/gocql-astra"
	"github.com/gocql/gocql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"time"
	"strconv"
)

func ConnectDB(cfg *config.Config, zapLogger *zap.Logger) (*gorm.DB, error) {
	schema := cfg.PostgresSchema
	if schema == "" {
		schema = "public"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable options='-c search_path=%s'",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPass,
		cfg.PostgresDB,
		cfg.PostgresPort,
		schema,
	)

	if cfg.PostgresDSN != "" {
		dsn = cfg.PostgresDSN
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Use silent mode in production
	if cfg.APP_MODE == "production" {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	zapLogger.Info("Database connection established")

	return db, nil
}


func NewCassandraSession(cfg *config.Config) (*gocql.Session, error) {
	cluster := gocql.NewCluster(cfg.CassandraHost)

	if cfg.CassandraPort != "" {
		port, err := strconv.Atoi(cfg.CassandraPort)
		if err != nil {
			log.Fatalf("Invalid PORT value: %v", err)
	}
		cluster.Port = port
	}

	if cfg.CassandraUsername != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cfg.CassandraUsername,
			Password: cfg.CassandraPassword,
		}
	}
	if cfg.AstraAppToken != "" {
		var err error
		cluster, err = gocqlastra.NewClusterFromURL("https://api.astra.datastax.com", cfg.AstraDbId, cfg.AstraAppToken, 50*time.Second)

    	if err != nil {
        	log.Fatalf("unable to load cluster in region %s from astra: %v", cfg.AstraDbRegion, err)
    }
	}
	cluster.Consistency = gocql.LocalQuorum
	cluster.Timeout = 5 * time.Second
	cluster.ConnectTimeout = 5 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	keyspace := cfg.CassandraKeyspace

	if cfg.AstraAppToken == "" {
		cql := fmt.Sprintf(`
CREATE KEYSPACE IF NOT EXISTS %s
WITH replication = {
  'class': 'SimpleStrategy',
  'replication_factor': '1'
};
`, keyspace)
		if err := session.Query(cql).Exec(); err != nil {
	log.Fatal("Failed to create keyspace:", err)
	}
	}

	session.Close()

	cluster.Keyspace = cfg.CassandraKeyspace

	return cluster.CreateSession()
}
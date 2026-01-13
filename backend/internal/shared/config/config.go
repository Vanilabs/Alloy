package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func parseEnv() error {

	err := godotenv.Overload(".env")
	if err != nil && !os.IsNotExist(err) {
		log.Println("godotenv.Load().Error:", err)
		return nil
	}

	missing := make([]string, 0)
	envVars := []string{}

	for _, v := range envVars {
		envVal := os.Getenv(v)
		if envVal == "" {
			missing = append(missing, v)
		}
	}

	if len(missing) == 0 {
		return nil
	}

	return fmt.Errorf("missing env vars: %v", missing)
}

func LoadAppConfig() *Config {
	env := parseEnv()
	if env != nil {
		log.Fatal("Failed to load app config:", env)
	}

	jwtExpiry, _ := strconv.ParseInt(os.Getenv("JWT_EXPIRY"), 10, 64)

	if jwtExpiry == 0 {
		jwtExpiry = 1 // Default to 1 hour
	}

	jwtRefreshExpiry, _ := strconv.ParseInt(os.Getenv("JWT_REFRESH_EXPIRY"), 10, 64)
	if jwtRefreshExpiry == 0 {
		jwtRefreshExpiry = 7 // Default to 1 week
	}

	config := Config{
		APP_MODE:              os.Getenv("APP_MODE"),
		PostgresDB:            os.Getenv("POSTGRES_DB"),
		PostgresUser:          os.Getenv("POSTGRES_USER"),
		PostgresPass:          os.Getenv("POSTGRES_PASSWORD"),
		PostgresDSN:           os.Getenv("POSTGRES_DSN"),
		RedisAddr:             os.Getenv("REDIS_ADDR"),
		RedisPassword:         os.Getenv("REDIS_PASSWORD"),
		RedisUrl:              os.Getenv("REDIS_URL"),
		PORT:                  os.Getenv("PORT"),
		ORIGINS:               os.Getenv("ORIGINS"),
		PostgresHost:          os.Getenv("POSTGRES_HOST"),
		PostgresPort:          os.Getenv("POSTGRES_PORT"),
		PostgresSchema:        os.Getenv("POSTGRES_SCHEMA"),
		JwtSecret:             os.Getenv("JWT_SECRET"),
		VerboseRequestLogging: os.Getenv("VERBOSE_REQUEST_LOGGING") == "true",
		CassandraHost:    		os.Getenv("CASSANDRA_HOST"),
		CassandraKeyspace: 		os.Getenv("CASSANDRA_KEYSPACE"),
		CassandraPassword:		os.Getenv("CASSANDRA_PASSWORD"),
		CassandraUsername:		os.Getenv("CASSANDRA_USERNAME"),
		CassandraPort : 		os.Getenv("CASSANDRA_PORT"),
		MailjetPublicKey:      os.Getenv("MAILJET_PUBLIC_KEY"),
		MailjetPrivateKey:     os.Getenv("MAILJET_PRIVATE_KEY"),
		EmailFrom:             os.Getenv("EMAIL_FROM"),
		EmailFromName:         os.Getenv("EMAIL_FROM_NAME"),
		AstraDbId: os.Getenv("ASTRA_DB_ID"),
		AstraDbRegion: os.Getenv("ASTRA_DB_REGION"),
		AstraAppToken: os.Getenv("ASTRA_APP_TOKEN"),
	}

	return &config
}

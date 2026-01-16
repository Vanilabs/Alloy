package config

type Config struct {
	APP_MODE              string `mapstructure:"APP_MODE"`
	PostgresDB            string `mapstructure:"POSTGRES_DB"`
	PostgresUser          string `mapstructure:"POSTGRES_USER"`
	PostgresPass          string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDSN           string `mapstructure:"POSTGRES_DSN"`
	PostgresHost          string `mapstructure:"POSTGRES_HOST"`
	PostgresPort          string `mapstructure:"POSTGRES_PORT"`
	PostgresSchema        string `mapstructure:"POSTGRES_SCHEMA"`
	RedisAddr             string `mapstructure:"REDIS_ADDR"`
	RedisUrl              string `mapstructure:"REDIS_URL"`
	RedisPassword         string `mapstructure:"REDIS_PASSWORD"`
	PORT                  string `mapstructure:"PORT"`
	ORIGINS               string `mapstructure:"ORIGINS"`
	JwtSecret             string `mapstructure:"JWT_SECRET"`
	VerboseRequestLogging bool   `mapstructure:"VERBOSE_REQUEST_LOGGING"`
	CassandraHost         string `mapstructure:"CASSANDRA_HOST"`
	CassandraKeyspace     string `mapstructure:"CASSANDRA_KEYSPACE"`
	CassandraPassword     string `mapstructure:"CASSANDRA_PASSWORD"`
	CassandraUsername     string `mapstructure:"CASSANDRA_USERNAME"`
	CassandraPort         string `mapstructure:"CASSANDRA_PORT"`
	MailjetPublicKey      string `mapstructure:"MAILJET_PUBLIC_KEY"`
	MailjetPrivateKey     string `mapstructure:"MAILJET_PRIVATE_KEY"`
	EmailFrom             string `mapstructure:"EMAIL_FROM"`
	EmailFromName         string `mapstructure:"EMAIL_FROM_NAME"`
	AstraDbId             string `mapstructure:"ASTRA_DB_ID"`
	AstraDbRegion         string `mapstructure:"ASTRA_DB_REGION"`
	AstraAppToken         string `mapstructure:"ASTRA_APP_TOKEN"`
	RefreshTokenSecret    string `mapstructure:"REFRESH_TOKEN_SECRET"`
}

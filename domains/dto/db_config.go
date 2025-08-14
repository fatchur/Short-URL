package dto

type DBConfig struct {
	DSN      string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Timezone string
	LogLevel string
}
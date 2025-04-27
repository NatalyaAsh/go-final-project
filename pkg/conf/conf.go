package conf

import "os"

const (
	DBFILE = "scheduler.db"
	PORT   = "7540"
)

type Configuration struct {
	DBFile   string
	Port     string
	Password string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func New() *Configuration {
	var cfg Configuration
	cfg.DBFile = getEnv("TODO_DBFILE", DBFILE)
	cfg.Port = getEnv("TODO_PORT", PORT)
	cfg.Password = os.Getenv("TODO_PASSWORD")
	return &cfg
}

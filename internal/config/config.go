package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Database struct {
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     uint16 `env:"POSTGRES_PORT"`
	PostgresDatabase string `env:"POSTGRES_DATABASE"`
	TestDatabaseURL  string `env:"TEST_DATABASE_URL"`
}

func (db Database) DSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", db.PostgresUser, db.PostgresPassword, db.PostgresHost, db.PostgresPort, db.PostgresDatabase)
}

type Auth struct {
	JwtSecret   string `env:"JWT_SECRET"`
	SecretToken string `env:"SECRET_TOKEN"`
}

type Config struct {
	Database
	Auth
}

func FromEnv() (Config, error) {
	var cfg Config
	return cfg, cleanenv.ReadEnv(&cfg)
}

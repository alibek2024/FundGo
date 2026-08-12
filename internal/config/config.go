package config

import (
	"log/slog"
	"time"
)

type Config struct {
	App    AppConfig
	DB     DBConfig
	JWT    JWTConfig
	Redis  RedisConfig
	Logger LoggerConfig
}

type LoggerConfig struct {
	logger slog.Logger
}

type AppConfig struct {
	Env  string
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

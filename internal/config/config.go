package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	App    AppConfig
	DB     DBConfig
	JWT    JWTConfig
	Redis  RedisConfig
	Logger LoggerConfig
}

type LoggerConfig struct {
	Level  string `env:"LOGGER_LEVEL" envDefault:"info"`
	Format string `env:"LOGGER_FORMAT" envDefault:"json"`
}

type AppConfig struct {
	Env  string `env:"APP_ENV" envDefault:"development"`
	Host string `env:"APP_HOST" envDefault:"-"`
	Port string `env:"APP_PORT" envDefault:"8080"`
}

type DBConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Name     string `env:"DB_NAME,required"`
	SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
}

type JWTConfig struct {
	PrivateKeyPath  string        `env:"JWT_PRIVATE_KEY_PATH,required"`
	PublicKeyPath   string        `env:"JWT_PUBLIC_KEY_PATH,required"`
	AccessTokenTTL  time.Duration `env:"JWT_ACCESS_TOKEN_TTL" envDefault:"15m"`
	RefreshTokenTTL time.Duration `env:"JWT_REFRESH_TOKEN_TTL" envDefault:"720h"`
}

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

func (c JWTConfig) RsaKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKeyData, err := os.ReadFile(c.PrivateKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read private key: %w", err)
	}

	privateBlock, _ := pem.Decode(privateKeyData)
	if privateBlock == nil {
		return nil, nil, errors.New("failed to decode private key PEM")
	}

	parsedPrivateKey, err := x509.ParsePKCS8PrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse private key: %w", err)
	}

	privateKey, ok := parsedPrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, errors.New("private key is not RSA")
	}

	publicKeyData, err := os.ReadFile(c.PublicKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read public key: %w", err)
	}

	publicBlock, _ := pem.Decode(publicKeyData)
	if publicBlock == nil {
		return nil, nil, errors.New("failed to decode public key PEM")
	}

	parsedPublicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse public key: %w", err)
	}

	publicKey, ok := parsedPublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, nil, errors.New("public key is not RSA")
	}

	return privateKey, publicKey, nil
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("failed to load .env: %w", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

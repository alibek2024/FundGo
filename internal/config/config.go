package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

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
	Level  string
	Format string
}

type AppConfig struct {
	Env  string
	Host string
	Port string
}

type DBConfig struct {
	URL string
}

type JWTConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string

	PrivateKeyB64 string
	PublicKeyB64  string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type RedisConfig struct {
	URL string
	TLS bool
}

func (c JWTConfig) readPrivateKey() ([]byte, error) {
	if c.PrivateKeyB64 != "" {
		data, err := base64.StdEncoding.DecodeString(c.PrivateKeyB64)
		if err != nil {
			return nil, fmt.Errorf("decode private key: %w", err)
		}

		return data, nil
	}

	return os.ReadFile(c.PrivateKeyPath)
}

func (c JWTConfig) readPublicKey() ([]byte, error) {
	if c.PublicKeyB64 != "" {
		data, err := base64.StdEncoding.DecodeString(c.PublicKeyB64)
		if err != nil {
			return nil, fmt.Errorf("decode public key: %w", err)
		}

		return data, nil
	}

	return os.ReadFile(c.PublicKeyPath)
}

func (c JWTConfig) RsaKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKeyData, err := c.readPrivateKey()
	if err != nil {
		return nil, nil, err
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

	publicKeyData, err := c.readPublicKey()
	if err != nil {
		return nil, nil, err
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

func getEnv(key string) string {
	return os.Getenv(key)
}

func Load() (Config, error) {
	_ = godotenv.Load()

	accessTTL, err := time.ParseDuration(getEnv("JWT_ACCESS_TOKEN_TTL"))
	if err != nil {
		return Config{}, fmt.Errorf("parse JWT_ACCESS_TOKEN_TTL: %w", err)
	}

	refreshTTL, err := time.ParseDuration(getEnv("JWT_REFRESH_TOKEN_TTL"))
	if err != nil {
		return Config{}, fmt.Errorf("parse JWT_REFRESH_TOKEN_TTL: %w", err)
	}

	redisTLS, err := strconv.ParseBool(getEnv("REDIS_TLS"))
	if err != nil {
		return Config{}, fmt.Errorf("parse REDIS_TLS: %w", err)
	}

	cfg := Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV"),
			Host: getEnv("APP_HOST"),
			Port: getEnv("PORT"),
		},

		DB: DBConfig{
			URL: getEnv("DATABASE_URL"),
		},

		JWT: JWTConfig{
			PrivateKeyPath:  getEnv("JWT_PRIVATE_KEY_PATH"),
			PublicKeyPath:   getEnv("JWT_PUBLIC_KEY_PATH"),
			PrivateKeyB64:   getEnv("JWT_PRIVATE_KEY_B64"),
			PublicKeyB64:    getEnv("JWT_PUBLIC_KEY_B64"),
			AccessTokenTTL:  accessTTL,
			RefreshTokenTTL: refreshTTL,
		},

		Redis: RedisConfig{
			URL: getEnv("REDIS_URL"),
			TLS: redisTLS,
		},

		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL"),
			Format: getEnv("LOG_FORMAT"),
		},
	}

	if cfg.DB.URL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	if cfg.Redis.URL == "" {
		return Config{}, errors.New("REDIS_URL is required")
	}

	if cfg.App.Port == "" {
		cfg.App.Port = "8080"
	}

	if cfg.App.Env == "" {
		cfg.App.Env = "development"
	}

	return cfg, nil
}

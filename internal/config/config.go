package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppEnv   string
	HTTPAddr string

	DBHost string
	DBPort int
	DBName string
	DBUser string
	DBPass string

	JWTSecret     string
	JWTTTLMinutes int
}

func Load() (Config, error) {
	c := Config{
		AppEnv:   getenv("APP_ENV", "dev"),
		HTTPAddr: getenv("HTTP_ADDR", ":8080"),

		DBHost: getenv("DB_HOST", "127.0.0.1"),
		DBName: getenv("DB_NAME", "daycare"),
		DBUser: getenv("DB_USER", "daycare_user"),
		DBPass: getenv("DB_PASS", "daycare_pass"),

		JWTSecret: getenv("JWT_SECRET", "change_me"),
	}
	var err error
	c.DBPort, err = atoi(getenv("DB_PORT", "3306"))
	if err != nil {
		return Config{}, fmt.Errorf("DB_PORT: %w", err)
	}
	c.JWTTTLMinutes, err = atoi(getenv("JWT_TTL_MINUTES", "720"))
	if err != nil {
		return Config{}, fmt.Errorf("JWT_TTL_MINUTES: %w", err)
	}
	return c, nil
}

func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName,
	)
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func atoi(s string) (int, error) { return strconv.Atoi(s) }

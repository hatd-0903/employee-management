package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port      string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	ExportDir string
	AuthUsers map[string]string
}

// Load reads configuration from environment variables, applying defaults
// suitable for local docker-compose usage.
func Load() Config {
	return Config{
		Port:      getenv("APP_PORT", "8080"),
		DBHost:    getenv("DB_HOST", "127.0.0.1"),
		DBPort:    getenv("DB_PORT", "3306"),
		DBUser:    getenv("DB_USER", "root"),
		DBPass:    getenv("DB_PASSWORD", ""),
		DBName:    getenv("DB_NAME", "employee_management"),
		ExportDir: getenv("EXPORT_DIR", "exports"),
		AuthUsers: parseUsers(getenv("AUTH_USERS", "admin:admin123")),
	}
}

// DSN builds the go-sql-driver/mysql data source name.
func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// parseUsers parses a "user1:pass1,user2:pass2" string into a lookup map.
func parseUsers(raw string) map[string]string {
	users := make(map[string]string)
	for _, pair := range strings.Split(raw, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			continue
		}
		users[parts[0]] = parts[1]
	}
	return users
}

package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	Env                 string
	CookieDomain        string
	DatabaseURL         string
	RedisURL            string
	RedisPassword       string
	JWTSecret           string
	JWTAccessTokenTTL   time.Duration
	JWTRefreshTokenTTL  time.Duration
	S3Bucket            string
	S3Endpoint          string
	AWSAccessKeyID      string
	AWSSecretAccessKey  string
	AWSS3Region         string
	S3UsePathStyle      bool
	CORSOrigin          string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	accessTTL, err := time.ParseDuration(getEnv("JWT_ACCESS_TOKEN_TTL", "2h"))
	if err != nil {
		accessTTL = 2 * time.Hour
	}

	refreshTTL, err := time.ParseDuration(getEnv("JWT_REFRESH_TOKEN_TTL", "168h"))
	if err != nil {
		refreshTTL = 7 * 24 * time.Hour
	}

	usePathStyle, _ := strconv.ParseBool(getEnv("S3_USE_PATH_STYLE", "false"))

	cfg := &Config{
		Port:               getEnv("PORT", "4242"),
		Env:                getEnv("ENV", getEnv("NODE_ENV", "development")),
		CookieDomain:       getEnv("COOKIE_DOMAIN", "localhost"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgresql://postgres:password@localhost:5432/schedule?sslmode=disable"),
		RedisURL:           getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		JWTSecret:          getEnv("JWT_SECRET", "super-secret-schedule-jwt-key"),
		JWTAccessTokenTTL:  accessTTL,
		JWTRefreshTokenTTL: refreshTTL,
		S3Bucket:           getEnv("S3_BUCKET", "schedule-test-bucket"),
		S3Endpoint:         getEnv("S3_ENDPOINT", "https://storage.yandexcloud.net"),
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSS3Region:        getEnv("AWS_S3_REGION", "ru-central1-a"),
		S3UsePathStyle:     usePathStyle,
		CORSOrigin:         getEnv("CORS_ORIGIN", "http://localhost:3000"),
	}

	return cfg, nil
}

func (c *Config) IsDev() bool {
	return c.Env == "development" || c.Env == "dev"
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

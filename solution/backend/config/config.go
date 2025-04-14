package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerAddress         string
	DBHost                string
	DBPort                string
	DBUser                string
	DBPassword            string
	DBName                string
	RedisHost             string
	RedisPort             string
	S3Bucket              string
	S3Region              string
	S3AccessKey           string
	S3SecretKey           string
	S3Endpoint            string
	S3Public              string
	GigaChatAuthBaseUrl   string
	GigaChatModelsBaseUrl string
	GigaChatAuthKey       string
	OpenAIToken           string
	OpenAIBaseUrl         string
	Moderation            bool
}

func LoadConfig() Config {
	return Config{
		ServerAddress:         os.Getenv("SERVER_ADDRESS"),
		DBHost:                os.Getenv("MONGODB_HOST"),
		DBPort:                os.Getenv("MONGODB_PORT"),
		DBUser:                os.Getenv("MONGODB_USER"),
		DBPassword:            os.Getenv("MONGODB_PASSWORD"),
		DBName:                os.Getenv("MONGODB_DB"),
		RedisHost:             os.Getenv("REDIS_HOST"),
		RedisPort:             os.Getenv("REDIS_PORT"),
		S3Bucket:              os.Getenv("S3_BUCKET"),
		S3Region:              os.Getenv("S3_REGION"),
		S3AccessKey:           os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey:           os.Getenv("S3_SECRET_KEY"),
		S3Endpoint:            os.Getenv("S3_ENDPOINT"),
		S3Public:              os.Getenv("S3_PUBLIC"),
		GigaChatAuthBaseUrl:   os.Getenv("GIGA_CHAT_AUTH_BASE_URL"),
		GigaChatModelsBaseUrl: os.Getenv("GIGA_CHAT_MODELS_BASE_URL"),
		GigaChatAuthKey:       os.Getenv("GIGA_CHAT_AUTH_KEY"),
		OpenAIToken:           os.Getenv("OPENAI_TOKEN"),
		OpenAIBaseUrl:         os.Getenv("OPENAI_BASE_URL"),
		Moderation:            os.Getenv("MODERATION") == "true",
	}
}

func (cfg *Config) BuildDsn() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)
}

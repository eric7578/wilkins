package storage

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var redisClient *redis.Client

// InitClient create single instance of redisClient
func InitClient(host string, password string, db int) *redis.Client {
	if redisClient != nil {
		return redisClient
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	return redisClient
}

// LoadRedisEnv load redis related env from .env file
func LoadRedisEnv() (string, string, int) {
	env := os.ExpandEnv("$GOPATH/src/github.com/eric7578/wilkins/.env")
	godotenv.Load(env)

	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if redisHost == "" {
		log.Fatal(fmt.Errorf("invalid redis config: REDIS_HOST = %s", redisHost))
	}

	if err != nil {
		log.Info("REDIS_DB is invalid, use 0 instead")
		redisDB = 0
	}

	return redisHost, redisPassword, redisDB
}

func mustGetRedisClient() *redis.Client {
	if redisClient == nil {
		log.Fatal("redisClient is not initialized")
	}
	return redisClient
}

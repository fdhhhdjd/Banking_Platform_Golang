package db

import (
	"log"

	"github.com/fdhhhdjd/Banking_Platform_Golang/worker"
	"github.com/hibiken/asynq"
)

// Config holds the application configuration
type Config struct {
	Cache CacheConfig
}

// CacheConfig holds the cache configuration
type CacheConfig struct {
	Link string
}

// Global variables
var AppConfig Config
var TaskDistributor worker.TaskDistributor

func initTaskDistributor() {
	// Cache connection setup
	redisOpt := asynq.RedisClientOpt{
		Addr: AppConfig.Cache.Link,
	}

	// Initialize task distributor
	TaskDistributor = worker.NewRedisTaskDistributor(redisOpt)

	log.Println("Task distributor initialized with Redis:", AppConfig.Cache.Link)
}

// Initialize function to be called in the main file
func Initialize() {
	initTaskDistributor()
}

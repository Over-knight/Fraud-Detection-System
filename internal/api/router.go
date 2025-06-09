package api

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// NewRouter initializes a new Gin router with the given Redis client.
func NewRouter(rdb *redis.Client) *gin.Engine {
	r := gin.Default()
	r.GET("/health", healthHandler)
	r.POST("/transactions", transactionHandler(rdb))
	return r
}
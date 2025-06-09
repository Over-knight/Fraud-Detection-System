package api

import (
	"context"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"frauddetector/kpg/model"
)

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func transactionHandler(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tx model.Transaction
		if err := c.ShouldBindJSON(&tx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// TODO: apply ML rules or stub logic
		verdict := "ok"
		if tx.Amount > 1000 {
			verdict = "fraud"
		}
		// Store the transaction in Redis
		_ = rdb.Set(context.Background(), tx.ID, verdict, 0).Err()
		c.JSON(http.StatusOK, gin.H{
			"id": tx.ID,
			"verdict":verdict})

		// c.JSON(http.StatusOK, gin.H{"message": "Transaction stored successfully"})
	}
}
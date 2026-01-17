package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gine_engine := gin.Default()
	gine_engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	gine_engine.Run()
}

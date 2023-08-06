package service

import (
	"raychat/chat"
	"raychat/middlewares"
	"raychat/service/models"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()
	v1 := r.Group("/v1", middlewares.Auth)
	{
		v1.GET("/models", models.GetModelsEndpoint)
		v1.POST("/chat/completions", chat.ChatEndpoint)
		v1.OPTIONS("/chat/completions", OptionsHandler)
	}
	r.Run(":8080")
}

func OptionsHandler(c *gin.Context) {
	// Set headers for CORS
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Headers", "*")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

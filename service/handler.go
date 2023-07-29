package service

import (
	"raychat/chat"
	"raychat/service/models"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()
	v1 := r.Group("/v1")
	v1.GET("/models", models.GetModelsEndpoint)
	v1.POST("/chat/completions", chat.ChatEndpoint)
	r.Run(":8080")
}

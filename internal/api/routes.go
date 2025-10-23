package api

import (
	"chat-app-backend/internal/controllers"
	"chat-app-backend/internal/ws"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// websocket routing 
	hub1 := ws.NewHub()
	go hub1.Run()

	// auth routing
	router.POST("/signup", controllers.SignUp)
	router.POST("/login", controllers.Login)

	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWS(hub1, c.Writer, c.Request)
	})
}
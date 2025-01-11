package main

import (
	"github.com/Frisbon/hungrymonke/service/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	r := gin.Default() // creo un router Gin per gestire l'HTTP

	r.Static("/doc", "./doc")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/doc/api.yaml")))

	// nota: puoi fare una fuzione per tutti i get e dentro la funzione separare per parametro passato nel body c
	r.GET("/admin/listUsers", func(c *gin.Context) {
		handlers.ListUsers(c)
	})

	r.POST("/login", func(c *gin.Context) { handlers.Login(c) })

	r.PUT("/users/me/username", func(c *gin.Context) {
		handlers.SetMyUsername(c)
	})

	r.PUT("/users/me/photo", func(c *gin.Context) {
		handlers.SetMyPhoto(c)
	})

	r.GET("/conversations", func(c *gin.Context) {
		handlers.GetMyConversations(c)
	})

	r.GET("/conversations/:ID", func(c *gin.Context) {
		handlers.GetMyConversation(c)
	})

	r.POST("/conversations/messages", func(c *gin.Context) {
		handlers.SendMessage(c)
	})

	r.POST("/messages/:ID/forward", func(c *gin.Context) {
		handlers.ForwardMSG(c)
	})

	r.GET("/debug", func(c *gin.Context) {
		handlers.DebugPrintDatabases()
	})

	r.Run(":8080")

}

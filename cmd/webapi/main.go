package main

import (
	"github.com/Frisbon/hungrymonke/service/api/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

/*

	TODO:
		- a quanto pare posso mandare un messaggio in un gruppo anche se non ne faccio parte ma se so l'ID...

*/

func main() {

	r := gin.Default() // creo un router Gin per gestire l'HTTP

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081", "http://localhost:8082"}, // Permetti richieste dal frontend e swagger
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true, // Se usi i cookie/token di autenticazione
	}))

	r.Static("/doc", "./doc")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/doc/api.yaml")))

	r.GET("api/utils/listUsers", func(c *gin.Context) {
		handlers.ListUsers(c)
	})

	r.POST("api/login", func(c *gin.Context) { handlers.Login(c) })

	r.PUT("/api/users/me/username", func(c *gin.Context) {
		handlers.SetMyUsername(c)
	})

	r.PUT("/api/users/me/photo", func(c *gin.Context) {
		handlers.SetMyPhoto(c)
	})

	r.GET("/api/conversations", func(c *gin.Context) {
		handlers.GetMyConversations(c)
	})

	r.GET("/api/conversations/:ID", func(c *gin.Context) {
		handlers.GetMyConversation(c)
	})

	r.POST("/api/conversations/messages", func(c *gin.Context) {
		handlers.SendMessage(c)
	})

	r.POST("/api/messages/:ID/forward", func(c *gin.Context) {
		handlers.ForwardMSG(c)
	})

	r.POST("/api/messages/:ID/comments", func(c *gin.Context) {
		handlers.CommentMSG(c)
	})

	r.DELETE("/api/messages/:ID/comments", func(c *gin.Context) {
		handlers.UncommentMSG(c)
	})

	r.DELETE("/api/messages/:ID", func(c *gin.Context) {
		handlers.DeleteMSG(c)
	})

	r.PUT("/api/groups/members", func(c *gin.Context) {
		handlers.AddToGroup(c)
	})

	r.DELETE("/api/groups/members", func(c *gin.Context) {
		handlers.LeaveGroup(c)
	})

	r.PUT("/api/groups/:ID/name", func(c *gin.Context) {
		handlers.SetGroupName(c)
	})

	r.PUT("/api/groups/:ID/photo", func(c *gin.Context) {
		handlers.SetGroupPhoto(c)
	})

	r.GET("/debug", func(c *gin.Context) {
		handlers.DebugPrintDatabases()
	})

	r.GET("/api/utils/getconvoinfo/:ID", func(c *gin.Context) {
		handlers.GetConvoInfo(c)
	})

	r.POST("/api/utils/createConvo", func(c *gin.Context) {
		handlers.CreatePrivateConvo(c)
	})

	r.POST("/api/utils/createGroup", func(c *gin.Context) {
		handlers.CreateGroupConvo(c)
	})

	r.Run(":8082")

}

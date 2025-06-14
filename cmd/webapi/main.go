package main

import (
	"log"

	"github.com/Frisbon/hungrymonke/service/api/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	r := gin.Default() //  creo un router Gin per gestire l'HTTP

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:8081", "http://localhost:8082"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// This is the most important line for this fix.
		// It explicitly allows the browser to send these headers.
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
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

	if err := r.Run(":8082"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}

/*

LINTER COMMANDS:
	go (in root) -> golangci-lint run
	js (webui) -> yarn lint
	api (in root) -> spectral lint ./doc/api.yaml

*/

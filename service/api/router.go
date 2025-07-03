package api

import (
	"net/http"

	"github.com/Frisbon/hungrymonke/service/api/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	router *gin.Engine
}

func New() *Router {
	r := gin.Default()
	return &Router{router: r}
}

func (rt *Router) Handler() http.Handler {
	r := rt.router

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Static("/doc", "./doc")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/doc/api.yaml")))

	// --- Routes updated to match api.yaml ---

	r.POST("/api/login", handlers.Login)

	// User routes
	r.PUT("/api/users/me/username", handlers.SetMyUsername)
	r.PUT("/api/users/me/photo", handlers.SetMyPhoto)

	// Conversation routes
	r.GET("/api/conversations", handlers.GetMyConversations)
	r.GET("/api/conversations/:ID", handlers.GetConversation)
	r.POST("/api/conversations/messages", handlers.SendMessage)

	// Message routes
	r.POST("/api/messages/:ID/forward", handlers.ForwardMSG)
	r.POST("/api/messages/:ID/comments", handlers.CommentMSG)
	r.DELETE("/api/messages/:ID/comments", handlers.UncommentMSG)
	r.DELETE("/api/messages/:ID", handlers.DeleteMSG)

	// Group routes
	r.PUT("/api/groups/members", handlers.AddToGroup)
	r.DELETE("/api/groups/members", handlers.LeaveGroup)
	r.PUT("/api/groups/:ID/name", handlers.SetGroupName)
	r.PUT("/api/groups/:ID/photo", handlers.SetGroupPhoto)

	// Admin and Utility routes
	r.GET("/admin/listUsers", handlers.ListUsers)
	r.GET("/api/utils/getconvoinfo/:ID", handlers.GetConvoInfo)
	r.POST("/api/utils/createConvo", handlers.CreatePrivateConvo)
	r.POST("/api/utils/createGroup", handlers.CreateGroupConvo)

	// Debug route
	r.GET("/debug", func(c *gin.Context) {
		handlers.DebugPrintDatabases()
	})

	return r
}

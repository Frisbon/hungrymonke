package main

import (
	"github.com/Frisbon/hungrymonke/service/handlers"
	"github.com/Frisbon/hungrymonke/service/structures"
	"github.com/gin-gonic/gin"
)

var users = make(map[string]structures.User)

func main() {
	r := gin.Default() // creo un router Gin per gestire l'HTTP

	r.POST("/users", func(c *gin.Context) {
		handlers.CreateUser(c, users)
	})

	r.Run(":8080")
}

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

	r.Run(":8080")

}

/*

Collaudo CreateUser handler

COMANDO UNICO PER CREARE DUE UTENTI

	curl -X POST http://localhost:8080/admin \
		-H "Content-Type: application/json" \
		-d '{
			"username": {"userid": "Primo1"},
			"photo": {"photofile": []}
			}' && \

	curl -X POST http://localhost:8080/admin \
		-H "Content-Type: application/json" \
		-d '{
			"username": {"userid": "Secondo2"},
			"photo": {"photofile": []}
			}'

e delle conversazioni tra loro....



COMANDO PER...

login
	curl -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{"username": {"userid": "Secondo2"}}'

setmyusername
	curl -X POST http://localhost:8080/users/me/username \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer <token>" \
	-d "francesco totti"

setmyphoto
	curl -X POST http://localhost:8080/users/me/photo \
    -H "Authorization: Bearer <token>" \
    -F "file=@service/pictureslol/2x2 pink pixel.jpg"

getmyconvos
	curl -X GET http://localhost:8080/conversations \
	-H "Authorization: Bearer <token>" \

getmyconvo


sendmessage


	SENZA ID

	curl -X POST http://localhost:8080/conversations//messages \
	-H "Authorization: Bearer <your_token>" \
	-H "Content-Type: application/json" \
	-d '{
	"recipientusername": "Secondo2",
	"message": {
		"text": "Hello!",
		"timestamp": "2025-01-06T00:00:00Z"
	}
	}'

	CON ID

	curl -X POST http://localhost:8080/conversations//messages \
	-H "Authorization: Bearer <your_token>" \
	-H "Content-Type: application/json" \
	-d '{
	"recipientusername": "Secondo2",
	"message": {
		"text": "Hello!",
		"timestamp": "2025-01-06T00:00:00Z"
	}
	}'

*/

// TODO: Funzioni per debug e per visualizzare correttamente liste di messaggi conversazioni ecc.

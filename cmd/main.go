package main

import (
	"github.com/Frisbon/hungrymonke/service/handlers"
	"github.com/Frisbon/hungrymonke/service/structures"
	"github.com/gin-gonic/gin"
)

var UserDB = make(map[string]structures.User)                   // nickname : struct utente.
var ConversationsDB = make(map[string]structures.Conversations) //  nickname : lista conversazioni di quell'utente.

var PrivateDB []structures.Private // Database con tutte le chat "1v1"

func main() {

	r := gin.Default() // creo un router Gin per gestire l'HTTP

	// LISTA DI HANDLER CON I PATH

	r.POST("/admin", func(c *gin.Context) {
		handlers.CreateUser(c, UserDB)
	})

	// nota: puoi fare una fuzione per tutti i get e dentro la funzione separare per parametro passato nel body c
	r.GET("/admin", func(c *gin.Context) {
		handlers.ListUsers(c, UserDB)
	})

	r.POST("/login", func(c *gin.Context) { handlers.Login(c, UserDB) })

	r.POST("/users/me/username", func(c *gin.Context) {
		handlers.SetMyUsername(c, UserDB)
	})

	r.POST("/users/me/photo", func(c *gin.Context) {
		handlers.SetMyPhoto(c, UserDB)
	})

	r.GET("/conversations", func(c *gin.Context) {
		handlers.GetMyConversations(c, UserDB, ConversationsDB)
	})

	r.GET("/conversations/:ID", func(c *gin.Context) {
		handlers.GetMyConversation(c, UserDB, ConversationsDB)
	})

	r.POST("/conversations/:ID/messages", func(c *gin.Context) {
		handlers.SendMessage(c, UserDB, ConversationsDB, PrivateDB)
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
	-H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzYyMjAwMzEsInN1YiI6Im5pbGwga2lnZ2VycyJ9.yPGIaj_Vh7ituvZwb1TFzd5RpfttirG324rMnwc4ADQ" \


eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzYyMTk3NTQsInN1YiI6IlByaW1vMSJ9.cLZSJLZKHN896aC0fpFHKxmto8jCFWdRlshgKh8NsnY

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

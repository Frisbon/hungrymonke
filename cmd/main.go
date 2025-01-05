package main

import (
	"github.com/Frisbon/hungrymonke/service/handlers"
	"github.com/Frisbon/hungrymonke/service/structures"
	"github.com/gin-gonic/gin"
)

var UserDB = make(map[string]structures.User)

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

COMANDO PER...


login
	curl -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{"username": {"userid": "Terzo3"}}'

setmyusername
	curl -X POST http://localhost:8080/users/me/username \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer <token>" \
	-d '{"newUsername": "francesco totti"}'


*/

package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Frisbon/hungrymonke/service/structures"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

/* ALCUNE FUNZIONI NON SONO DOCUMENTATE POICHÈ FATTE PER PRATICA, SONO NEL PATH /admin */

/*
POST, path /admin (usato per il debug)
- c contiene tutte le info sull'HTTP (body, parametri e metodi)
- UserDB è una mappa Nickname:User
*/
func CreateUser(c *gin.Context, UserDB map[string]structures.User) {

	var newUser structures.User
	// vedo se il JSON ricevuto si binda correttamente alla struct del User dichiarata.
	if err := c.ShouldBindJSON(&newUser); err != nil {
		//se errore rispondo con 400 bad req.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, exists := UserDB[newUser.Username.UserID]; !exists {
		UserDB[newUser.Username.UserID] = newUser
		c.JSON(http.StatusCreated, gin.H{"message": "Utente creato.", "user": newUser})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Utente esiste già, comando ignorato.", "user": newUser})
	}

}

// GET, path /admin
func ListUsers(c *gin.Context, UserDB map[string]structures.User) {
	fmt.Println("||----- USER LIST -----||")
	if len(UserDB) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Non trovo Utenti..."})
	}

	for k := range UserDB {
		fmt.Println(k)
	}

}

// POST, path /users/me/username
func SetMyUsername(c *gin.Context, UserDB map[string]structures.User) {

	//siccome lavoro con il current user, estraggo il token e leggo il nome dal claim
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token non trovato, te sei loggato ve?"})
		return
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims := &jwt.StandardClaims{}

	/*
		ParseWithClaims:
			-prende tokenstring, lo decodifica, mette nella struct claims
			- ci sarebbe anche la convalida del token tramite il func sotto:
				-prendo il token in input e restituisco la jwtkey se presente o errore se qualcosa va storto.
	*/
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) { return jwtKey, nil })

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token non valido"})
		return
	}

	// ora posso usare claims
	username := claims.Subject

	user, exists := UserDB[username]
	if !exists { // se non esiste utente dai claims
		c.JSON(http.StatusNotFound, gin.H{"error": "Utente non trovato"})
		return
	}

	// aggiurnamento username con quello che si trova nel body (c)
	var newUsername structures.UserID
	if err := c.ShouldBindJSON(&newUsername); err != nil { // bindo al json
		c.JSON(http.StatusBadRequest, gin.H{"error": "input non valido, come fai a NON darmi una stringa???"})
		return
	}

	user.Username = newUsername
	UserDB[username] = user

	c.JSON(http.StatusOK, gin.H{"message": "Username aggiornato lessgo", "user": user})
}

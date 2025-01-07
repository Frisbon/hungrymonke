package handlers

// login handlers + functions

import (
	"net/http"
	"time"

	"github.com/Frisbon/hungrymonke/service/structures"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("non_so_perchè_complico_cosi_tanto_questa_roba...")

/* NB Devi passare il nome utente come stringa*/
func GeneraToken(username string) (string, error) {

	expirationTime := time.Now().Add(12 * time.Hour) // token scade 12 ore

	claims := &jwt.StandardClaims{
		Subject:   username,              // Username o ID dell'utente
		ExpiresAt: expirationTime.Unix(), // tempo di scadenza del token
	}

	// creo token con i parametri sopra
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey) //firma e restituisci stringa token

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// handler per il login
func Login(c *gin.Context, UserDB map[string]structures.User) {
	var reqUserID string

	//creo una struttura user e vedo se nel body mi arriva quello.
	if err := c.ShouldBind(&reqUserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "il dato del curl non funge fra non è stringa"}) // se non funge
		return
	}

	// vedo se esiste user nel sistema, se non esiste lo creo
	if _, exists := UserDB[reqUserID]; !exists {
		UserDB[reqUserID] = structures.User{Username: reqUserID}
	}

	//genero token x user e lo returno
	tokenString, err := GeneraToken(reqUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "non riesco a crea il token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "user": UserDB[reqUserID]})
}

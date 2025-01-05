package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/Frisbon/hungrymonke/service/structures"
)

/* ALCUNE FUNZIONI NON SONO DOCUMENTATE POICHÈ FATTE PER PRATICA*/

/*
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
	if _, exists := UserDB[newUser.Username.UserID]; exists == false {
		UserDB[newUser.Username.UserID] = newUser
		c.JSON(http.StatusCreated, gin.H{"message": "Utente creato.", "user": newUser})
	}else 
	{
		c.JSON(http.StatusBadRequest, gin.H{"message": "Utente esiste già, comando ignorato.", "user": newUser})
	}
		
}

func GetUser(c *gin.Context) {
    id := c.Param("id")
    // Logica per recuperare l'utente dall'archivio dati.
    user, found := /* cerca l'utente usando l'id */
    if !found {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}

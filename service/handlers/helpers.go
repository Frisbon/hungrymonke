package handlers

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	scs "github.com/Frisbon/hungrymonke/service/structures"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// helper che ritorna (nome utente, struct utente, errore). Utile per funzioni dove sono il current_user.
func QuickAuth(c *gin.Context) (string, *scs.User, string) {

	//siccome lavoro con il current user, estraggo il token e leggo il nome dal claim
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token non trovato, te sei loggato ve?"})
		return "", nil, "err"
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
		return "", nil, "err"
	}

	// ora posso usare claims
	username := claims.Subject // username è il nome del current user

	user, exists := scs.UserDB[username]
	if !exists { // se non esiste utente dai claims
		c.JSON(http.StatusNotFound, gin.H{"error": "Utente non trovato"})
		return "", nil, "err"
	}

	return username, user, ""

}

// NB: Genera SEMPRE un ID univoco. Aggiorna anche GenericDB
func GenerateRandomString(length int) string {
	for {

		letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

		randGen := rand.New(rand.NewSource(time.Now().UnixNano())) // Nuovo generatore di numeri casuali
		var sb strings.Builder
		for i := 0; i < length; i++ {
			sb.WriteByte(letters[randGen.Intn(len(letters))]) // Sceglie un carattere casuale dal set
		}

		if isUniversalIdUnique(sb.String()) {
			scs.GenericDB[sb.String()] = struct{}{}
			return sb.String()
		}
	}

}

// gli passo tutti i DB che usano generic ID e controllo se si ha uno li dentro
func isUniversalIdUnique(ID string) bool {

	if _, exists := scs.GenericDB[ID]; exists {
		return false
	}
	return true

}

// costruisce la stringa di preview in base al contenuto passatogli come parametro
func PreviewMaker(Content scs.Content) string {

	if Content.Text != nil && Content.Photo != nil {

		if len(*Content.Text) > 8 {
			return "📷 " + (*Content.Text)[:8] + "..."
		}
		// uso tutta la stringa se è più corta di 8 caratteri
		return "📷 " + *Content.Text

	} else if Content.Text != nil {
		if len(*Content.Text) > 8 {
			return (*Content.Text)[:8] + "..."
		}
		return *Content.Text

	} else if Content.Photo != nil {
		return "📷 Photo..." // Ritorna una stringa indicante la presenza di un'immagine

	} else {
		return "ERROR - INVALID CONTENT" // Nel caso in cui non ci sia né testo né foto
	}

}

/*
Aggiorna la convo con le info dell'ultimo messaggio, ossia la preview e la last msg timestamp
*/
func UpdateConversationWLastMSG(convo *scs.ConversationELT) {

	// recupero ultimo messaggio
	msglist := convo.Messages
	lst_msg := msglist[len(msglist)-1]

	convo.DateLastMessage = lst_msg.Timestamp
	convo.Preview = PreviewMaker(lst_msg.Content)

}

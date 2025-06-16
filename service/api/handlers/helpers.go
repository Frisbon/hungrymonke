package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	scs "github.com/Frisbon/hungrymonke/service/api/structures"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// helper che ritorna (nome utente, struct utente, errore). Utile per funzioni dove sono il current_user.
// NB devo gestire i lock con le funzioni esterne a quickauth.
func QuickAuth(c *gin.Context) (string, *scs.User, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token non trovato, te sei loggato ve?"})
		return "", nil, errors.New("token not found")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) { return jwtKey, nil })

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token non valido"})
		return "", nil, errors.New("invalid token")
	}

	username := strings.Trim(claims.Subject, "\"")

	user, exists := scs.UserDB[username]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utente non trovato"})
		return "", nil, errors.New("user not found")
	}

	return username, user, nil
}

// NB: Genera SEMPRE un ID univoco. Aggiorna anche GenericDB
func GenerateRandomString(length int) string {
	for {
		letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
		var sb strings.Builder
		for i := 0; i < length; i++ {
			sb.WriteByte(letters[randGen.Intn(len(letters))])
		}

		if isUniversalIdUnique(sb.String()) {
			scs.GenericDB[sb.String()] = struct{}{}
			return sb.String()
		}
	}
}

// gli passo tutti i DB che usano generic ID e controllo se si ha uno li dentro
func isUniversalIdUnique(id string) bool {
	_, exists := scs.GenericDB[id]
	return !exists
}

// costruisce la stringa di preview in base al contenuto passatogli come parametro
func PreviewMaker(Content scs.Content) string {

	if Content.Text != nil && Content.Photo != nil {

		if len(*Content.Text) > 16 {
			return "📷 " + (*Content.Text)[:16] + "..."
		}
		// uso tutta la stringa se è più corta di 16 caratteri
		return "📷 " + *Content.Text

	} else if Content.Text != nil {
		if len(*Content.Text) > 16 {
			return (*Content.Text)[:16] + "..."
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
	if len(convo.Messages) > 0 {
		lst_msg := convo.Messages[len(convo.Messages)-1]
		convo.DateLastMessage = lst_msg.Timestamp
		convo.Preview = PreviewMaker(lst_msg.Content)
	}
}

/*
Aggiorna i messaggi "delivered" come "seen" dell'utente opposto.
Basically implies che ho visualizzato i suoi messaggi
*/
func PrivateMsgStatusUpdater(convo *scs.ConversationELT, logged_user *scs.User) {
	for _, msg := range convo.Messages {
		if msg.Author.Username != logged_user.Username && msg.Status != scs.Seen {
			msg.Status = scs.Seen
		}
	}
}

func ContainsUser(users []*scs.User, target *scs.User) bool {
	for _, user := range users {
		if user.Username == target.Username {
			return true
		}
	}
	return false
}

/*
Aggiorna i messaggi (x)"delivered" come "seen"
<=>
Ogni utente del gruppo ha inviato un messaggio dopo x oppure ha visualizzato la chat
HINT: Usa mappa boolean [group_user: seen (bool)]
*/
func GroupMsgStatusUpdater(group_convo *scs.ConversationELT, logged_user *scs.User) {
	for _, msg := range group_convo.Messages {
		if !ContainsUser(msg.SeenBy, logged_user) {
			msg.SeenBy = append(msg.SeenBy, logged_user)
		}

		seen := true
		if msg.Status != scs.Seen {
			for _, convoUser := range scs.GroupDB[group_convo.ConvoID].Users {
				if !ContainsUser(msg.SeenBy, convoUser) {
					seen = false
					break
				}
			}
		}

		if msg.Author.Username != logged_user.Username && seen {
			msg.Status = scs.Seen
		}
	}
}

/*returns true if error!!!!*/
func statusUpdater(convo *scs.ConversationELT, current_user *scs.User) bool {
	if _, exists := scs.PrivateDB[convo.ConvoID]; exists {
		PrivateMsgStatusUpdater(convo, current_user)
		return false
	} else if _, exists := scs.GroupDB[convo.ConvoID]; exists {
		GroupMsgStatusUpdater(convo, current_user)
		return false
	} else {
		return true
	}
}

func DebugPrintDatabases() {
	scs.DBMutex.RLock()         //*
	defer scs.DBMutex.RUnlock() //*

	printMap := func(name string, data interface{}) {
		if len(fmt.Sprintf("%v", data)) == 0 {
			fmt.Printf("%s è vuoto.\n", name)
			return
		}
		fmt.Printf("Contenuto di %s:\n", name)
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Printf("Errore durante la conversione in JSON di %s: %v\n", name, err)
			return
		}
		fmt.Println(string(jsonData))
	}

	printMap("GenericDB", scs.GenericDB)
	printMap("UserDB", scs.UserDB)
	printMap("PrivateDB", scs.PrivateDB)
	printMap("GroupDB", scs.GroupDB)
	printMap("MsgDB", scs.MsgDB)
	printMap("ConvoDB", scs.ConvoDB)
	printMap("UserConvosDB", scs.UserConvosDB)
}

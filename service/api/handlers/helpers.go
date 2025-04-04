package handlers

import (
	"encoding/json"
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
func QuickAuth(c *gin.Context) (string, *scs.User, string) {

	//siccome lavoro con il current user, estraggo il token e leggo il nome dal claim
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token non trovato, te sei loggato ve?"})
		return "", nil, "err"
	}

	fmt.Println("Token ricevuto:", tokenString)

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
	// Rimuove eventuali virgolette extra dal claim Subject
	username := strings.Trim(claims.Subject, "\"") // username è il nome del current user

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

	// recupero ultimo messaggio
	msglist := convo.Messages
	lst_msg := msglist[len(msglist)-1]

	convo.DateLastMessage = lst_msg.Timestamp
	convo.Preview = PreviewMaker(lst_msg.Content)

}

/*
Aggiorna i messaggi "delivered" come "seen" dell'utente opposto.
Basically implies che ho visualizzato i suoi messaggi
*/
func PrivateMsgStatusUpdater(convo *scs.ConversationELT, logged_user *scs.User) {

	for _, msg := range convo.Messages {

		if msg.Author != logged_user && msg.Status != "seen" {
			msg.Status = "seen"
		}

	}

}

func ContainsUser(users []*scs.User, target *scs.User) bool {
	for _, user := range users {
		if user == target {
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

		// prendo un messaggio

		// aggiungo il mio "seen"

		// se il seen array ha tutti gli utenti allora status = seen
		if !ContainsUser(msg.SeenBy, logged_user) {
			msg.SeenBy = append(msg.SeenBy, logged_user)
		}

		// se non è seen comincio a controllare se è stato visto da tutti gli utenti
		seen := true
		if msg.Status != "seen" { // Compare with the Seen constant, not the string "seen"
			for _, convoUser := range scs.GroupDB[group_convo.ConvoID].Users { // Added 'range' keyword
				if !ContainsUser(msg.SeenBy, convoUser) {
					seen = false
					break // Optional: break early since we already know it's not seen by all
				}
			}
		}

		if msg.Author != logged_user && seen {
			msg.Status = "seen"
		}

	}
	// TODO, clearup inconsistency between private "seen" and group one.
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
	// Funzione helper per stampare le mappe in formato JSON leggibile
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

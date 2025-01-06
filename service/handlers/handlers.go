package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Frisbon/hungrymonke/service/structures"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

/*
Tutti i Helper si troveranno qui.

*/

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

	for _, v := range UserDB {
		fmt.Println(v)
	}

}

// helper che ritorna (nome utente, struct utente, errore)
func quickAuth(c *gin.Context, UserDB map[string]structures.User) (string, structures.User, string) {

	//siccome lavoro con il current user, estraggo il token e leggo il nome dal claim
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token non trovato, te sei loggato ve?"})
		return "", structures.User{}, "err"
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
		return "", structures.User{}, "err"
	}

	// ora posso usare claims
	username := claims.Subject // username è il nome del current user

	user, exists := UserDB[username]
	if !exists { // se non esiste utente dai claims
		c.JSON(http.StatusNotFound, gin.H{"error": "Utente non trovato"})
		return "", structures.User{}, "err"
	}

	return username, user, ""

}

// POST, path /users/me/username
func SetMyUsername(c *gin.Context, UserDB map[string]structures.User) {

	// AUTENTICAZIONE UTENTE
	username, user, er1 := quickAuth(c, UserDB)
	if len(er1) != 0 {
		return
	}

	// aggiurnamento username con quello che si trova nel body "text/plain" (c)
	newUsername, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nun riesc a legg u body"})
		return
	}

	usernameString := string(newUsername)
	if usernameString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input non valido, stringa vuota"})
		return
	}

	delete(UserDB, username) // elimino il current user dal db
	user.Username.UserID = usernameString
	UserDB[usernameString] = user // lo riaggiungo con il nome diverso.

	//azzo devo anche aggiornarlo con il token altrimenti si ricorda ancora il nome precedente
	// rigeneriamo il token con il nuovo username
	tokenString, err := GeneraToken(usernameString) // Supponiamo che la funzione generaToken prenda il nuovo username
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Non riesco a creare il token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Username aggiornato lessgo", "user": user, "new_token": tokenString})

}

// POST, path /users/me/photo
func SetMyPhoto(c *gin.Context, UserDB map[string]structures.User) {

	// AUTENTICAZIONE UTENTE
	username, user, er1 := quickAuth(c, UserDB)
	if len(er1) != 0 {
		return
	}

	// aggiurnamento foto con quella che si trova nel body (c)
	// estraggo il file dalla richiesta multipart/form-data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File invalido :P"})
		return
	}

	// apro il file appena estratto
	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "non riesco ad apri, ao me apri"})
		return
	}
	defer openedFile.Close()

	// Leggo il file in modo da salvarlo come []byte dentro al content. Lo faccio perchè PhotoFile è di tipo []byte
	content, err := io.ReadAll(openedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "che hai scritto nel file?? non leggo..."})
		return
	}

	// Aggiorna l'utente con il percorso del file salvato
	user.Photo.PhotoFile = content
	UserDB[username] = user

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"message": "Foto aggiornata lessgo",
	})
}

// GET, path /conversations
// TODO - COLLAUDO CON CONVERSAZIONI GIÀ CREATE.
func GetMyConversations(c *gin.Context, UserDB map[string]structures.User, ConversationsDB map[string]structures.Conversations) {

	// AUTENTICAZIONE UTENTE
	username, _, er1 := quickAuth(c, UserDB)
	if len(er1) != 0 {
		return
	}

	fmt.Println("||----- CURRENT USER CONVOs LIST -----||")

	// estraggo da conversationsDB[username]
	fmt.Println(ConversationsDB[username])

}

// GET, path /conversations/:ID
// TODO - COLLAUDO CON CONVERSAZIONI GIÀ CREATE.
func GetMyConversation(c *gin.Context, UserDB map[string]structures.User, ConversationsDB map[string]structures.Conversations) structures.ConversationELT {

	// AUTENTICAZIONE UTENTE
	username, _, er1 := quickAuth(c, UserDB)
	if len(er1) != 0 {
		return structures.ConversationELT{}
	}

	//recupero le conversazioni di quell'utente
	user_conversations, exists := ConversationsDB[username]
	if !exists {
		c.JSON(http.StatusNoContent, gin.H{"error": "Non ho conversazioni per questo utente."})
		return structures.ConversationELT{}
	}

	//estraggo l'ID dal body
	conversationID := c.Param("ID")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID della conversazione mancante"})
		return structures.ConversationELT{}
	}

	//ricerco la conversazione nel ConversationsDB tramite l'id estratto.
	var found_conversation structures.ConversationELT
	for _, conversation := range user_conversations {
		if conversation.ID.GenericID == conversationID {
			found_conversation = conversation
			break
		}
	}
	if len(found_conversation.ID.GenericID) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversazione non trovata"})
		return structures.ConversationELT{}
	}

	// Risposta con la conversazione trovata
	c.JSON(http.StatusOK, gin.H{
		"message":      "Conversazione trovata",
		"conversation": found_conversation,
	})

	fmt.Printf("||----- Conversation ID: %s -----||\n", conversationID)
	fmt.Println(found_conversation)
	return found_conversation

}

// helper ( x SendMessage) per generare un ID stringa casuale
func generateRandomString(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randGen := rand.New(rand.NewSource(time.Now().UnixNano())) // Nuovo generatore di numeri casuali
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(letters[randGen.Intn(len(letters))]) // Sceglie un carattere casuale dal set
	}
	return sb.String()
}

func SendMessage(c *gin.Context, UserDB map[string]structures.User, ConversationsDB map[string]structures.Conversations, PrivateDB []structures.Private) {

	//autorizzo il current user
	SenderUsername, _, er1 := quickAuth(c, UserDB)
	if len(er1) != 0 {
		return
	}

	// leggo il messaggio e il username (se c'è) dal body

	type requestLol struct {
		RecipientUsername string             `json:"recipientusername"`
		Message           structures.Message `json:"message"`
	}

	var req requestLol
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON non valido", "details": err.Error()})
		return
	}

	//estraggo l'ID dal body (se c'è)
	conversationID := c.DefaultQuery("ID", "") // se c'è ID lo estrae altrimenti lo mette come ""

	if conversationID != "" { // se tra parametri ho ID invio su ID direttamente

		convo := GetMyConversation(c, UserDB, ConversationsDB)

		// se non esiste convo do errore
		if len(convo.ID.GenericID) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Convo non trovata con l'ID che mi hai fornito..."})
			return
		}

		// se esiste invio messaggio su quella convo e la salvo sul DB
		convo.Messages = append(convo.Messages, req.Message)
		for nome, convolist := range ConversationsDB {
			for i, convo2 := range convolist {
				if convo.ID.GenericID == convo2.ID.GenericID {
					// ho trovato la convo nel db
					convolist[i] = convo              // update Conversations
					ConversationsDB[nome] = convolist // update ConversationsDB
					fmt.Println("Convo salvata nel DB:", convo)
					return
				}
			}
		}

	} else if conversationID == "" && req.RecipientUsername != "" { // se context NON ha ID ma ha username

		// cerco convo privata tra i due username
		var found_private structures.Private
		for _, private := range PrivateDB {

			// se (Utente1 = Sender AND Utente2 = Reciever) OR Viceversa OR Sender == Reciever
			if private.FirstUser.Username.UserID == SenderUsername && private.SecondUser.Username.UserID == req.RecipientUsername {
				found_private = private
				break
			} else if private.FirstUser.Username.UserID == req.RecipientUsername && private.SecondUser.Username.UserID == SenderUsername {
				found_private = private
				break
			} else if private.FirstUser.Username.UserID == private.SecondUser.Username.UserID {
				found_private = private
				break
			}

		}

		// se non esiste
		if len(found_private.Conversation.ID.GenericID) == 0 {
			// creo convo (e anche privata).
			var newConvo structures.ConversationELT

			/*
					genero ID randomico di 5 caratteri (a caso) che non sia presente in ConversationsDB
				    ricerco la conversazione nel ConversationsDB tramite l'id generato.

					Simulo un DO-WHILE perchè sono maestro del Bydon
					Il MEGA for controlla se la stringa appena generata sia univoca,
					quindi itera tutta la mappa ConversationsDB e le Conversations (array)
					all'interno per vedere se ci sia una con ID combaciante.
					Se trova un ID già presente, lo rigenera, altrimenti esce dal ciclo.

					NB: possono esistere due convo con lo stesso ID nel ConversationsDB
					se A chatta con B
					DB[A] = id0			perchè id0 è la convo di A con B
					DB[B] = id0			perchè id0 è anche la convo di B con A
			*/
			var found bool
			var convoID string
			for {
				convoID = generateRandomString(5)
				found = true

				// controllo se l'ID è già nel DB
				for _, convolist := range ConversationsDB {
					for _, conversation := range convolist {
						if conversation.ID.GenericID == convoID {
							// se già cho quell'ID nel DB, rigenero la stringa
							found = false
							break
						}
					}
					if !found {
						break
					}
				}

				// Se l'ID non è stato trovato, esci dal ciclo
				if found {
					break
				}
			}

			newConvo.ID.GenericID = convoID                            //imposto id convo
			newConvo.Messages = append(newConvo.Messages, req.Message) // invio messaggio su quella convo

			ConversationsDB[SenderUsername] = append(ConversationsDB[SenderUsername], newConvo) // aggiorno ConversationsDB per entrambi i lati.
			ConversationsDB[req.RecipientUsername] = append(ConversationsDB[req.RecipientUsername], newConvo)

			var newPrivate structures.Private
			newPrivate.Conversation = newConvo
			newPrivate.FirstUser = UserDB[SenderUsername]
			newPrivate.SecondUser = UserDB[req.RecipientUsername]

			PrivateDB = append(PrivateDB, newPrivate)

		} else { // se esiste

			// invio messaggio su quella convo
			found_private.Conversation.Messages = append(found_private.Conversation.Messages, req.Message)

		}

		fmt.Println("Messaggio Inviato! :)")
		fmt.Println(PrivateDB)

		return
	}

}

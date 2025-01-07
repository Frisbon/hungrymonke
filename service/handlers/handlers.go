package handlers

import (
	"fmt"
	"io"
	"net/http"

	scs "github.com/Frisbon/hungrymonke/service/structures"
	"github.com/gin-gonic/gin"
)

// GET, path /admin
func ListUsers(c *gin.Context) {

	if len(scs.UserDB) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Non trovo Utenti..."})
	} else {

		var array []scs.User
		for _, v := range scs.UserDB {
			array = append(array, scs.User{Username: v.Username, Photo: v.Photo})
		}
		c.JSON(http.StatusOK, gin.H{"Users": array})

		for _, elt := range array {
			fmt.Printf("%+v\n", elt)

		}
	}

}

// POST, path /users/me/username
func SetMyUsername(c *gin.Context) {

	// AUTENTICAZIONE UTENTE
	username, user, er1 := QuickAuth(c)
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

	delete(scs.UserDB, username) // elimino il current user dal db
	user.Username = usernameString
	scs.UserDB[usernameString] = user // lo riaggiungo con il nome diverso.

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
func SetMyPhoto(c *gin.Context) {

	// AUTENTICAZIONE UTENTE
	username, user, er1 := QuickAuth(c)
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
	user.Photo = content
	scs.UserDB[username] = user

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"message": "Foto aggiornata lessgo",
	})
}

// GET, path /conversations
func GetMyConversations(c *gin.Context) {

	// AUTENTICAZIONE UTENTE
	username, _, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	if len(scs.UserConvosDB[username]) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Non ci sono conversazioni per questo utente", "Username": username})
	}

	c.JSON(http.StatusOK, gin.H{"Username": username, "User Conversations": scs.UserConvosDB[username]})

}

// GET, path /conversations/:ID
func GetMyConversation(c *gin.Context) {

	// AUTENTICAZIONE UTENTE
	username, _, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	//recupero le conversazioni di quell'utente visto che devo cercare l'id SOLO PER LE CONVERSAZIONI DI QUELL'UTENTE!
	user_conversations, exists := scs.UserConvosDB[username]
	if !exists {
		c.JSON(http.StatusNoContent, gin.H{"error": "Non ho conversazioni per questo utente."})
		return
	}

	//estraggo l'ID dal body
	conversationID := c.Param("ID")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID della conversazione mancante"})
		return
	}

	//ricerco la conversazione nel UserConvosDB tramite l'id estratto.
	var found_conversation scs.ConversationELT

	for _, conversation := range user_conversations {
		if conversation.ConvoID == conversationID {
			found_conversation = conversation
			break
		}
	}
	if len(found_conversation.ConvoID) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversazione non trovata"})
		return
	}

	// Risposta con la conversazione trovata
	c.JSON(http.StatusOK, gin.H{
		"message":      "Conversazione trovata",
		"conversation": found_conversation,
	})

}

// POST, path /conversations/messages
func SendMessage(c *gin.Context) {

	//autorizzo il current user
	SenderUsername, sender_struct, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	// leggo il messaggio e il username a chi inviarlo (se c'è) dal body
	type requestLol struct {
		RecipientUsername string      `json:"recipientusername"`
		Message           scs.Message `json:"message"`
	}

	var req requestLol
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON non valido", "details": err.Error()})
		return
	}

	msgc := MsgCONSTR{
		Timestamp: req.Message.Timestamp,
		Content:   req.Message.Content,
		Status:    req.Message.Status,
		Reactions: req.Message.Reactions,
		Author:    sender_struct,
	}

	newMessage := ConstrMessage(msgc)

	//estraggo l'ID dal body (se c'è)
	conversationID := c.DefaultQuery("ID", "")
	// se tra parametri ho ID
	if conversationID != "" {

		// prendo conversazione con l'ID che mi è stato fornito.
		convo := scs.ConvoDB[conversationID]
		//aggiungo il messaggio alla lista dei messaggi su quella convo
		convo.Messages = append(convo.Messages, newMessage)

		// se context NON ha ID ma ha username
	} else if conversationID == "" && req.RecipientUsername != "" {

		// cerco convo privata tra i due username
		var found_private *scs.Private
		for _, private := range scs.PrivateDB {

			// se (Utente1 = Sender AND Utente2 = Reciever) OR Viceversa OR Sender == Reciever
			if private.FirstUser.Username == SenderUsername && private.SecondUser.Username == req.RecipientUsername {
				found_private = private
				break
			} else if private.FirstUser.Username == req.RecipientUsername && private.SecondUser.Username == SenderUsername {
				found_private = private
				break
			} else if private.FirstUser.Username == private.SecondUser.Username {
				found_private = private
				break
			}

		}

		// se non esiste
		if len(found_private.Conversation.ConvoID) == 0 {

			//recupero puntatore al secondo utente

			//creo convo
			x := ConvoCONSTR{
				DateLastMessage: newMessage.Timestamp,
				Preview:         "",                         //TODO: implementa preview maker, qui lo devi prendere da newMessage.text
				Messages:        []*scs.Message{newMessage}, // array con solo il messaggio mandato al momento.
			}

			newConversation := ConstrConvo(x)

			p := &scs.Private{
				Conversation: newConversation,
				FirstUser:    sender_struct,
				SecondUser:   scs.UserDB[req.RecipientUsername],
			}

			scs.PrivateDB[newConversation.ConvoID] = p

			c.JSON(http.StatusOK, gin.H{"Status": "Non ho trovato una convo con questo utente... Perciò l'ho creata!", "Private_Conversation": &p})

		} else { // se esiste

			// invio messaggio su quella convo
			found_private.Conversation.Messages = append(found_private.Conversation.Messages, newMessage)
			c.JSON(http.StatusOK, gin.H{"Status": "Ho trovato una convo esistente con questo utente... Messaggio inviato!", "Conversation": found_private.Conversation})

		}

		fmt.Println("Messaggio Inviato! :)")

		return
	}

}

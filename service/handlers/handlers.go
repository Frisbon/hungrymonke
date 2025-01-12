package handlers

import (
	"fmt"
	"io"
	"net/http"
	"time"

	scs "github.com/Frisbon/hungrymonke/service/structures"
	"github.com/gin-gonic/gin"
)

// GET, path /admin/users
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

	// aggiornamento foto con quella che si trova nel body (c)
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
	username, logged_struct, er1 := QuickAuth(c)
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
	var found_conversation *scs.ConversationELT

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

	/*

		dopo aver trovato la conversazione, mi assicuro che lo status
		*dei messaggi "delivered" degli altri utenti* sia impostato a seen
		se è una convo 1v1

		in caso fosse una convo di gruppo
		OGNI utente deve visualizzare per poter mettere lo status seen.
		creo una mappa [user:boolean] e se ogni ogni boolean è true allora
		status seen (per chat di gruppo.)


	*/
	PrivateMsgStatusUpdater(found_conversation, logged_struct)

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
		Message           scs.Content `json:"message"`
	}

	var req requestLol
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON non valido", "details": err.Error()})
		return
	}

	var reactlist []scs.Reaction

	// DEFINISCO NUOVO MESSAGGIO CON PARAMETRI GIUSTI
	msgc := MsgCONSTR{
		Timestamp: time.Now(),
		Content:   req.Message,
		Status:    "delivered",
		Reactions: reactlist,
		Author:    sender_struct,
	}

	// USO COSTRUTTORE PER GENERARE ANCHE L'ID MESSAGGIO UNIVOCO
	newMessage := ConstrMessage(msgc)

	//estraggo l'ID dalla query (se c'è)
	conversationID := c.DefaultQuery("ID", "")

	// se mando sia ID che Username, do errore
	if conversationID != "" && req.RecipientUsername != "" {

		c.JSON(http.StatusBadRequest, gin.H{"Status": "Non puoi inviare entrambi ID e Nickname nel json! :O"})
		return
	}

	// se tra parametri ho ID e il nome è clear
	if conversationID != "" && req.RecipientUsername == "" {

		// prendo conversazione con l'ID che mi è stato fornito.
		convo := scs.ConvoDB[conversationID]
		//aggiungo il messaggio alla lista dei messaggi su quella convo
		convo.Messages = append(convo.Messages, newMessage)
		UpdateConversationWLastMSG(convo)

		// se mando un msg => leggo i messaggi dell'altro utente => aggiorno i loro status.
		PrivateMsgStatusUpdater(convo, sender_struct)

		// aggiorno le liste personali di entrambi gli utenti...
		scs.UserConvosDB[SenderUsername] = append(scs.UserConvosDB[SenderUsername], convo)
		scs.UserConvosDB[req.RecipientUsername] = append(scs.UserConvosDB[req.RecipientUsername], convo)

		c.JSON(http.StatusOK, gin.H{"Status": "Inviato sulla convo indicata dall'ID!", "Conversation": &convo})

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
			} else if private.FirstUser.Username == req.RecipientUsername && req.RecipientUsername == private.SecondUser.Username {
				found_private = private
				break
			}

		}

		// se non esiste
		if found_private == nil {

			//recupero puntatore al secondo utente

			//creo convo
			x := ConvoCONSTR{
				DateLastMessage: newMessage.Timestamp,
				Preview:         "",
				Messages:        []*scs.Message{newMessage}, // array con solo il messaggio mandato al momento.
			}

			newConversation := ConstrConvo(x)
			UpdateConversationWLastMSG(newConversation)

			// aggiorno le liste personali di entrambi gli utenti...
			scs.UserConvosDB[SenderUsername] = append(scs.UserConvosDB[SenderUsername], newConversation)
			scs.UserConvosDB[req.RecipientUsername] = append(scs.UserConvosDB[req.RecipientUsername], newConversation)

			p := &scs.Private{
				Conversation: newConversation,
				FirstUser:    sender_struct,
				SecondUser:   scs.UserDB[req.RecipientUsername],
			}

			scs.PrivateDB[newConversation.ConvoID] = p

			c.JSON(http.StatusOK, gin.H{"Status": "Non ho trovato una convo con questo utente... Perciò l'ho creata!", "Private_Conversation": &p})
			return

		} else { // se esiste

			found_private.Conversation.Messages = append(found_private.Conversation.Messages, newMessage)

			// se mando un msg => leggo i messaggi dell'altro utente => aggiorno i loro status.
			PrivateMsgStatusUpdater(found_private.Conversation, sender_struct)
			UpdateConversationWLastMSG(found_private.Conversation)
			c.JSON(http.StatusOK, gin.H{"Status": "Ho trovato una convo esistente con questo utente... Messaggio inviato!", "Conversation": found_private.Conversation})
			return

		}

	}

}

// POST, path /messages/{ID}/forward
func ForwardMSG(c *gin.Context) {
	// Note: Se inoltro non implica che leggo per forza il messaggio dell'altro.

	//autorizzo il current user
	_, sender_struct, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	// leggo l'ID messggio dal body
	type requestLol struct {
		ConvoID string `json:"ConvoID"`
	}

	var req requestLol
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON non valido", "details": err.Error()})
		return
	}

	//estraggo l'ID messaggio dal path (se c'è)
	MsgID := c.Param("ID")

	msg, exists := scs.MsgDB[MsgID]

	convo, exists2 := scs.ConvoDB[req.ConvoID]

	if MsgID == "" || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "MSG not found"})
		return
	}

	if !exists2 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	// Devo "clonare" il messaggio e aggiungerlo alla conversazione.
	var react_list []scs.Reaction
	newMsg := ConstrMessage(MsgCONSTR{
		Timestamp: time.Now(),
		Content:   msg.Content,
		Author:    sender_struct, //maschero da chi inoltro
		Status:    "delivered",
		Reactions: react_list, //resetto lista reazioni

	})

	convo.Messages = append(convo.Messages, newMsg)
	c.JSON(http.StatusOK, gin.H{"success": "Messaggio inoltrato", "MSG": newMsg})

	//bwt non faccio controlli delle fonti del messaggio ecc ecc perchè non mi va, basta che funge :)

}

// POST, path /messages/{ID}/comments
func CommentMSG(c *gin.Context) {

	//autorizzo il current user
	_, sender_struct, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	//estraggo l'ID messaggio dal path (se c'è)
	MsgID := c.Param("ID")

	if MsgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID vuoto"})
		return
	}

	// leggo la reaction dal body
	type requestLol struct {
		Emoticon string `json:"Emoticon"`
	}

	var req requestLol
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON non valido", "details": err.Error()})
		return
	}

	MsgStruct, exists := scs.MsgDB[MsgID]

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nun trovo u msg man"})
		return
	}

	//NB posso commentare max 1 volta....

	for _, R := range MsgStruct.Reactions {

		if R.Author == sender_struct {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bro u already commented, delete ur last reaction"})
			return
		}

	}

	MsgStruct.Reactions = append(MsgStruct.Reactions, scs.Reaction{
		Timestamp: time.Now(),
		Author:    sender_struct,
		Emoticon:  req.Emoticon,
	})

	c.JSON(http.StatusOK, gin.H{"Success": "Messaggio inoltrato", "MSG": MsgStruct})

}

// DELETE, path /messages/{ID}/comments
func UncommentMSG(c *gin.Context) {

	//autorizzo il current user
	_, sender_struct, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	//estraggo l'ID messaggio dal path (se c'è)
	MsgID := c.Param("ID")

	if MsgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID messaggio vuoto"})
		return
	}

	MsgStruct, exists := scs.MsgDB[MsgID]

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nun trovo u msg man"})
		return
	}

	//NB posso commentare max 1 volta....

	found := false
	toDelete := 0
	for i, R := range MsgStruct.Reactions {

		if R.Author == sender_struct {
			found = true
			toDelete = i
			break
		}

	}
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you left no reaction on that msg man"})
		return
	}

	MsgStruct.Reactions = append(MsgStruct.Reactions[:toDelete], MsgStruct.Reactions[toDelete+1:]...)

	c.JSON(http.StatusOK, gin.H{"Success": "Reaction deleted", "MSG": MsgStruct})

}

// DELETE, path /message/{ID}
func DeleteMSG(c *gin.Context) {

	//autorizzo il current user
	_, sender_struct, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	//estraggo l'ID messaggio dal path (se c'è)
	MsgID := c.Param("ID")

	if MsgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID messaggio vuoto"})
		return
	}

	// devo trovare il messaggio nella conversazione, eliminarlo da lì e successivamente dal DB dei Messaggi + libero l'ID
	MsgStruct, exists := scs.MsgDB[MsgID]

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nun trovo ur msg in the MsgDB man"})
		return
	}

	if _, exists2 := scs.GenericDB[MsgID]; !exists2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unregistered GenericID"})
		return
	}

	if MsgStruct.Author != sender_struct {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is not YOUR message!"})
		return
	}

	//cerco il messaggio nella convo.
	found := false
	var foundConvo *scs.ConversationELT
	toDelete := 0
	for _, convo := range scs.ConvoDB {

		for i2, msg := range convo.Messages {
			if msg.MsgID == MsgID {
				found = true
				foundConvo = convo
				toDelete = i2
				break
			}
		}
	}
	if !found {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant find ID through all convos."})
		return
	}

	// elimino il messaggio dall'array...
	foundConvo.Messages = append(foundConvo.Messages[:toDelete], foundConvo.Messages[toDelete+1:]...)
	// elimino dai DBs...
	delete(scs.MsgDB, MsgID)
	delete(scs.GenericDB, MsgID)

	// devo fixare il preview della convo ora...
	if len(foundConvo.Messages) != 0 {
		UpdateConversationWLastMSG(foundConvo)
	} else {
		foundConvo.DateLastMessage = time.Now()
		foundConvo.Preview = "*msg deleted* :)"
	}

	c.JSON(http.StatusOK, gin.H{"Success": "Message deleted", "MSG": MsgStruct})
}

// PUT, path /groups/members
func AddToGroup(c *gin.Context) {

	//autorizzo il current user
	logged_username, logged_struct, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	//estraggo l'ID conversazione_gruppo dalla query (se c'è)
	ConvoID := c.DefaultQuery("ID", "")

	// mi passa l'utente nel request body (text/plain)

	userToAdd, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nun riesc a legg u body"})
		return
	}

	if string(userToAdd) == logged_username {
		c.JSON(http.StatusBadRequest, gin.H{"error": "stai ad aggiunge te stesso ao"})
		return
	}

	structToAdd, exists := scs.UserDB[string(userToAdd)]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "sto utente non è nel db..."})
		return
	}

	if ConvoID == "" {

		// Creo Gruppo e aggiungo i due user alla lista utenti...
		userList := []*scs.User{logged_struct, structToAdd}
		g := ConstrGroup(userList)
		scs.UserConvosDB[logged_username] = append(scs.UserConvosDB[logged_username], g.Conversation)
		scs.UserConvosDB[structToAdd.Username] = append(scs.UserConvosDB[structToAdd.Username], g.Conversation)

		c.JSON(http.StatusOK, gin.H{"success": "creato gruppo + aggiunto utenti", "group": g})
		return

	} else {

		if _, exists := scs.GroupDB[ConvoID]; !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cant find convo (group) w that ID"})
			return
		}

		// controllo se appartengo al gruppo...
		found := false
		for _, user := range scs.GroupDB[ConvoID].Users {
			if user == logged_struct {
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bro ur not in that group... wyd?"})
			return
		}

		// controllo se l'utente da aggiungere sta già nel gruppo...
		found2 := false
		for _, user := range scs.GroupDB[ConvoID].Users {
			if user == logged_struct {
				found2 = true
				break
			}
		}
		if found2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already in the group"})
			return
		}
		// Trovo Gruppo dal GroupDB e aggiungo l'utente

		scs.GroupDB[ConvoID].Users = append(scs.GroupDB[ConvoID].Users, structToAdd)
		c.JSON(http.StatusOK, gin.H{"success": "aggiunto utente al gruppo indicato"})
		return

	}

}

// DELETE, path /groups/members
func LeaveGroup(c *gin.Context) {

	//autorizzo il current user
	_, logged_struct, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	//estraggo l'ID conversazione_gruppo dalla query (se c'è)
	ConvoID := c.DefaultQuery("ID", "")

	g, exists := scs.GroupDB[ConvoID]

	if exists {
		toDelete := -1
		found := false
		for i, user := range g.Users {
			if user == logged_struct {
				toDelete = i
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusForbidden, gin.H{"Error": " You are not in that group! >:O "})
			return

		}
		g.Users = append(g.Users[:toDelete], g.Users[toDelete+1:]...)
		delete(scs.UserConvosDB, ConvoID)
		c.JSON(http.StatusOK, gin.H{"Success": "You left the group and the convo was removed from your list..."})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Convo (group) ID!"})

}

// PUT, path /groups/{ID}/name
func SetGroupName(c *gin.Context) {

	//autorizzo il current user
	_, logged_struct, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	ConvoID := c.Param("ID")

	group, exists := scs.GroupDB[ConvoID]

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Cant find convo (group) ID"})
		return
	}

	groupName, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nun riesc a legg u body"})
		return
	}

	found := false
	for _, user := range group.Users {
		if user == logged_struct {
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusForbidden, gin.H{"Error": " You are not in that group! >:O "})
		return

	}

	group.Name = string(groupName)
	c.JSON(http.StatusOK, gin.H{"Success": "Name Updated!", "Group": group})

}

// PUT, path /groups/{ID}/photo
func SetGroupPhoto(c *gin.Context) {

	//autorizzo il current user
	_, logged_struct, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	ConvoID := c.Param("ID")

	group, exists := scs.GroupDB[ConvoID]

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Cant find convo (group) ID"})
		return
	}

	found := false
	for _, user := range group.Users {
		if user == logged_struct {
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusForbidden, gin.H{"Error": " You are not in that group! >:O "})
		return

	}

	// aggiornamento foto con quella che si trova nel body (c)
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
	nudePic, err := io.ReadAll(openedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "che hai scritto nel file?? non leggo..."})
		return
	}

	group.GroupPhoto = nudePic

	c.JSON(http.StatusOK, gin.H{"Success": "Photo Updated!", "Group": group})

}

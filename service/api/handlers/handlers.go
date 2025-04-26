package handlers

import (
	"fmt"
	"io"
	"net/http"
	"time"

	scs "github.com/Frisbon/hungrymonke/service/api/structures"
	"github.com/gin-gonic/gin"
)

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
	usernameString = usernameString[1 : len(usernameString)-1]
	println("Ho ricevuto [" + usernameString + "] come nuovo username...")

	if usernameString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input non valido, stringa vuota"})
		return
	}
	if _, exists := scs.UserDB[usernameString]; exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Questo nome è già occupato! :("})
		return
	}

	//trasferisco le convo sull'altro nickname
	var convos = scs.UserConvosDB[username]
	delete(scs.UserConvosDB, username)
	scs.UserConvosDB[usernameString] = convos

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

// PUT, path /users/me/photo
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

	contentType := http.DetectContentType(content)

	// Aggiorna l'utente con il percorso del file salvato
	user.Photo = content
	user.PhotoMimeType = contentType // tipo MIME
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
		c.JSON(http.StatusOK, gin.H{"Error": "Non ci sono conversazioni per questo utente", "Username": username})
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

	*/

	statusErr := statusUpdater(found_conversation, logged_struct)
	if statusErr {
		c.JSON(http.StatusInternalServerError, gin.H{"Status": "oh sh*t"}) // todo add to open api (below too)
	}

	// Risposta con la conversazione trovata
	c.JSON(http.StatusOK, gin.H{
		"message":      "Conversazione trovata",
		"conversation": found_conversation,
	})

}

// POST, path /conversations/messages/?=ID
func SendMessage(c *gin.Context) {

	//autorizzo il current user
	SenderUsername, sender_struct, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	// leggo il messaggio e il username a chi inviarlo (se c'è) dal body
	type requestLol struct {
		RecipientUsername string       `json:"recipientusername"`
		Message           scs.Content  `json:"message"`
		ReplyingTo        *scs.Message `json:"replyingto,omitempty"`
	}

	var req requestLol
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON non valido", "details": err.Error()})
		return
	}

	var reactlist []scs.Reaction
	var repMsg *scs.Message

	if req.ReplyingTo != nil {

		x, exists := scs.MsgDB[req.ReplyingTo.MsgID]

		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Non esiste messaggio a cui stai rispondendo"})
			return
		}

		repMsg = x
	} else {
		repMsg = nil
	}

	// DEFINISCO NUOVO MESSAGGIO CON PARAMETRI GIUSTI
	msgc := MsgCONSTR{
		Timestamp:  time.Now(),
		Content:    req.Message,
		Status:     "delivered",
		Reactions:  reactlist,
		Author:     sender_struct,
		ReplyingTo: repMsg,
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

	// se tra parametri ho solo ID
	if conversationID != "" && req.RecipientUsername == "" {

		// prendo conversazione con l'ID che mi è stato fornito.
		convo := scs.ConvoDB[conversationID]
		//aggiungo il messaggio alla lista dei messaggi su quella convo
		convo.Messages = append(convo.Messages, newMessage)
		UpdateConversationWLastMSG(convo)

		// se mando un msg => leggo i messaggi dell'altro utente => aggiorno i loro status.
		//devo distinguere se la convo è privata o di gruppo
		statusErr := statusUpdater(convo, sender_struct)
		if statusErr {
			c.JSON(http.StatusInternalServerError, gin.H{"Status": "oh sh*t"}) // todo add to open api (below too)
		}

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
	var seenby_list []scs.User
	newMsg := ConstrMessage(MsgCONSTR{
		Timestamp:   time.Now(),
		Content:     msg.Content,
		Author:      sender_struct, //maschero da chi inoltro
		Status:      "delivered",
		Reactions:   react_list, //resetto lista reazioni
		SeenBy:      seenby_list,
		IsForwarded: true,
	})

	convo.Messages = append(convo.Messages, newMsg)

	c.JSON(http.StatusOK, gin.H{"success": "Messaggio inoltrato", "MSG": newMsg})
	UpdateConversationWLastMSG(convo)
	//bwt non faccio controlli delle fonti del messaggio ecc ecc perchè non mi va, basta che funge :)
}

// POST, path /messages/{ID}/comments
func CommentMSG(c *gin.Context) {

	//nota: se commento => visualizzo il messaggio

	//autorizzo il current user
	_, sender_struct, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	//estraggo l'ID messaggio dal path (se c'è)
	MsgID := c.Param("ID")

	fmt.Println("Mi hai passato l'ID: ", MsgID)
	fmt.Println("Controllo se esiste: ", scs.MsgDB[MsgID])

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
		c.JSON(http.StatusNotFound, gin.H{"error": "nun trovo u msg man"})
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

	//recupero la convo in cui si trova il messaggio
	for _, convo := range scs.ConvoDB {
		for _, imsg := range convo.Messages {
			if MsgStruct.MsgID == imsg.MsgID {

				statusErr := statusUpdater(convo, sender_struct)
				if statusErr {
					c.JSON(http.StatusInternalServerError, gin.H{"Status": "oh sh*t"}) // todo add to open api (below too)
				}
				break
			}
		}
	}

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
		c.JSON(http.StatusNotFound, gin.H{"error": "nun trovo ur msg in the MsgDB man"})
		return
	}

	if _, exists2 := scs.GenericDB[MsgID]; !exists2 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unregistered GenericID"})
		return
	}

	if MsgStruct.Author != sender_struct {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "This is not YOUR message!"})
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
// PUT, path /groups/members
func AddToGroup(c *gin.Context) {

	//autorizzo il current user
	logged_username, logged_struct, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	//estraggo l'ID conversazione_gruppo dalla query (se c'è)
	ConvoID := c.DefaultQuery("ID", "")

	// Definisci una struct per il body della richiesta JSON (array di usernames)
	type AddUsersRequest struct {
		Usernames []string `json:"usernames"`
	}

	var req AddUsersRequest
	// Binda il body JSON alla struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Errore nel leggere o parsare il body JSON", "details": err.Error()})
		return
	}

	// Raccogli gli utenti da aggiungere/includere (converti nomi utente in struct User)
	usersToAddStructs := []*scs.User{}
	addedUsernames := []string{} // Per la risposta

	for _, username := range req.Usernames {
		if username == "" {
			continue // Salta username vuoti
		}
		if username == logged_username {
			continue // Ignora il creatore se è nella lista
		}
		userStruct, exists := scs.UserDB[username]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Utente '" + username + "' non trovato nel DB"})
			return // Interrompi se un utente non esiste
		}
		usersToAddStructs = append(usersToAddStructs, userStruct)
		addedUsernames = append(addedUsernames, username) // Aggiungi solo se l'utente esiste
	}

	// Se non ci sono utenti validi nella richiesta (oltre il creatore)
	if len(usersToAddStructs) == 0 && ConvoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nessun utente valido specificato per creare un nuovo gruppo"})
		return
	}

	if ConvoID == "" {
		// Creo Gruppo e aggiungo i due user alla lista utenti... (creatore + utenti richiesti)
		initialMembers := []*scs.User{logged_struct}
		initialMembers = append(initialMembers, usersToAddStructs...)

		if len(initialMembers) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Un nuovo gruppo richiede almeno due membri (tu + almeno uno)"})
			return
		}

		g := ConstrGroup(initialMembers)

		c.JSON(http.StatusOK, gin.H{
			"success":     "Nuovo gruppo creato con successo",
			"convoID":     g.Conversation.ConvoID,
			"added_users": addedUsernames,
		})
		return

	} else {
		// Trovo Gruppo dal GroupDB e aggiungo l'utente (o gli utenti)

		group, exists := scs.GroupDB[ConvoID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conversazione (gruppo) non trovata con quell'ID"})
			return
		}

		// controllo se appartengo al gruppo...
		isMember := false
		for _, user := range group.Users {
			if user.Username == logged_username {
				isMember = true
				break
			}
		}
		if !isMember {
			c.JSON(http.StatusForbidden, gin.H{"error": "Non fai parte di questo gruppo"})
			return
		}

		// Aggiungi gli utenti validi al gruppo esistente
		usersActuallyAdded := []string{}
		for _, userStructToAdd := range usersToAddStructs {

			// controllo se l'utente da aggiungere sta già nel gruppo...
			alreadyInGroup := false
			for _, existingUser := range group.Users {
				if existingUser.Username == userStructToAdd.Username {
					alreadyInGroup = true
					break
				}
			}

			if !alreadyInGroup {
				// Aggiungi l'utente al gruppo
				group.Users = append(group.Users, userStructToAdd)
				// Aggiorna UserConvosDB per il nuovo membro
				if _, exists := scs.UserConvosDB[userStructToAdd.Username]; !exists {
					scs.UserConvosDB[userStructToAdd.Username] = scs.Conversations{}
				}
				scs.UserConvosDB[userStructToAdd.Username] = append(scs.UserConvosDB[userStructToAdd.Username], group.Conversation)
				usersActuallyAdded = append(usersActuallyAdded, userStructToAdd.Username) // Aggiunto con successo
			}
		}

		// Se nessun utente è stato effettivamente aggiunto (es. erano già tutti dentro)
		if len(usersActuallyAdded) == 0 && len(usersToAddStructs) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"success":     "Tutti gli utenti specificati erano già membri del gruppo",
				"convoID":     ConvoID,
				"added_users": []string{}, // Lista vuota perché nessuno è stato aggiunto ORA
			})
			return // Esci qui se non è stato aggiunto nessuno di nuovo
		}

		c.JSON(http.StatusOK, gin.H{
			"success":     "Utenti aggiunti al gruppo esistente",
			"convoID":     ConvoID,
			"added_users": usersActuallyAdded,
		})
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
		c.JSON(http.StatusNoContent, gin.H{"Success": "You left the group and the convo was removed from your list..."})
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

	contentType := http.DetectContentType(nudePic)

	group.GroupPhoto = nudePic
	group.PhotoMimeType = contentType

	c.JSON(http.StatusOK, gin.H{"Success": "Photo Updated!", "Group": group})

}

/* HANDLER NON RICHIESTI */

// PUT /utils/getconvoinfo/{ID}
/* Serve al front-end per capire se una convo è privata o di gruppo, per renderizzare le info nella convo-list*/
func GetConvoInfo(c *gin.Context) {

	//autorizzo il current user
	_, logged_struct, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	ConvoID := c.Param("ID")

	group, exists := scs.GroupDB[ConvoID]
	private, exists2 := scs.PrivateDB[ConvoID]

	// se è una convo di gruppo
	if exists {
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

		c.JSON(http.StatusOK, gin.H{"Group": group})
	} else if exists2 {
		if logged_struct != private.FirstUser && logged_struct != private.SecondUser {
			c.JSON(http.StatusForbidden, gin.H{"Error": " You are not in that convo! >:O "})
			return
		}

		c.JSON(http.StatusOK, gin.H{"PrivateConvo": private})

	} else {
		//todo check if its at least in ConvoDB and theres an alignment error.\
		c.JSON(http.StatusNoContent, gin.H{"Error": "no convo bro"})
	}

}

// GET, path /utils/listUsers
func ListUsers(c *gin.Context) {

	if len(scs.UserDB) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Non trovo Utenti..."})
	} else {

		var array []scs.User
		for _, v := range scs.UserDB {
			array = append(array, scs.User{Username: v.Username, Photo: v.Photo, PhotoMimeType: v.PhotoMimeType})
		}
		c.JSON(http.StatusOK, gin.H{"Users": array})

	}

}

// POST, path /utils/createConvo
func CreatePrivateConvo(c *gin.Context) {

	//autorizzo il current user
	Current_Username, _, er1 := QuickAuth(c)
	if len(er1) != 0 {
		return
	}

	// leggo l'ID messggio dal body
	type requestLol struct {
		SecondUsername string `json:"SecondUsername"`
	}

	var req requestLol
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON non valido", "details": err.Error()})
		return
	}

	fmt.Println("Current:" + Current_Username + " req.SecondUsername=" + req.SecondUsername)
	// Recupera le struct User per i due usernames
	user1, user1Exists := scs.UserDB[Current_Username]
	user2, user2Exists := scs.UserDB[req.SecondUsername]

	if !user1Exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": user1.Username + " not exists"})
		return
	}
	if !user2Exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": user2.Username + " not exists"})
		return
	}

	// Per trovare la conversazione esistente in modo efficiente,
	// iteriamo sulle conversazioni di uno degli utenti e controlliamo se l'altra parte è l'altro utente.
	// Assumiamo che UserConvosDB sia popolato correttamente.

	user1Convos, user1HasConvos := scs.UserConvosDB[user1.Username]
	if user1HasConvos {
		for _, convoELT := range user1Convos {
			// Recupera la struct Private associata a questa conversazione
			privateConvo, exists := scs.PrivateDB[convoELT.ConvoID]
			if exists {
				// Controlla se questa conversazione privata include entrambi gli utenti
				isBetweenUsers := (privateConvo.FirstUser.Username == user1.Username && privateConvo.SecondUser.Username == user2.Username) ||
					(privateConvo.FirstUser.Username == user2.Username && privateConvo.SecondUser.Username == user1.Username)
				fmt.Println("isBetweenUsers=", isBetweenUsers)
				if isBetweenUsers {
					// Trovata conversazione esistente
					c.JSON(http.StatusOK, gin.H{"convoID": convoELT.ConvoID})
					return
				}
			}
			// Nota: Se una ConvoID esiste in UserConvosDB ma non in PrivateDB,
			// c'è un'inconsistenza nei dati. Qui assumiamo che i DB siano consistenti.
		}
	}
	fmt.Println("Non ho trovato convo, creo...")
	// Se non è stata trovata nessuna conversazione, creane una nuova
	newConvoID := GenerateRandomString(6)

	// Crea la nuova struct ConversationELT
	newConvoELT := &scs.ConversationELT{
		ConvoID:         newConvoID,
		DateLastMessage: time.Now(), // O potresti inizializzarla a zero o quando arriva il primo messaggio
		Preview:         "",         // Inizializza con stringa vuota o placeholder
		Messages:        []*scs.Message{},
	}

	// Aggiungi la nuova conversazione al ConvoDB
	scs.ConvoDB[newConvoID] = newConvoELT

	// Crea la nuova struct Private
	// Convenzione: ordina gli utenti per username per avere consistenza nella chiave (anche se usiamo ConvoID)
	// e potenzialmente per futuri indici/ricerche. Qui li mettiamo semplicemente come recuperati.
	newPrivate := &scs.Private{
		Conversation: newConvoELT,
		FirstUser:    user1,
		SecondUser:   user2,
	}

	// Aggiungi la nuova conversazione privata al PrivateDB
	scs.PrivateDB[newConvoID] = newPrivate

	// Aggiungi la nuova conversazione alle liste di conversazioni degli utenti
	// Inizializza la slice se l'utente non ha ancora conversazioni in UserConvosDB
	if !user1HasConvos {
		scs.UserConvosDB[user1.Username] = scs.Conversations{}
	}
	scs.UserConvosDB[user1.Username] = append(scs.UserConvosDB[user1.Username], newConvoELT)

	// Controlla e aggiungi anche per il secondo utente
	_, user2HasConvos := scs.UserConvosDB[user2.Username]
	if !user2HasConvos {
		scs.UserConvosDB[user2.Username] = scs.Conversations{}
	}
	scs.UserConvosDB[user2.Username] = append(scs.UserConvosDB[user2.Username], newConvoELT)
	c.JSON(http.StatusOK, gin.H{"convoID": newConvoID})

}

// POST, path /utils/createGroup
func CreateGroupConvo(c *gin.Context) {

	_, Current_UserStruct, er1 := QuickAuth(c)
	if len(er1) != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type requestLol struct {
		Users              []string `json:"Users"` // Assumi che l'utente invii un array con almeno l'Username, SENZA il creatore
		GroupName          string   `json:"GroupName"`
		GroupPicture       []byte   `json:"GroupPicture,omitempty"` // Base64 decodificata in []byte
		GroupPhotoMimeType string   `json:"GroupPhotoMimeType,omitempty"`
	}

	var req requestLol
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON non valido", "details": err.Error()})
		return
	}

	if req.GroupName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group name cannot be empty"})
		return
	}

	groupMembersSlice := []*scs.User{Current_UserStruct}

	// Itera sugli utenti forniti dal frontend (gli altri membri)
	for _, userReq := range req.Users {
		// Verifica se l'utente fornito esiste nel UserDB
		user, exists := scs.UserDB[userReq]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User '" + userReq + "' not found"})
			return
		}
		// Aggiungi l'utente verificato direttamente alla slice
		groupMembersSlice = append(groupMembersSlice, user)
	}

	// Un gruppo deve avere almeno 2 membri (il creatore + almeno un altro utente dalla richiesta)
	if len(groupMembersSlice) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group must have at least two members"})
		return
	}

	// Genera un ID univoco
	newConvoID := GenerateRandomString(6)

	// Crea la ConversationELT
	newConvoELT := &scs.ConversationELT{
		ConvoID:         newConvoID,
		DateLastMessage: time.Now(),
		Preview:         "[New Group Created! 🥳]",
		Messages:        []*scs.Message{},
	}

	// Crea la struct Group
	newGroup := &scs.Group{
		Conversation:  newConvoELT,
		Name:          req.GroupName,
		GroupPhoto:    req.GroupPicture,
		PhotoMimeType: req.GroupPhotoMimeType,
		Users:         groupMembersSlice, // Slice di membri verificati (creatore + altri)
	}

	// Aggiungi ai DB globali
	scs.ConvoDB[newConvoID] = newConvoELT
	scs.GroupDB[newConvoID] = newGroup

	// Aggiorna UserConvosDB per ogni membro
	for _, member := range groupMembersSlice {
		if _, exists := scs.UserConvosDB[member.Username]; !exists {
			scs.UserConvosDB[member.Username] = scs.Conversations{}
		}
		scs.UserConvosDB[member.Username] = append(scs.UserConvosDB[member.Username], newConvoELT)
	}

	// Risposta di successo
	c.JSON(http.StatusOK, gin.H{"convoID": newConvoID, "message": "Group created successfully"})

}

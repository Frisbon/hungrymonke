package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	scs "github.com/Frisbon/hungrymonke/service/api/structures"
	"github.com/gin-gonic/gin"
)

// POST, path /users/me/username
func SetMyUsername(c *gin.Context) {
	scs.DBMutex.Lock()         //*
	defer scs.DBMutex.Unlock() //*

	username, user, err := QuickAuth(c)
	if err != nil {
		return
	}

	newUsername, errBody := c.GetRawData()
	if errBody != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nun riesc a legg u body"})
		return
	}

	usernameString := string(newUsername)
	usernameString = usernameString[1 : len(usernameString)-1]
	log.Println("Ho ricevuto [" + usernameString + "] come nuovo username...")

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
	tokenString, errToken := GeneraToken(usernameString) // Supponiamo che la funzione generaToken prenda il nuovo username
	if errToken != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Non riesco a creare il token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Username aggiornato lessgo", "user": user, "new_token": tokenString})
}

// PUT, path /users/me/photo
func SetMyPhoto(c *gin.Context) {
	scs.DBMutex.Lock()         //*
	defer scs.DBMutex.Unlock() //*

	username, user, err := QuickAuth(c)
	if err != nil {
		return
	}

	// aggiornamento foto con quella che si trova nel body (c)
	// estraggo il file dalla richiesta multipart/form-data
	file, errFile := c.FormFile("file")
	if errFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File invalido :P"})
		return
	}

	// apro il file appena estratto
	openedFile, errOpen := file.Open()
	if errOpen != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "non riesco ad apri, ao me apri"})
		return
	}
	defer openedFile.Close()

	// Leggo il file in modo da salvarlo come []byte dentro al content. Lo faccio perchè PhotoFile è di tipo []byte
	content, errRead := io.ReadAll(openedFile)
	if errRead != nil {
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
	scs.DBMutex.RLock()         //*
	defer scs.DBMutex.RUnlock() //*

	username, _, err := QuickAuth(c)
	if err != nil {
		return
	}

	if len(scs.UserConvosDB[username]) == 0 {
		c.JSON(http.StatusOK, gin.H{"Error": "Non ci sono conversazioni per questo utente", "Username": username})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Username": username, "User Conversations": scs.UserConvosDB[username]})
}

// GET, path /conversations/:ID
func GetMyConversation(c *gin.Context) {
	scs.DBMutex.Lock()         //*
	defer scs.DBMutex.Unlock() //*

	username, logged_struct, err := QuickAuth(c)
	if err != nil {
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
	user_conversations, exists := scs.UserConvosDB[username]
	if !exists {
		c.JSON(http.StatusNoContent, gin.H{"error": "Non ho conversazioni per questo utente."})
		return
	}

	for _, conversation := range user_conversations {
		if conversation.ConvoID == conversationID {
			found_conversation = conversation
			break
		}
	}
	if found_conversation == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversazione non trovata"})
		return
	}

	statusErr := statusUpdater(found_conversation, logged_struct)
	if statusErr {
		c.JSON(http.StatusInternalServerError, gin.H{"Status": "oh sh*t"})
		return
	}

	// Risposta con la conversazione trovata
	c.JSON(http.StatusOK, gin.H{
		"message":      "Conversazione trovata",
		"conversation": found_conversation,
	})
}

// POST, path /conversations/messages/?=ID
func SendMessage(c *gin.Context) {
	scs.DBMutex.Lock()         //*
	defer scs.DBMutex.Unlock() //*

	SenderUsername, sender_struct, err := QuickAuth(c)
	if err != nil {
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

	msgc := MsgCONSTR{
		Timestamp:  time.Now(),
		Content:    req.Message,
		Status:     scs.Delivered,
		Reactions:  reactlist,
		Author:     sender_struct,
		ReplyingTo: repMsg,
	}

	newMessage := ConstrMessage(msgc)

	conversationID := c.DefaultQuery("ID", "")

	if conversationID != "" && req.RecipientUsername != "" {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Non puoi inviare entrambi ID e Nickname nel json! :O"})
		return
	}

	if conversationID != "" && req.RecipientUsername == "" {
		convo := scs.ConvoDB[conversationID]
		convo.Messages = append(convo.Messages, newMessage)
		UpdateConversationWLastMSG(convo)
		statusErr := statusUpdater(convo, sender_struct)
		if statusErr {
			c.JSON(http.StatusInternalServerError, gin.H{"Status": "oh sh*t"})
		}
		c.JSON(http.StatusOK, gin.H{"Status": "Inviato sulla convo indicata dall'ID!", "Conversation": &convo})
		return
	}

	if conversationID == "" && req.RecipientUsername != "" {
		var found_private *scs.Private
		for _, private := range scs.PrivateDB {
			isMatch := (private.FirstUser.Username == SenderUsername && private.SecondUser.Username == req.RecipientUsername) ||
				(private.FirstUser.Username == req.RecipientUsername && private.SecondUser.Username == SenderUsername)
			if isMatch {
				found_private = private
				break
			}
		}

		if found_private == nil {
			x := ConvoCONSTR{
				DateLastMessage: newMessage.Timestamp,
				Messages:        []*scs.Message{newMessage},
			}
			newConversation := ConstrConvo(x)
			UpdateConversationWLastMSG(newConversation)
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
		} else {
			found_private.Conversation.Messages = append(found_private.Conversation.Messages, newMessage)
			PrivateMsgStatusUpdater(found_private.Conversation, sender_struct)
			UpdateConversationWLastMSG(found_private.Conversation)
			c.JSON(http.StatusOK, gin.H{"Status": "Ho trovato una convo esistente con questo utente... Messaggio inviato!", "Conversation": found_private.Conversation})
			return
		}
	}
}

// POST, path /messages/{ID}/forward
func ForwardMSG(c *gin.Context) {
	scs.DBMutex.Lock()         //*
	defer scs.DBMutex.Unlock() //*

	_, sender_struct, err := QuickAuth(c)
	if err != nil {
		return
	}

	type requestLol struct {
		ConvoID string `json:"ConvoID"`
	}

	var req requestLol
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON non valido", "details": err.Error()})
		return
	}

	MsgID := c.Param("ID")

	msg, exists := scs.MsgDB[MsgID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "MSG not found"})
		return
	}

	convo, exists2 := scs.ConvoDB[req.ConvoID]
	if !exists2 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	var react_list []scs.Reaction
	var seenby_list []scs.User
	newMsg := ConstrMessage(MsgCONSTR{
		Timestamp:   time.Now(),
		Content:     msg.Content,
		Author:      sender_struct,
		Status:      scs.Delivered,
		Reactions:   react_list,
		SeenBy:      seenby_list,
		IsForwarded: true,
	})

	convo.Messages = append(convo.Messages, newMsg)
	UpdateConversationWLastMSG(convo)

	c.JSON(http.StatusOK, gin.H{"success": "Messaggio inoltrato", "MSG": newMsg})
}

// POST, path /messages/{ID}/comments
func CommentMSG(c *gin.Context) {
	scs.DBMutex.Lock()         //*
	defer scs.DBMutex.Unlock() //*

	_, sender_struct, err := QuickAuth(c)
	if err != nil {
		return
	}

	MsgID := c.Param("ID")
	if MsgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID vuoto"})
		return
	}

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

	for _, R := range MsgStruct.Reactions {
		if R.Author.Username == sender_struct.Username {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bro u already commented, delete ur last reaction"})
			return
		}
	}

	MsgStruct.Reactions = append(MsgStruct.Reactions, scs.Reaction{
		Timestamp: time.Now(),
		Author:    sender_struct,
		Emoticon:  req.Emoticon,
	})

	for _, convo := range scs.ConvoDB {
		for _, imsg := range convo.Messages {
			if MsgStruct.MsgID == imsg.MsgID {
				statusErr := statusUpdater(convo, sender_struct)
				if statusErr {
					c.JSON(http.StatusInternalServerError, gin.H{"Status": "oh sh*t"})
				}
				break
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"Success": "Messaggio inoltrato", "MSG": MsgStruct})
}

// DELETE, path /messages/{ID}/comments
func UncommentMSG(c *gin.Context) {
	scs.DBMutex.Lock()         //*
	defer scs.DBMutex.Unlock() //*

	_, sender_struct, err := QuickAuth(c)
	if err != nil {
		return
	}

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

	found := false
	toDelete := 0
	for i, R := range MsgStruct.Reactions {
		if R.Author.Username == sender_struct.Username {
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
	scs.DBMutex.Lock()         //*
	defer scs.DBMutex.Unlock() //*

	_, sender_struct, err := QuickAuth(c)
	if err != nil {
		return
	}

	MsgID := c.Param("ID")
	if MsgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID messaggio vuoto"})
		return
	}

	MsgStruct, exists := scs.MsgDB[MsgID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "nun trovo ur msg in the MsgDB man"})
		return
	}

	if _, exists2 := scs.GenericDB[MsgID]; !exists2 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unregistered GenericID"})
		return
	}

	if MsgStruct.Author.Username != sender_struct.Username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "This is not YOUR message!"})
		return
	}

	var foundConvo *scs.ConversationELT
	toDelete := -1
	for _, convo := range scs.ConvoDB {
		for i2, msg := range convo.Messages {
			if msg.MsgID == MsgID {
				foundConvo = convo
				toDelete = i2
				break
			}
		}
		if foundConvo != nil {
			break
		}
	}
	if foundConvo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant find ID through all convos."})
		return
	}

	foundConvo.Messages = append(foundConvo.Messages[:toDelete], foundConvo.Messages[toDelete+1:]...)
	delete(scs.MsgDB, MsgID)
	delete(scs.GenericDB, MsgID)

	UpdateConversationWLastMSG(foundConvo)

	c.JSON(http.StatusOK, gin.H{"Success": "Message deleted", "MSG": MsgStruct})
}

// PUT, path /groups/members
func AddToGroup(c *gin.Context) {
	scs.DBMutex.Lock()         //*
	defer scs.DBMutex.Unlock() //*

	logged_username, _, err := QuickAuth(c)
	if err != nil {
		return
	}

	ConvoID := c.DefaultQuery("ID", "")

	type AddUsersRequest struct {
		Usernames []string `json:"Users"`
	}

	var req AddUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Errore nel leggere o parsare il body JSON", "details": err.Error()})
		return
	}

	usersToAddStructs := []*scs.User{}
	for _, username := range req.Usernames {
		if username == "" || username == logged_username {
			continue
		}
		userStruct, exists := scs.UserDB[username]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Utente '" + username + "' non trovato nel DB"})
			return
		}
		usersToAddStructs = append(usersToAddStructs, userStruct)
	}

	if len(usersToAddStructs) == 0 && ConvoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nessun utente valido specificato per creare un nuovo gruppo"})
		return
	}

	group, exists := scs.GroupDB[ConvoID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversazione (gruppo) non trovata con quell'ID"})
		return
	}

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

	usersActuallyAdded := []string{}
	for _, userStructToAdd := range usersToAddStructs {
		alreadyInGroup := false
		for _, existingUser := range group.Users {
			if existingUser.Username == userStructToAdd.Username {
				alreadyInGroup = true
				break
			}
		}

		if !alreadyInGroup {
			group.Users = append(group.Users, userStructToAdd)
			scs.UserConvosDB[userStructToAdd.Username] = append(scs.UserConvosDB[userStructToAdd.Username], group.Conversation)
			usersActuallyAdded = append(usersActuallyAdded, userStructToAdd.Username)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     "Utenti aggiunti al gruppo esistente",
		"convoID":     ConvoID,
		"added_users": usersActuallyAdded,
	})
}

// DELETE, path /groups/members
func LeaveGroup(c *gin.Context) {
	scs.DBMutex.Lock()         //*
	defer scs.DBMutex.Unlock() //*

	logged_username, _, err := QuickAuth(c)
	if err != nil {
		return
	}

	ConvoID := c.DefaultQuery("ID", "")
	g, exists := scs.GroupDB[ConvoID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Convo (group) ID!"})
		return
	}

	toDelete := -1
	found := false
	for i, user := range g.Users {
		if user.Username == logged_username {
			toDelete = i
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusForbidden, gin.H{"Error": "You are not in that group! >:O"})
		return
	}

	g.Users = append(g.Users[:toDelete], g.Users[toDelete+1:]...)

	userConvos, ok := scs.UserConvosDB[logged_username]
	convoIndexToDelete := -1
	for i, convo := range userConvos {
		if convo.ConvoID == ConvoID {
			convoIndexToDelete = i
			break
		}
	}

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Error": "cant find u mate"})
		return
	}

	if convoIndexToDelete != -1 {
		scs.UserConvosDB[logged_username] = append(userConvos[:convoIndexToDelete], userConvos[convoIndexToDelete+1:]...)
		fmt.Printf("Conversation %s removed from user %s list.\n", ConvoID, logged_username)
	} else {
		fmt.Printf("Warning: Conversation %s not found in user %s conversation list.\n", ConvoID, logged_username)
	}

	if len(g.Users) == 0 {
		delete(scs.GroupDB, ConvoID)
		delete(scs.ConvoDB, ConvoID)
		fmt.Printf("Group %s is now empty and has been deleted from global DBs.\n", ConvoID)
	}

	c.JSON(http.StatusOK, gin.H{"Success": "You left the group and the convo was removed from your list."})
}

// PUT, path /groups/{ID}/name
func SetGroupName(c *gin.Context) {
	scs.DBMutex.Lock()         //*
	defer scs.DBMutex.Unlock() //*

	_, logged_struct, err := QuickAuth(c)
	if err != nil {
		return
	}

	ConvoID := c.Param("ID")
	group, exists := scs.GroupDB[ConvoID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Cant find convo (group) ID"})
		return
	}

	groupName, errBody := c.GetRawData()
	if errBody != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nun riesc a legg u body"})
		return
	}
	groupName = groupName[1 : len(groupName)-1]

	found := false
	for _, user := range group.Users {
		if user.Username == logged_struct.Username {
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
	scs.DBMutex.Lock()         //*
	defer scs.DBMutex.Unlock() //*

	_, logged_struct, err := QuickAuth(c)
	if err != nil {
		return
	}

	ConvoID := c.Param("ID")

	file, errFile := c.FormFile("file")
	if errFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File invalido :P"})
		return
	}

	openedFile, errOpen := file.Open()
	if errOpen != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "non riesco ad apri, ao me apri"})
		return
	}
	defer openedFile.Close()

	nudePic, errRead := io.ReadAll(openedFile)
	if errRead != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "che hai scritto nel file?? non leggo..."})
		return
	}
	contentType := http.DetectContentType(nudePic)

	group, exists := scs.GroupDB[ConvoID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Cant find convo (group) ID"})
		return
	}

	found := false
	for _, user := range group.Users {
		if user.Username == logged_struct.Username {
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusForbidden, gin.H{"Error": " You are not in that group! >:O "})
		return
	}

	group.GroupPhoto = nudePic
	group.PhotoMimeType = contentType

	c.JSON(http.StatusOK, gin.H{"Success": "Photo Updated!", "Group": group})
}

/* HANDLER NON RICHIESTI */

// PUT /utils/getconvoinfo/{ID}
func GetConvoInfo(c *gin.Context) {
	scs.DBMutex.RLock()         //*
	defer scs.DBMutex.RUnlock() //*

	_, logged_struct, err := QuickAuth(c)
	if err != nil {
		return
	}

	ConvoID := c.Param("ID")

	group, exists := scs.GroupDB[ConvoID]
	private, exists2 := scs.PrivateDB[ConvoID]

	if exists {
		found := false
		for _, user := range group.Users {
			if user.Username == logged_struct.Username {
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusForbidden, gin.H{"Error": " You are not in that group! >:O "})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Group": group})
		return
	}

	if exists2 {
		if logged_struct.Username != private.FirstUser.Username && logged_struct.Username != private.SecondUser.Username {
			c.JSON(http.StatusForbidden, gin.H{"Error": " You are not in that convo! >:O "})
			return
		}
		c.JSON(http.StatusOK, gin.H{"PrivateConvo": private})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"Error": "no convo bro"})
}

// GET, path /utils/listUsers
func ListUsers(c *gin.Context) {
	scs.DBMutex.RLock()         //*
	defer scs.DBMutex.RUnlock() //*

	if len(scs.UserDB) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Non trovo Utenti..."})
		return
	}

	var array []scs.User
	for _, v := range scs.UserDB {
		array = append(array, scs.User{Username: v.Username, Photo: v.Photo, PhotoMimeType: v.PhotoMimeType})
	}
	c.JSON(http.StatusOK, gin.H{"Users": array})
}

// POST, path /utils/createConvo
func CreatePrivateConvo(c *gin.Context) {
	scs.DBMutex.Lock()         //*
	defer scs.DBMutex.Unlock() //*

	Current_Username, _, err := QuickAuth(c)
	if err != nil {
		return
	}

	type requestLol struct {
		SecondUsername string `json:"SecondUsername"`
	}

	var req requestLol
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON non valido", "details": err.Error()})
		return
	}

	user1, user1Exists := scs.UserDB[Current_Username]
	user2, user2Exists := scs.UserDB[req.SecondUsername]
	if !user1Exists || !user2Exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "one or more users do not exist"})
		return
	}

	user1Convos, user1HasConvos := scs.UserConvosDB[user1.Username]
	if user1HasConvos {
		for _, convoELT := range user1Convos {
			privateConvo, exists := scs.PrivateDB[convoELT.ConvoID]
			if exists {
				isBetweenUsers := (privateConvo.FirstUser.Username == user1.Username && privateConvo.SecondUser.Username == user2.Username) ||
					(privateConvo.FirstUser.Username == user2.Username && privateConvo.SecondUser.Username == user1.Username)
				if isBetweenUsers {
					c.JSON(http.StatusOK, gin.H{"convoID": convoELT.ConvoID})
					return
				}
			}
		}
	}

	newConvoID := GenerateRandomString(6)
	newConvoELT := &scs.ConversationELT{
		ConvoID:         newConvoID,
		DateLastMessage: time.Now(),
		Preview:         "",
		Messages:        []*scs.Message{},
	}
	scs.ConvoDB[newConvoID] = newConvoELT

	newPrivate := &scs.Private{
		Conversation: newConvoELT,
		FirstUser:    user1,
		SecondUser:   user2,
	}
	scs.PrivateDB[newConvoID] = newPrivate

	addConvoToUserIfNotExists(user1.Username, newConvoELT)
	addConvoToUserIfNotExists(user2.Username, newConvoELT)

	c.JSON(http.StatusOK, gin.H{"convoID": newConvoID})
}

// POST, path /utils/createGroup
func CreateGroupConvo(c *gin.Context) {
	scs.DBMutex.Lock()         //*
	defer scs.DBMutex.Unlock() //*

	_, Current_UserStruct, err := QuickAuth(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type requestLol struct {
		Users              []string `json:"Users"`
		GroupName          string   `json:"GroupName"`
		GroupPicture       []byte   `json:"GroupPicture,omitempty"`
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
	for _, userReq := range req.Users {
		user, exists := scs.UserDB[userReq]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User '" + userReq + "' not found"})
			return
		}
		groupMembersSlice = append(groupMembersSlice, user)
	}

	if len(groupMembersSlice) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group must have at least two members"})
		return
	}

	newConvoELT := ConstrConvo(ConvoCONSTR{
		DateLastMessage: time.Now(),
		Preview:         "[New Group Created! 🥳]",
		Messages:        []*scs.Message{},
	})

	newGroup := &scs.Group{
		Conversation:  newConvoELT,
		Name:          req.GroupName,
		GroupPhoto:    req.GroupPicture,
		PhotoMimeType: req.GroupPhotoMimeType,
		Users:         groupMembersSlice,
	}
	scs.GroupDB[newConvoELT.ConvoID] = newGroup

	for _, member := range groupMembersSlice {
		addConvoToUserIfNotExists(member.Username, newConvoELT)
	}

	c.JSON(http.StatusOK, gin.H{"convoID": newConvoELT.ConvoID, "message": "Group created successfully"})
}

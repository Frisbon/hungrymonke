package handlers

import (
	"encoding/json"
	"errors"

	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"

	scs "github.com/Frisbon/hungrymonke/service/api/structures"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func QuickAuth(c *gin.Context) (string, *scs.User, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
		return "", nil, errors.New("token not found")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) { return jwtKey, nil })
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return "", nil, errors.New("invalid token")
	}
	username := strings.Trim(claims.Subject, "\"")
	user, exists := scs.UserDB[username]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return "", nil, errors.New("user not found")
	}
	return username, user, nil
}

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

func isUniversalIdUnique(id string) bool {
	_, exists := scs.GenericDB[id]
	return !exists
}

func PreviewMaker(Content scs.Content) string {
	if Content.Text != nil && Content.Photo != nil {
		if len(*Content.Text) > 16 {
			return "📷 " + (*Content.Text)[:16] + "..."
		}
		return "📷 " + *Content.Text
	} else if Content.Text != nil {
		if len(*Content.Text) > 16 {
			return (*Content.Text)[:16] + "..."
		}
		return *Content.Text
	} else if Content.Photo != nil {
		return "📷 Photo..."
	} else {
		return "ERROR - INVALID CONTENT"
	}
}

func UpdateConversationWLastMSG(convo *scs.ConversationELT) {
	if len(convo.Messages) > 0 {
		lst_msg := convo.Messages[len(convo.Messages)-1]
		convo.DateLastMessage = lst_msg.Timestamp
		convo.Preview = PreviewMaker(lst_msg.Content)
	}
}

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
	scs.DBMutex.RLock()
	defer scs.DBMutex.RUnlock()
	printMap := func(name string, data interface{}) {
		v := reflect.ValueOf(data)
		if v.Kind() != reflect.Map {
			log.Printf("%s is not a map", name)
			return
		}

		if v.Len() == 0 { // <-- This is the correct and universal check
			log.Printf("%s è vuoto.\n", name)
			return
		}

		log.Printf("Contenuto di %s:\n", name)
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Printf("Errore durante la conversione in JSON di %s: %v\n", name, err)
			return
		}
		log.Println(string(jsonData))
	}
	printMap("GenericDB", scs.GenericDB)
	printMap("UserDB", scs.UserDB)
	printMap("PrivateDB", scs.PrivateDB)
	printMap("GroupDB", scs.GroupDB)
	printMap("MsgDB", scs.MsgDB)
	printMap("ConvoDB", scs.ConvoDB)
	printMap("UserConvosDB", scs.UserConvosDB)
}

// il nome è self explanatory brotha
func addConvoToUserIfNotExists(username string, convoToAdd *scs.ConversationELT) {
	// se la lista convo non esiste, la aggiungo (ma non succede mai)
	if _, exists := scs.UserConvosDB[username]; !exists {
		scs.UserConvosDB[username] = scs.Conversations{}
	}

	// se convo già in db, adios
	isAlreadyPresent := false
	for _, existingConvo := range scs.UserConvosDB[username] {
		if existingConvo.ConvoID == convoToAdd.ConvoID {
			isAlreadyPresent = true
			break
		}
	}

	// se no add brotha
	if !isAlreadyPresent {
		scs.UserConvosDB[username] = append(scs.UserConvosDB[username], convoToAdd)
	}
}

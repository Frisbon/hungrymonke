package structures

import (
	"math/rand"
	"strings"
	"time"
)

/*

NB: Gli adders aggiungono un elemento NUOVO, sse non esiste un ID già esistente!
	ALTRIMENTI CREANO UN'ALTRO ID E SALVANO COMUNQUE
	NON aggiornano il database con ID già presenti!!!

*/

func GenerateRandomString(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randGen := rand.New(rand.NewSource(time.Now().UnixNano())) // Nuovo generatore di numeri casuali
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(letters[randGen.Intn(len(letters))]) // Sceglie un carattere casuale dal set
	}
	return sb.String()
}

func AddToPrivateDB(ID string, toAdd Private) {

	for {

		if isUniversalIdUnique(ID) {
			PrivateDB[ID] = toAdd
			return
		}

		ID = GenerateRandomString(6)
	}

}

func AddToGroupDB(ID string, toAdd Group) {

	for {

		if isUniversalIdUnique(ID) {
			GroupDB[ID] = toAdd
			return
		}

		ID = GenerateRandomString(6)
	}

}

func AddToMsgDB(ID string, toAdd Message) {

	for {

		if isUniversalIdUnique(ID) {
			MsgDB[ID] = toAdd
			return
		}

		ID = GenerateRandomString(6)
	}

}

func AddToConvoDB(ID string, toAdd ConversationELT) {

	for {

		if isUniversalIdUnique(ID) {
			ConvoDB[ID] = toAdd
			return
		}

		ID = GenerateRandomString(6)
	}

}

// gli passo tutti i DB che usano generic ID e controllo se si ha uno li dentro
func isUniversalIdUnique(ID string) bool {

	exists := false

	if _, exists = PrivateDB[ID]; exists {
		return false
	}
	if _, exists = GroupDB[ID]; exists {
		return false
	}
	if _, exists = MsgDB[ID]; exists {
		return false
	}
	if _, exists = ConvoDB[ID]; exists {
		return false
	}

	return true

}

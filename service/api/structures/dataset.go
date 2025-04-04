package structures

import (
	"os"
	"time"
)

/*
In questo file creo dei record di utenti e messaggi
in modo da non dover ri-registrarli per testare l'applicazione.

Una rappresentazione grafica dei contenuti è stampata su dataset.png
*/
var DATASET_INITIALIZED = false

func init() {

	// load the pfps
	pfp, errpfp := os.ReadFile("service/api/pictureslol/cat2 Arturo.jpg")
	if errpfp != nil {
		panic("Failed to read image " + errpfp.Error()) // Panic during init if the file can't be read
	}

	// load the pfps
	pfp2, errpfp2 := os.ReadFile("service/api/pictureslol/Betta.png")
	if errpfp2 != nil {
		panic("Failed to read image " + errpfp2.Error()) // Panic during init if the file can't be read
	}

	// load the pfps
	pfp3, errpfp3 := os.ReadFile("service/api/pictureslol/Carlo.jpg")
	if errpfp3 != nil {
		panic("Failed to read image " + errpfp3.Error()) // Panic during init if the file can't be read
	}

	// Initialize users
	arturo := User{Username: "Arturo", Photo: pfp}
	betta := User{Username: "Betta", Photo: pfp2}
	carlo := User{Username: "Carlo", Photo: pfp3}

	// arturo := User{Username: "arturo", Photo: []byte{}}
	// betta := User{Username: "betta", Photo: []byte{}}
	// carlo := User{Username: "carlo", Photo: []byte{}}

	// Add users to UserDB
	UserDB["Arturo"] = &arturo
	UserDB["Betta"] = &betta
	UserDB["Carlo"] = &carlo

	// **Group Conversation: "FRENS*"**
	// Messages for the group conversation
	msg1 := Message{
		Timestamp: time.Date(2023, 10, 10, 12, 30, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("Sposarsi in autostrada, sì o no?")},
		Author:    &carlo,
		Status:    Seen,
		Reactions: []Reaction{
			{Author: &arturo, Emoticon: "👍", Timestamp: time.Date(2023, 10, 10, 12, 30, 0, 0, time.UTC)},
			{Author: &betta, Emoticon: "👎", Timestamp: time.Date(2023, 10, 10, 12, 30, 0, 0, time.UTC)},
		},
		MsgID:  "msg1",
		SeenBy: []*User{&betta, &arturo},
	}
	msg2 := Message{
		Timestamp: time.Date(2023, 10, 10, 12, 31, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("Carlo... hai preso le pillole oggi?")},
		Author:    &betta,
		Status:    Seen,
		Reactions: []Reaction{},
		MsgID:     "msg2",
		SeenBy:    []*User{&carlo, &arturo},
	}
	msg3 := Message{
		Timestamp: time.Date(2023, 10, 10, 12, 32, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("Eccolo di nuovo con i sintomi schizofrenici")},
		Author:    &arturo,
		Status:    Seen,
		Reactions: []Reaction{
			{Author: &carlo, Emoticon: "😂", Timestamp: time.Date(2023, 10, 10, 12, 32, 0, 0, time.UTC)},
			{Author: &betta, Emoticon: "😂", Timestamp: time.Date(2023, 10, 10, 12, 32, 0, 0, time.UTC)},
		},
		MsgID:  "msg3",
		SeenBy: []*User{&betta, &carlo},
	}

	// Load the photo file for msg4
	photoData, err := os.ReadFile("service/api/pictureslol/cat.jpg") // Adjust the path to match your file location
	if err != nil {
		panic("Failed to read image " + err.Error()) // Panic during init if the file can't be read
	}

	msg4 := Message{
		Timestamp: time.Date(2023, 10, 10, 12, 33, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("Vi mando il mio gatto"), Photo: &photoData, PhotoMimeType: "image/jpg"}, // todo add photo
		Author:    &arturo,
		Status:    Seen,
		Reactions: []Reaction{},
		MsgID:     "msg4",
		SeenBy:    []*User{&betta, &carlo},
	}

	// Initialize group conversation "FRENS*"
	groupFrens := ConversationELT{
		ConvoID:         "group_frens",
		DateLastMessage: msg4.Timestamp,
		Preview:         "[📷 Photo] Vi mando il mio...",
		Messages:        []*Message{&msg1, &msg2, &msg3, &msg4},
	}
	ConvoDB["group_frens"] = &groupFrens

	// load the group photo
	pfpg, errpfpg := os.ReadFile("service/api/pictureslol/cat frens.jpg")
	if errpfpg != nil {
		panic("Failed to read image " + errpfpg.Error()) // Panic during init if the file can't be read
	}

	// Initialize group
	group := Group{
		Conversation: &groupFrens,
		GroupPhoto:   pfpg,
		Name:         "FRENS 😼",
		Users:        []*User{&arturo, &betta, &carlo},
	}
	GroupDB["group_frens"] = &group

	// **Private Conversation: Arturo and Betta**
	// Messages for Arturo and Betta
	msg5 := Message{
		Timestamp: time.Date(2023, 10, 10, 16, 0, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("Ciao Betta, come stai?")},
		Author:    &arturo,
		Status:    Seen,
		Reactions: []Reaction{},
		MsgID:     "msg5",
	}
	msg6 := Message{
		Timestamp: time.Date(2023, 10, 10, 16, 1, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("Ciao Arturo spacca muro, benone, te?")},
		Author:    &betta,
		Status:    Seen,
		Reactions: []Reaction{},
		MsgID:     "msg6",
	}
	msg7 := Message{
		Timestamp: time.Date(2023, 10, 10, 16, 2, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("Bene, ti ricordi che usciamo con Carlo domani sera?")},
		Author:    &arturo,
		Status:    Seen,
		Reactions: []Reaction{
			{Author: &betta, Emoticon: "👍", Timestamp: time.Date(2023, 10, 10, 16, 2, 0, 0, time.UTC)},
		},
		MsgID: "msg7",
	}

	// Load the photo file for msg8 and 11
	photoData2, err := os.ReadFile("service/api/pictureslol/ritrovo.png") // Adjust the path to match your file location
	if err != nil {
		panic("Failed to read image " + err.Error()) // Panic during init if the file can't be read
	}

	msg8 := Message{
		Timestamp: time.Date(2023, 10, 10, 16, 2, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("Vediamoci alle 18:00 in Via Obelisco 3, in questo punto qui:"), Photo: &photoData2, PhotoMimeType: "image/png"},
		Author:    &arturo,
		Status:    Seen,
		Reactions: []Reaction{},
		MsgID:     "msg8",
	}
	msg9 := Message{
		Timestamp: time.Date(2023, 10, 10, 16, 3, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("Lui non ne sa ancora nulla, ora gli inoltro le indicazioni...")},
		Author:    &betta,
		Status:    Delivered,
		Reactions: []Reaction{},
		MsgID:     "msg9",
	}

	// Initialize private conversation between Arturo and Betta
	privateArturoBetta := ConversationELT{ // TODO: Siccome il pov varia tra betta e arturo, se last msg è di autore, compare anche status in UI
		ConvoID:         "private_arturo_betta",
		DateLastMessage: msg9.Timestamp,
		Preview:         "Lui non ne sa niente...",
		Messages:        []*Message{&msg5, &msg6, &msg7, &msg8, &msg9},
	}

	ConvoDB["private_arturo_betta"] = &privateArturoBetta

	// Initialize private chat
	privateAB := Private{
		Conversation: &privateArturoBetta,
		FirstUser:    &arturo,
		SecondUser:   &betta,
	}
	PrivateDB["private_arturo_betta"] = &privateAB

	msg10 := Message{
		Timestamp: time.Date(2023, 10, 10, 16, 4, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("Oi Carlè, ti lascio le info per l'uscita di domani")},
		Author:    &betta,
		Status:    Seen,
		Reactions: []Reaction{},
		MsgID:     "msg10",
	}
	msg11 := Message{
		Timestamp: time.Date(2023, 10, 10, 16, 5, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("Vediamoci alle 18:00 in Via Obelisco 3, in questo punto qui:"), Photo: &photoData2, PhotoMimeType: "image/png"},
		Author:    &betta,
		Status:    Seen,
		Reactions: []Reaction{},
		MsgID:     "msg11",
	}
	msg12 := Message{
		Timestamp: time.Date(2023, 10, 10, 16, 6, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("Grazie donna 😘")},
		Author:    &carlo,
		Status:    Seen,
		Reactions: []Reaction{
			{Author: &betta, Emoticon: "✨", Timestamp: time.Date(2023, 10, 10, 16, 6, 0, 0, time.UTC)},
		},
		MsgID: "msg12",
	}

	// Initialize private conversation between Betta and Carlo
	privateBettaCarlo := ConversationELT{
		ConvoID:         "private_betta_carlo",
		DateLastMessage: msg12.Timestamp,
		Preview:         "Grazie donna 😘",
		Messages:        []*Message{&msg10, &msg11, &msg12},
	}

	ConvoDB["private_betta_carlo"] = &privateBettaCarlo

	// Initialize private chat
	privateBC := Private{
		Conversation: &privateBettaCarlo,
		FirstUser:    &betta,
		SecondUser:   &carlo,
	}
	PrivateDB["private_betta_carlo"] = &privateBC

	// **Private Conversation: Arturo and Carlo**
	msg13 := Message{
		Timestamp: time.Date(2023, 10, 10, 13, 30, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("A fra guarda che bona")},
		Author:    &arturo,
		Status:    Seen,
		Reactions: []Reaction{},
		MsgID:     "msg13",
	}

	// Load the photo file for msg14
	photoData3, err := os.ReadFile("service/api/pictureslol/sorella.png") // Adjust the path to match your file location
	if err != nil {
		panic("Failed to read image " + err.Error()) // Panic during init if the file can't be read
	}

	msg14 := Message{
		Timestamp: time.Date(2023, 10, 10, 13, 31, 0, 0, time.UTC),
		Content:   Content{Photo: &photoData3, PhotoMimeType: "image/png"},
		Author:    &arturo,
		Status:    Seen,
		Reactions: []Reaction{},
		MsgID:     "msg14",
	}
	msg15 := Message{
		Timestamp: time.Date(2023, 10, 10, 13, 32, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("ao sta bono zi")},
		Author:    &carlo,
		Status:    Seen,
		Reactions: []Reaction{},
		MsgID:     "msg15",
	}
	msg16 := Message{
		Timestamp: time.Date(2023, 10, 10, 13, 33, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("è mi sorella")},
		Author:    &carlo,
		Status:    Seen,
		Reactions: []Reaction{},
		MsgID:     "msg16",
	}
	msg17 := Message{
		Timestamp: time.Date(2023, 10, 10, 13, 34, 0, 0, time.UTC),
		Content:   Content{Text: stringPtr("fammela conoscere 💀😂")},
		Author:    &carlo,
		Status:    Seen,
		Reactions: []Reaction{
			{Author: &betta, Emoticon: "🤦", Timestamp: time.Date(2023, 10, 10, 13, 35, 0, 0, time.UTC)},
		},
		MsgID: "msg17",
	}

	// Initialize private conversation between Arturo and Carlo
	privateArturoCarlo := ConversationELT{
		ConvoID:         "private_arturo_carlo",
		DateLastMessage: msg17.Timestamp,
		Preview:         "fammela conoscere 💀😂",
		Messages:        []*Message{&msg13, &msg14, &msg15, &msg16, &msg17},
	}
	ConvoDB["private_arturo_carlo"] = &privateArturoCarlo

	// Initialize private chat
	privateAC := Private{
		Conversation: &privateArturoCarlo,
		FirstUser:    &arturo,
		SecondUser:   &carlo,
	}
	PrivateDB["private_arturo_carlo"] = &privateAC

	// Add all messages to MsgDB
	messages := []*Message{&msg1, &msg2, &msg3, &msg4, &msg5, &msg6, &msg7, &msg8, &msg9, &msg10, &msg11, &msg12, &msg13, &msg14, &msg15, &msg16, &msg17}
	for _, msg := range messages {
		MsgDB[msg.MsgID] = msg
	}

	// Initialize UserConvosDB
	UserConvosDB["Arturo"] = Conversations{&groupFrens, &privateArturoBetta, &privateArturoCarlo}
	UserConvosDB["Betta"] = Conversations{&groupFrens, &privateArturoBetta, &privateBettaCarlo}
	UserConvosDB["Carlo"] = Conversations{&groupFrens, &privateBettaCarlo, &privateArturoCarlo}

	// Add all IDs to GenericDB to ensure uniqueness
	ids := []string{"group_frens", "private_arturo_betta", "private_betta_carlo", "private_arturo_carlo", "msg1", "msg2", "msg3", "msg4", "msg5", "msg6", "msg7", "msg8", "msg9", "msg10", "msg11", "msg12", "msg13", "msg14", "msg15", "msg16", "msg17"}
	for _, id := range ids {
		GenericDB[id] = struct{}{}
	}

	DATASET_INITIALIZED = true

}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

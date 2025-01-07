package structures

import "time"

type User struct {
	Username string `json:"username"`
	Photo    []byte `json:"photo"`
}

type Status string // enumerativo per Message
const (
	Delivered Status = "delivered"
	Recieved  Status = "recieved"
	Seen      Status = "seen"
)

type Reaction struct {
	ID       string `json:"id"`
	Author   string `json:"author"`
	Emoticon string `json:"emoticon"`
}

type Content struct {
	Text  *string `json:"text,omitempty"`  // Usa puntatori per distinguere i valori non forniti
	Photo *[]byte `json:"photo,omitempty"` // Alternativa per il contenuto
}

type Message struct {
	Timestamp time.Time `json:"timestamp"`
	/* NB Content type is either Text o Photo*/
	Content   Content    `json:"content"`
	Author    string     `json:"author"`
	Status    Status     `json:"status"`
	Reactions []Reaction `json:"reactions"`
	ID        string     `json:"id"` //id messaggio
}

/*La struct MessagePreviewXorPhoto era inutile da dichiarare (vedi Preview sotto.) */
type ConversationELT struct {
	ID              string    `json:"id"`
	DateLastMessage time.Time `json:"datelastmessage"` //timestamp
	/*NB il Preview è una stringa variabile (messaggio) or una stringa prefissata ("📷 Photo")*/
	Preview  string    `json:"preview"`
	Messages []Message `json:"messages"`
}

type Conversations []ConversationELT

type Group struct {
	Conversation ConversationELT `json:"conversation"`
	GroupPhoto   []byte          `json:"groupphoto"`
	Name         string          `json:"name"`
	Users        []User          `json:"users"`
}

type Private struct {
	Conversation ConversationELT `json:"conversation"`
	FirstUser    User            `json:"firstuser"`
	SecondUser   User            `json:"seconduser"`
}

package structures

import "time"

type Photo struct {
	PhotoFile []byte `json:"photofile"` // a quanto pare []byte si usa per i file binari
}

type UserID struct {
	UserID string `json:"userid"`
}

type GenericID struct {
	GenericID string `json:"genericid"`
}

type User struct {
	Username UserID `json:"username"`
	Photo    Photo  `json:"photo"`
}

type Status string // enumerativo per Message
const (
	Recieved Status = "recieved"
	Seen     Status = "seen"
)

type Reaction struct {
	ID       GenericID `json:"id"`
	Author   UserID    `json:"author"`
	Emoticon string    `json:"emoticon"`
}

type Text string
type Message struct {
	Timestamp time.Time `json:"timestamp"`
	/* NB Content type is either Text o Photo*/
	Content   interface{} `json:"content"` // todo ricordati il type assertion dopo altrimenti può essere any type
	Type      string      `json:"type"`
	Author    UserID      `json:"author"`
	Status    Status      `json:"status"`
	Reactions []Reaction  `json:"reactions"`
	ID        GenericID   `json:"id"`
}

/*La struct MessagePreviewXorPhoto era inutile da dichiarare (vedi Preview sotto.) */
type ConversationELT struct {
	ID              GenericID `json:"id"`
	Photo           Photo     `json:"photo"`
	DateLastMessage time.Time `json:"datelastmessage"` //timestamp
	/*NB il Preview è una stringa variabile (messaggio) or una stringa prefissata ("📷 Photo")*/
	Preview  string    `json:"preview"`
	Messages []Message `json:"messages"`
}

type Conversations []ConversationELT

type Group struct {
	Conversation ConversationELT `json:"conversation"`
	Name         string          `json:"name"`
	Users        []User          `json:"users"`
}

type Private struct {
	Conversation ConversationELT `json:"conversation"`
	FirstUser    User            `json:"firstuser"`
	SecondUser   User            `json:"seconduser"`
}

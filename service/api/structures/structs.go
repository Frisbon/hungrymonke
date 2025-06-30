package structures

import "time"

/*
	NB: Ho ignorato il fatto che GenericID fosse readonly, per semplificare il codice GO
		sticavoli se viene modificato... basta che non ci siano doppioni, tutto qui.
*/

type User struct {
	Username      string `json:"username"`
	Photo         []byte `json:"photo"`
	PhotoMimeType string `json:"photoMimeType,omitempty"`
}

type Status string //  enumerativo per Message
const (
	Delivered Status = "delivered"
	Seen      Status = "seen"
)

type Reaction struct {
	Author    *User     `json:"author"`
	Emoticon  string    `json:"emoticon"`
	Timestamp time.Time `json:"timestamp"`
}

type Content struct {
	Text          *string `json:"text,omitempty"`  //  Uso i puntatori qui per distinguere i valori non forniti
	Photo         *[]byte `json:"photo,omitempty"` //  Alternativa per il contenuto
	PhotoMimeType string  `json:"photoMimeType,omitempty"`
}

type Message struct {
	Timestamp time.Time `json:"timestamp"`
	/* NB Content type is either Text o Photo or both*/
	Content     Content    `json:"content"`
	Author      *User      `json:"author"`
	Status      Status     `json:"status"`
	Reactions   []Reaction `json:"reactions"`
	MsgID       string     `json:"msgid"`
	SeenBy      []*User    `json:"seenby"` //  new record for group msgs.
	IsForwarded bool       `json:"isforwarded,omitempty"`
	ReplyingTo  *Message   `json:"replyingto,omitempty"`
}

/*La struct MessagePreviewXorPhoto era inutile da dichiarare (vedi Preview sotto.) */
type ConversationELT struct {
	ConvoID         string     `json:"convoid"`         //  id conversazione // readonly, implemento un costruttore.
	DateLastMessage time.Time  `json:"datelastmessage"` // timestamp
	Preview         string     `json:"preview"`         /*NB il Preview è una stringa variabile (messaggio) or una stringa prefissata ("📷 Photo") v2: oppure un mix tra i due? :)*/
	Messages        []*Message `json:"messages"`
}

type Conversations []*ConversationELT

type Group struct {
	Conversation  *ConversationELT `json:"conversation"`
	GroupPhoto    []byte           `json:"groupphoto"`
	PhotoMimeType string           `json:"photoMimeType,omitempty"`
	Name          string           `json:"name"`
	Users         []*User          `json:"users"`
}

// Questa struct serve per aiutare a decidere lo status "Seen" in una convo di gruppo.
type GroupMessage struct {
	Msg   Message `json:"message"`
	Users []*User `json:"users"`
}

type Private struct {
	Conversation *ConversationELT `json:"conversation"`
	FirstUser    *User            `json:"firstuser"`
	SecondUser   *User            `json:"seconduser"`
}

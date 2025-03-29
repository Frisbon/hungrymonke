package handlers

import (
	"time"

	scs "github.com/Frisbon/hungrymonke/service/api/structures"
)

/*
	Siccome ci sono struct (MESSAGE, CONVERSATIONELT) con ID univoci, per evitare di aggiornare manualmente ogni volta i database,
	creo dei costruttori che generano un ID univoco e mi aggiornano i DB automaticamente ogni volta
	che creo una struttura.

*/

// STRUCT SENZA ID PASSATAMI DALL'ESTERNO
type MsgCONSTR struct {
	Timestamp time.Time      `json:"timestamp"`
	Content   scs.Content    `json:"content"`
	Author    *scs.User      `json:"author"` // sender_struct pointer
	Status    scs.Status     `json:"status"`
	Reactions []scs.Reaction `json:"reactions"`
}

// Costruisce il messaggio, lo salva nel DB e ritorna il puntatore a come argomento.
func ConstrMessage(data MsgCONSTR) *scs.Message {

	m := &scs.Message{

		Timestamp: data.Timestamp,
		Content:   data.Content,
		Author:    data.Author,
		Status:    data.Status,
		Reactions: data.Reactions,
		MsgID:     GenerateRandomString(5),
	}

	scs.MsgDB[m.MsgID] = m
	return m

}

// STRUCT SENZA ID PASSATAMI DALL'ESTERNO
type ConvoCONSTR struct {
	DateLastMessage time.Time      `json:"datelastmessage"` //timestamp
	Preview         string         `json:"preview"`         /*NB il Preview è una stringa variabile (messaggio) or una stringa prefissata ("📷 Photo") v2: oppure un mix tra i due? :)*/
	Messages        []*scs.Message `json:"messages"`
}

// Costruisce la convo, la salva nel DB e la ritorna come argomento.
func ConstrConvo(data ConvoCONSTR) *scs.ConversationELT {

	c := &scs.ConversationELT{

		DateLastMessage: data.DateLastMessage,
		Preview:         data.Preview,
		Messages:        data.Messages,
		ConvoID:         GenerateRandomString(5),
	}

	scs.ConvoDB[c.ConvoID] = c
	return c

}

func ConstrGroup(Users []*scs.User) *scs.Group {

	g := &scs.Group{

		Conversation: ConstrConvo(ConvoCONSTR{}),
		Name:         "NewGroup :)",
		Users:        Users,
	}

	scs.GroupDB[g.Conversation.ConvoID] = g
	return g

}

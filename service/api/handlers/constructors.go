package handlers

import (
	"time"

	scs "github.com/Frisbon/hungrymonke/service/api/structures"
)

type MsgCONSTR struct {
	Timestamp   time.Time
	Content     scs.Content
	Author      *scs.User
	Status      scs.Status
	Reactions   []scs.Reaction
	SeenBy      []scs.User
	IsForwarded bool
	ReplyingTo  *scs.Message
}

func ConstrMessage(data MsgCONSTR) *scs.Message {
	m := &scs.Message{
		Timestamp:   data.Timestamp,
		Content:     data.Content,
		Author:      data.Author,
		Status:      data.Status,
		Reactions:   data.Reactions,
		ReplyingTo:  data.ReplyingTo,
		IsForwarded: data.IsForwarded,
		MsgID:       GenerateRandomString(5),
	}
	scs.MsgDB[m.MsgID] = m
	return m
}

type ConvoCONSTR struct {
	DateLastMessage time.Time
	Preview         string
	Messages        []*scs.Message
}

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

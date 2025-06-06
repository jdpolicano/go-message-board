package db

import (
	"errors"
	"time"
)

var NullPointer = errors.New("null pointer reciever")
var SessionAlreadyExists = errors.New("board already exists")
var NoSuchSession = errors.New("board doesn't exist")

type Message struct {
	User    string `json:"user"`    // user that posted the message.
	Time    int64  `json:"time"`    // unix timestamp of the message.
	Content string `json:"content"` // the actual message written.
}

type Session struct {
	Id       string // the name of this Session.
	Creator  string
	Messages []Message // messages in the board, sorted by time when they occured.
}

type MemoryDatabase struct {
	Sessions map[string]*Session // a list of available boards by name
}

func NewMemDatabase() MemoryDatabase {
	return MemoryDatabase{Sessions: make(map[string]*Session)}
}

func (db *MemoryDatabase) CreateSession(name string, id string) (*Session, error) {
	if db == nil {
		return nil, NullPointer
	}

	if db.GetSession(name) != nil {
		return nil, SessionAlreadyExists
	}

	session := &Session{name, id, make([]Message, 0)}
	db.Sessions[name] = session
	return session, nil
}

func (db *MemoryDatabase) GetSession(id string) *Session {
	if db == nil {
		return nil
	}
	return db.Sessions[id]
}

func (db *MemoryDatabase) GetSessionIds() []string {
	if db == nil {
		return nil
	}
	ids := make([]string, len(db.Sessions), 0)
	for id := range db.Sessions {
		ids = append(ids, id)
	}
	return ids
}

func (session *Session) AddMessage(content, user string) Message {
	message := Message{
		User:    user,
		Time:    time.Now().UnixMilli(),
		Content: content,
	}
	session.Messages = append(session.Messages, message)
	return message
}

func (session *Session) GetMessages() []Message {
	messages := make([]Message, 0, len(session.Messages))
	copy(messages, session.Messages)
	return messages
}

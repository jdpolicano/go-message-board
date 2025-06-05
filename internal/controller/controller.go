package controller

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jdpolicano/go-message-board/internal/db"
)

type ClientMessage struct {
	Username string
	Content  string
}

type Client struct {
	conn     *websocket.Conn
	userName string
}

func NewClient(c *websocket.Conn, name string) *Client {
	return &Client{c, name}
}

type SessionHandle struct {
	register   chan *Client
	unregister chan string
	message    chan ClientMessage
	sigterm    chan bool
	shouldQuit bool
	store      *db.Session
	clients    []*Client
}

func newSessionHandle(store *db.Session) *SessionHandle {
	return &SessionHandle{
		register:   make(chan *Client),
		unregister: make(chan string),
		message:    make(chan ClientMessage),
		sigterm:    make(chan bool),
		shouldQuit: false,
		clients:    make([]*Client, 0, 1024),
		store:      store,
	}
}

func (s *SessionHandle) spawn() {
	defer s.closeSession()
	for !s.shouldQuit {
		select {
		case client := <-s.register:
			{
				s.addClient(client)
			}
		case username := <-s.unregister:
			{
				s.removeClient(username)
			}
		case message := <-s.message:
			{
				s.addMessage(message)
			}
		case <-s.sigterm:
			{
				break
			}
		}
	}
}

func (s *SessionHandle) addClient(client *Client) {
	if s == nil {
		return
	}
	for _, c := range s.clients {
		if c.userName == client.userName {
			return
		}
	}
	s.clients = append(s.clients, client)
}

func (s *SessionHandle) removeClient(username string) {
	if s == nil {
		return
	}

	if s.store.Creator == username {
		s.shouldQuit = true
		return
	}

	start, end := 0, len(s.clients)-1
	for start <= end {
		if s.clients[start].userName == username {
			if start != end {
				s.clients[start], s.clients[end] = s.clients[end], s.clients[start]
			}
			break
		}
		start++
	}
	s.clients = s.clients[:end]
}

func (s *SessionHandle) addMessage(message ClientMessage) {
	if s == nil {
		return
	}
	payload := s.store.AddMessage(message.Content, message.Username)
	for _, c := range s.clients {
		if c.userName != message.Username {
			c.conn.WriteJSON(payload)
		}
	}
}

func (s *SessionHandle) closeSession() {
	if s == nil {
		return
	}
	for _, c := range s.clients {
		c.conn.Close()
	}
}

type Controller struct {
	sync.RWMutex
	dataHandle *db.MemoryDatabase
	sessions   map[string]*SessionHandle
}

func NewController(db *db.MemoryDatabase) *Controller {
	return &Controller{
		sync.RWMutex{},
		db,
		make(map[string]*SessionHandle),
	}
}

func (c *Controller) ListSessionIds() []string {
	c.RLock()
	defer c.RUnlock()
	sessions := make([]string, 0, len(c.sessions))
	for id := range c.sessions {
		sessions = append(sessions, id)
	}
	return sessions
}

func (c *Controller) CreateSession(client *Client) (*SessionHandle, error) {
	c.RLock()
	defer c.RUnlock()
	id := uuid.NewString()
	store, e := c.dataHandle.CreateSession(client.userName, id)
	if e != nil {
		return nil, fmt.Errorf("%s", e)
	}
	session := newSessionHandle(store)
	go session.spawn()
	return session, nil
}

// func SpawnBoardChat(db *db.Session) SessionHandle {
// 	handle := SessionHandle{
// 		register:   make(chan *Client),
// 		unregister: make(chan string),
// 		message:    make(chan ClientMessage),
// 		sigterm:    make(chan bool),
// 	}

// 	clients := make([]*Client, 0, 1024)

// 	go func() {
// 		for {
// 			select {
// 			case msg := <-handle.message:
// 				{
// 					final := db.AddMessage(msg.Content, msg.Name)
// 					payload, e := json.Marshal(final)
// 					if e == nil {
// 						for _, c := range clients {
// 							if c.userName != final.User {
// 								c.conn.WriteMessage(websocket.BinaryMessage, payload)
// 							}
// 						}
// 					} else {
// 						fmt.Println(db.Name, e)
// 					}
// 				}
// 			case newClient := <-handle.register:
// 				{
// 					exists := false
// 					for _, c := range clients {
// 						if c.userName == newClient.userName {
// 							exists = true
// 						}
// 					}
// 					if !exists {
// 						clients = append(clients, newClient)
// 					}
// 				}
// 			case unregName := <-handle.unregister:
// 				{
// 					curr, end := 0, len(clients)-1
// 					for curr <= end {
// 						if clients[curr].userName == unregName {
// 							clients[curr].conn.Close()
// 							clients[curr], clients[end] = clients[end], clients[curr]
// 							clients = clients[0:end]
// 							end--
// 							continue
// 						}
// 						curr++
// 					}
// 				}
// 			case <-handle.sigterm:
// 				{
// 					for _, c := range clients {
// 						c.conn.Close()
// 					}
// 					return
// 				}
// 			}

// 		}
// 	}()

// 	return handle
// }

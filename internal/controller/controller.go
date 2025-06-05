package controller

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jdpolicano/go-message-board/internal/db"
)

type ClientMessage struct {
	Name    string
	Content string
}

type Client struct {
	conn         *websocket.Conn
	userName     string
	ownedSession *SessionHandle
}

type SessionHandle struct {
	register   chan *Client
	unregister chan string
	message    chan ClientMessage
	sigterm    chan bool
	clients    []*Client
	creator    *Client
	id         string
}

type Controller struct {
	dataHandle    *db.MemoryDatabase
	sessions      map[string]*SessionHandle
	createSession chan *Client
	joinSession   chan *SessionHandle
	listSession   chan *[]string
	quit          chan bool
}

func NewController(db *db.MemoryDatabase) *Controller {
	return &Controller{
		db,
		make(map[string]*SessionHandle),
		make(chan *Client),
		make(chan *SessionHandle),
		make(chan *[]string),
		make(chan bool),
	}
}

func newSessionHandle(creator *Client, id string) *SessionHandle {
	return &SessionHandle{
		register:   make(chan *Client),
		unregister: make(chan string),
		message:    make(chan ClientMessage),
		sigterm:    make(chan bool),
		clients:    make([]*Client, 0, 1024),
		creator:    creator,
		id:         id,
	}
}

func NewClient(c *websocket.Conn, name string) *Client {
	return &Client{c, name, nil}
}

func (controller *Controller) ListSessionIds() *[]string {
	sessions := make([]string, 0, 1024)
	controller.listSession <- &sessions
	return &sessions
}

func (controller *Controller) CreateNewSession(client *Client) {
	controller.createSession <- client
}

func (controller *Controller) Spawn() {
	go controller.spawnHandler()
}

func (controller *Controller) spawnHandler() {
	for {
		select {
		case container := <-controller.listSession:
			{
				controller.listSessions(container)
			}
		case client := <-controller.createSession:
			{
				controller.spawnSession(client)
			}
		}
		case handle := <-controller.joinSession: {

		}
	}
}

func (controller *Controller) listSessions(container *[]string) {
	for sess := range controller.sessions {
		*container = append(*container, sess)
	}
}

func (controller *Controller) spawnSession(client *Client) {
	id := uuid.NewString()
	sessionHandle := newSessionHandle(client, id)
	controller.sessions[id] = sessionHandle
	client.ownedSession = sessionHandle
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

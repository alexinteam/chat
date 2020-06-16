// Serving client websocket connecions.

package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	pongWait       = 3 * time.Second
	pingPeriod     = 2 * time.Second
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	user    *User
	message chan *Message
}

// End client's session
func (c *Client) kill(msg *Message) {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		log.Println("JSON encoding error:", err)
		return
	}

	err = c.conn.WriteMessage(
		websocket.TextMessage,
		msgJson,
	)
	if err != nil {
		log.Println("Bye-message write error:", err)
		return
	}

	c.conn.Close()
}

func (c *Client) readWS() {
	defer func() {
		un := &Unreg{
			client: c,
			msg:    c.user.Username + " has gone",
		}
		c.hub.unregister <- un
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("Read error:", err)
			}
			return
		}

		var msg Message
		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Println("JSON decode error:", err)
			return
		}

		msg.Room = c.hub.room

		switch msg.Action {
		case "message":
			msg.Sender = c.user
			c.hub.message <- &msg
		}
	}
}

func (c *Client) writeWS() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		un := &Unreg{
			client: c,
			msg:    c.user.Username + " has gone",
		}
		c.hub.unregister <- un
	}()

	for {
		select {
		case msg := <-c.message:
			msgJson, err := json.Marshal(msg)
			if err != nil {
				log.Println("JSON encoding error:", err)
				continue
			}

			err = c.conn.WriteMessage(
				websocket.TextMessage,
				msgJson,
			)
			if err != nil {
				log.Println("Message write error:", err)
				return
			}

		// Heartbeat
		case <-ticker.C:
			err := c.conn.WriteMessage(
				websocket.PingMessage,
				[]byte{},
			)
			if err != nil {
				return
			}
		}
	}
}

func handlerWS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roomId, _ := strconv.Atoi(vars["id"])
	hub, ok := hubs[roomId]
	if !ok {
		log.Println("No hub for room ID" + string(roomId))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Open connection error:", err)
		return
	}
	log.Println("Successfully connected")

	user := context.Get(r, "User").(*User)

	client := &Client{
		hub:     hub,
		conn:    conn,
		user:    user,
		message: make(chan *Message),
	}
	client.hub.register <- &Reg{client: client}

	go client.writeWS()
	client.readWS()
}

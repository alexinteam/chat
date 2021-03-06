// Hub for transfering messages

package chat

import (
	"log"
	"time"
)

type Hub struct {
	clients    map[*Client]bool
	room       *Room
	message    chan *Message
	ctl        chan *Message
	register   chan *Reg
	unregister chan *Unreg
}

type Reg struct {
	client *Client
}

type Unreg struct {
	client *Client
	msg    string
}

func makeHub(room *Room) *Hub {
	h := &Hub{
		clients:    make(map[*Client]bool),
		room:       room,
		message:    make(chan *Message),
		ctl:        make(chan *Message),
		register:   make(chan *Reg),
		unregister: make(chan *Unreg),
	}

	// Make reverse relation
	h.room.hub = h

	return h
}

// Send message to all clients
func (h *Hub) send(msg *Message) {
	for client := range h.clients {
		// These messages are sent to everyone including sender
		toAll := contains([]string{"mute", "ban"}, msg.Action)
		// Don't send messages to sender
		toSelf := msg.Sender != nil && client.user.Id == msg.Sender.Id
		// Send only to recipient or if it is broadcast
		isBroadcast := msg.Recipient == nil
		isRecipient := !isBroadcast && (client.user.Id == msg.Recipient.Id)

		doSend := toAll || (!toSelf && (isBroadcast || isRecipient))

		if doSend {
			client.message <- msg
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		// Add client to chat
		case reg := <-h.register:
			client := reg.client
			log.Println("Registered:", client.user.Username)
			h.clients[client] = true
			u := client.user

			// Tell everyone about new user
			msg := &Message{
				Action:   "new_user",
				Sender:   u,
				Text:     client.user.Username + " joined the room",
				SendDate: time.Now().UTC(),
			}
			h.send(msg)

		// Remove client from chat
		case unreg := <-h.unregister:
			client := unreg.client
			t := unreg.msg
			_, alive := h.clients[client]
			if alive {
				log.Println("Unregistered:", client.user.Username)

				// Tell everyone about user has gone
				msg := &Message{
					Action:   "gone_user",
					Sender:   client.user,
					Text:     t,
					SendDate: time.Now().UTC(),
				}
				h.send(msg)

				client.kill(msg)
				delete(h.clients, client)
			}

		// Information messages
		case ctl := <-h.ctl:
			h.send(ctl)

		// Chat messages: send to all and save to DB
		case msg := <-h.message:
			err := msg.save()
			if err != nil {
				log.Println("Saving message error:", err)
				continue
			}

			if msg.Recipient != nil {
				log.Printf(
					"%s TO %s: %s",
					msg.Sender.Username,
					msg.Recipient.Username,
					msg.Text,
				)
			} else {
				log.Printf(
					"%s: %s",
					msg.Sender.Username,
					msg.Text,
				)
			}

			h.send(msg)
		}
	}
}

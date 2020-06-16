// Structure and methods for Room objects.

package chat

import (
	"database/sql"
)

type Room struct {
	Id   int
	Name string
	hub  *Hub
}

func (r *Room) getUsers() []*User {
	var users []*User
	for c := range r.hub.clients {
		u := c.user
		users = append(users, u)
	}

	return users
}

func (r *Room) getMessages(user *User, limit int) ([]*Message, error) {
	var messages []*Message

	rows, err := stmtGetRoomMessagesByUser.Query(user.Id, r.Id, limit)
	if err == sql.ErrNoRows {
		return []*Message{}, nil
	} else if err != nil {
		return []*Message{}, err
	} else {
		defer rows.Close()
	}

	var msg *Message
	var sId, rId *int
	var sUsername, sFullname, sEmail *string
	var rUsername, rFullname, rEmail *string

	for rows.Next() {
		msg = &Message{}
		err = rows.Scan(
			// Message info
			&msg.Id, &msg.Action, &msg.Text, &msg.SendDate,
			// Sender
			&sId, &sUsername, &sFullname, &sEmail,
			// Recipient
			&rId, &rUsername, &rFullname, &rEmail,
		)
		if err != nil {
			return []*Message{}, err
		}
		if sId != nil {
			msg.Sender = &User{
				Id:       *sId,
				Username: *sUsername,
				Fullname: *sFullname,
				Email:    *sEmail,
			}
		}
		if rId != nil {
			msg.Recipient = &User{
				Id:       *rId,
				Username: *rUsername,
				Fullname: *rFullname,
				Email:    *rEmail,
			}
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func getAllRooms() ([]*Room, error) {
	var rooms []*Room

	rows, err := stmtGetAllRooms.Query()
	if err == sql.ErrNoRows {
		return []*Room{}, nil
	} else if err != nil {
		return []*Room{}, err
	} else {
		defer rows.Close()
	}

	var room *Room
	for rows.Next() {
		room = &Room{}
		err = rows.Scan(&room.Id, &room.Name)
		if err != nil {
			return []*Room{}, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// Structure and methods for User objects.

package chat

type User struct {
	Id       int    `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

// Terminate all connections of this user
func (u *User) logout() {
	stmtDeleteUserSessions.Exec(u.Id)

	un := &Unreg{
		msg: u.Username + " has gone due to another connection",
	}

	// Send unregister message to each hub with this user
	for _, h := range hubs {
		for c := range h.clients {
			if c.user.Id == u.Id {
				un.client = c
				c.hub.unregister <- un
				break
			}
		}
	}
}

func getUserById(id int) (*User, error) {
	var user User
	err := stmtGetUserById.QueryRow(id).Scan(
		&user.Id, &user.Fullname, &user.Username, &user.Email,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// URL handlers for pages.

package chat

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func handlerLoginPage(w http.ResponseWriter, r *http.Request) {
	ctx := make(map[string]string)

	if r.Method == "POST" {
		r.ParseForm()
		login := r.Form["login"][0]
		password := r.Form["password"][0]

		user, err := authenticate(login, password)
		if err != nil {
			ctx["err"] = err.Error()
		} else {
			// Kill all other connections and make new one
			// because only 1 session per user allowed
			user.logout()
			makeSession(w, user)
			http.Redirect(w, r, "/", 302)
			return
		}
	}

	tplLogin.ExecuteTemplate(w, "base", ctx)
}

func handlerLogout(w http.ResponseWriter, r *http.Request) {
	removeSession(w, r)
	http.Redirect(w, r, "/login", 302)
	return
}

func handlerChatPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u := context.Get(r, "User").(*User)

	// No need to check errors - regex in mux validates input
	roomId, _ := strconv.Atoi(vars["id"])

	room := hubs[roomId].room

	ctx := struct {
		User *User
		Room *Room
	}{
		u,
		room,
	}
	tplChat.ExecuteTemplate(w, "base", ctx)
}

func handlerIndexPage(w http.ResponseWriter, r *http.Request) {
	// Since this is the default handler for all URLs, we
	// must check if it is correct URL
	if !contains([]string{"/", "/chat"}, r.URL.Path) {
		http.Error(w, "Not found", 404)
		return
	}
	u := context.Get(r, "User").(*User)

	rooms, err := getAllRooms()
	if err != nil {
		log.Println("Getting rooms error:", err)
		return
	}

	ctx := struct {
		User  *User
		Rooms []*Room
	}{
		u,
		rooms,
	}
	tplIndex.ExecuteTemplate(w, "base", ctx)
}

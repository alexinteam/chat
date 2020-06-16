// Bind URLs to handle functions.

package chat

import (
	"github.com/gorilla/mux"
	"net/http"
)

func makeRouter() {
	r := mux.NewRouter().StrictSlash(true)

	staticDir := "/static/"

	r.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	// AJAX
	r.HandleFunc(
		"/ajax/rooms/{id:[0-9]+}/users",
		authMiddleware(handlerAjaxGetRoomUsers),
	)
	r.HandleFunc(
		"/ajax/rooms/{id:[0-9]+}/messages",
		authMiddleware(handlerAjaxGetRoomMessages),
	)

	// WS
	r.HandleFunc(
		"/ws/rooms/{id:[0-9]+}",
		authMiddleware(handlerWS),
	)

	// Pages
	r.HandleFunc(
		"/login",
		handlerLoginPage,
	)
	r.HandleFunc(
		"/logout",
		authMiddleware(handlerLogout),
	)
	r.HandleFunc(
		"/chat/rooms/{id:[0-9]+}",
		authMiddleware(handlerChatPage),
	)
	r.HandleFunc(
		"/chat",
		authMiddleware(handlerIndexPage),
	)
	r.HandleFunc(
		"/",
		authMiddleware(handlerIndexPage),
	)

	http.Handle("/", r)
}

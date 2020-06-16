// Entry point.

package chat

import (
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

// Global storage of hubs (one per room)
var hubs = make(map[int]*Hub)

func RunServer(settings map[string]string) {
	var err error

	db = dbConnect(
		settings["dbUser"],
		settings["dbPass"],
		settings["dbName"],
	)
	log.Println("DB connected successfully")
	defer db.Close()

	// Prepare SQL statements
	initStmts()

	// Parse HTML-templates
	initTpls()

	// Run websockets for each room
	rooms, err := getAllRooms()
	if err != nil {
		log.Fatal("Could not get the list of rooms")
	}
	for _, room := range rooms {
		hub := makeHub(room)
		hubs[room.Id] = hub
		go hub.run()
	}

	// MAke router
	makeRouter()

	// Run server
	log.Printf("Server is running on %s port...\n", settings["port"])
	err = http.ListenAndServe(":"+settings["port"], nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}

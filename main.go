// Entry point.

package main

import "github.com/alexinteam/chat/chat"

var settings = map[string]string{
	// app port
	"port": "8080",

	// Database settings
	"dbUser": "postgres",
	"dbPass": "postgres",
	"dbName": "chat",
}

func main() {
	chat.RunServer(settings)
}

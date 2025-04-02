package main

import (
	"github.com/VincNT21/kallaxy/client/context"
	"github.com/VincNT21/kallaxy/client/gui"
)

func main() {
	// Load configuration ?

	const serverURL = "http://localhost:8080"

	// Call NewAppContext to initialize dependencies
	apiCtx := context.NewAppContext(serverURL)

	// Pass the AppContext to StartGui
	gui.StartGui(apiCtx)
}

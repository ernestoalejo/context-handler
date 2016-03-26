package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/tylerb/graceful.v1"
)

// ServeGracefully creates a new listener that will close gracefully without dropping
// client connections. If a PORT environment variable is defined it will use that
// port instead of the default one. If debugging is enabled a normal listener will be
// used instead to close it quicker in the development cycle.
func ServeGracefully(defaultPort int64, stopTimeout time.Duration) {
	listenPort := fmt.Sprintf(":%d", defaultPort)
	if os.Getenv("PORT") != "" {
		listenPort = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	log.Printf("[*] Server listening in %s\n", listenPort)
	if os.Getenv("DEBUG") == "true" {
		http.ListenAndServe(listenPort, nil)
	} else {
		graceful.Run(listenPort, stopTimeout, nil)
	}

	log.Println("[*] Graceful server shutdown")
}

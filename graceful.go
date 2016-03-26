package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/tylerb/graceful.v1"
)

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

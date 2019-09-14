package crossref

import (
	log "github.com/sirupsen/logrus"
)

const api string = "https://api.crossref.org/works"

// Debug checks the communication between the library and the API in detail.
func Debug() {
	log.SetLevel(log.DebugLevel)
	log.Info("Changed to debugging mode")
}

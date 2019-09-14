package crossref

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Client communicates with CrossRef API and ensures politeness.
type Client struct {
	mailto  string
	appname string
	http.Client
}

// NewClient returns a new crossref client with an attached mailto and app
// name. The request timeout is set to 5 seconds.
func NewClient(appname, malito string) *Client {
	c := Client{
		mailto:  malito,
		appname: appname,
		Client:  http.Client{Timeout: 5 * time.Second},
	}

	log.WithFields(log.Fields{
		"client": c,
	}).Debug("New client created")

	return &c
}

// String implements the Stringer interface.
func (c *Client) String() string {
	return fmt.Sprintf("App: %s, MailTo: %s", c.appname, c.mailto)
}

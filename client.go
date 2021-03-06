package crossref

import (
	"fmt"
	"io"
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
		Client:  http.Client{Timeout: 30 * time.Second},
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

func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	r.Header.Add("User-Agent", fmt.Sprintf("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36 %s/1.0 (https://github.com/cgxeiji/crossref; mailto:%s)", c.appname, c.mailto))
	log.Debug(r.UserAgent())

	return r, nil
}

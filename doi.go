package crossref

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"
)

// DOIJSON returns the raw JSON of a DOI search. Use this if you want to
// manually process the JSON data.
func (c *Client) DOIJSON(doi string) ([]byte, error) {
	log.WithFields(log.Fields{
		"client": c,
	}).Debug("DOI requested")

	if doi == "" {
		return nil, ErrEmptyQuery
	}

	url := fmt.Sprintf("%s/%s", api, doi)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", fmt.Sprintf("%s (mailto: %s)", c.appname, c.mailto))

	q := req.URL.Query()
	q.Add("mailto", c.mailto)
	req.URL.RawQuery = q.Encode()

	log.WithFields(log.Fields{
		"url": req.URL.String(),
	}).Debug("Requesting information")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

// DOI gets a processed metadata from a DOI.
func (c *Client) DOI(doi string) (*Work, error) {
	js, err := c.DOIJSON(doi)
	if err != nil {
		return &Work{}, err
	}

	var data doiJSON
	json.Unmarshal(js, &data)
	log.Debug(pretty.Sprint(data))

	if data.Status != "ok" {
		return &Work{}, ErrZeroWorks
	}

	data.Work.parse()

	return data.Work, nil
}

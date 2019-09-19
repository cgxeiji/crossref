package crossref

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// QueryJSON gets the raw JSON from CrossRef. A search term must be specified,
// otherwise it will return an ErrEmptyQuery error. Use this if you want to
// manually manipulate the JSON data.
func (c *Client) QueryJSON(search string) ([]byte, error) {
	log.WithFields(log.Fields{
		"client": c,
	}).Debug("Query requested")

	if search == "" {
		return nil, ErrEmptyQuery
	}

	req, err := c.newRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("query", search)
	q.Add("rows", "10")
	q.Add("mailto", c.mailto)
	q.Add("sort", "relevance")
	req.URL.RawQuery = q.Encode()

	log.WithFields(log.Fields{
		"url": req.URL.String(),
	}).Debug("Requesting information")

	resp, err := c.Do(req)
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

// Query returns a list of processed Work metadata from a search query. A
// search term must be specified, otherwise it will return an ErrEmptyQuery
// error. If no works are found, it returns an ErrZeroWorks error.
func (c *Client) Query(search string) ([]*Work, error) {
	js, err := c.QueryJSON(search)
	works := []*Work{}
	if err != nil {
		return works, err
	}

	var data queryJSON
	json.Unmarshal(js, &data)

	if data.Message.TotalResults == 0 {
		return works, ErrZeroWorks
	}

	for _, work := range data.Message.Items {
		work.parse()
	}

	return data.Message.Items, nil
}

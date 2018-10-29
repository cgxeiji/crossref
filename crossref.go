package crossref

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"
)

// Client communicates with CrossRef API.
type Client struct {
	mailto  string
	appname string
}

const api string = "https://api.crossref.org"

// Debug checks the communication between the library and the API in detail.
func Debug() {
	log.SetLevel(log.DebugLevel)
	log.Info("Changed to debugging mode")
}

// String implements the Stringer interface.
func (c *Client) String() string {
	return fmt.Sprintf("App: %s, MailTo: %s", c.appname, c.mailto)
}

// WorksJSON returns the raw JSON. Use this if you want to manually process the
// JSON data.
func (c *Client) WorksJSON(doi, sel string) ([]byte, error) {
	log.WithFields(log.Fields{
		"client": c,
	}).Debug("JSON requested")

	url := fmt.Sprintf("%s/works?filter=doi:%s", api, doi)
	if sel != "" {
		url = fmt.Sprintf("%s&select=%s", url, sel)
	}
	url = fmt.Sprintf("%s&mailto=%s", url, c.mailto)

	log.WithFields(log.Fields{
		"url": url,
	}).Debug("Requesting information")

	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

// Contributor saves the name of the author in First Last name.
type Contributor struct {
	First string
	Last  string
}

// Work is the processed JSON for direct access to the information.
type Work struct {
	Type       string
	Titles     []string
	Title      string
	BookTitles []string
	BookTitle  string
	Authors    []Contributor
	Date       string
	Publisher  string
	Editors    []Contributor
	Issue      string
	Volume     string
	Pages      string
	DOI        string
	ISSNs      []string
	ISSN       string
	ISBNs      []string
	ISBN       string
}

func (w *Work) populate(content map[string]interface{}) {
	for _, v := range content["title"].([]interface{}) {
		w.Titles = append(w.Titles, getS(v))
	}
	w.Title = w.Titles[0]

	for _, v := range content["container-title"].([]interface{}) {
		w.BookTitles = append(w.BookTitles, getS(v))
	}
	if len(w.BookTitles) > 0 {
		w.BookTitle = w.BookTitles[0]
	}

	issns, _ := content["ISSN"].([]interface{})
	for _, v := range issns {
		w.ISSNs = append(w.ISSNs, getS(v))
	}
	if len(w.ISSNs) > 0 {
		w.ISSN = w.ISSNs[0]
	}

	isbns, _ := content["ISBN"].([]interface{})
	for _, v := range isbns {
		w.ISBNs = append(w.ISBNs, getS(v))
	}
	if len(w.ISBNs) > 0 {
		w.ISBN = w.ISBNs[0]
	}

	authors, _ := content["author"].([]interface{})
	for _, v := range authors {
		a := v.(map[string]interface{})
		co := Contributor{
			Last:  getS(a["family"]),
			First: getS(a["given"]),
		}

		w.Authors = append(w.Authors, co)
	}

	editors, _ := content["editor"].([]interface{})
	for _, v := range editors {
		e := v.(map[string]interface{})
		co := Contributor{
			Last:  getS(e["family"]),
			First: getS(e["given"]),
		}

		w.Editors = append(w.Editors, co)
	}

	w.Type = getS(content["type"])
	w.Publisher = getS(content["publisher"])
	w.Issue = getS(content["issue"])
	w.Volume = getS(content["volume"])
	w.Pages = getS(content["page"])
	w.DOI = getS(content["DOI"])

	date := content["issued"].(map[string]interface{})["date-parts"].([]interface{})[0].([]interface{})
	for _, v := range date {
		w.Date = fmt.Sprintf("%s%d-", w.Date, getI(v))
	}
	w.Date = strings.TrimRight(w.Date, "-")
}

// Works gets a processed metadata from a DOI.
func (c *Client) Works(doi string) (*Work, error) {
	js, err := c.WorksJSON(doi, "title,container-title,author,issued,DOI,type,issue,volume,page")
	if err != nil {
		return &Work{}, err
	}

	var data map[string]interface{}
	json.Unmarshal(js, &data)

	items := data["message"].(map[string]interface{})["items"].([]interface{})
	content := items[0].(map[string]interface{})
	log.Debug(pretty.Sprint(content))

	w := Work{}
	w.populate(content)

	return &w, nil
}

// QueryJSON gets the raw JSON from CrossRef. Use this if you want to manually
// manipulate the JSON data.
func (c *Client) QueryJSON(search string) ([]byte, error) {
	log.WithFields(log.Fields{
		"client": c,
	}).Debug("Query requested")

	rx, err := regexp.Compile("[^[:alnum:][:space:]]+")
	if err != nil {
		return nil, err
	}

	s := rx.ReplaceAllString(search, "")
	s = strings.Replace(s, " ", "+", -1)

	url := fmt.Sprintf("%s/works?query=%s&rows=10", api, s)
	url = fmt.Sprintf("%s&mailto=%s", url, c.mailto)

	log.WithFields(log.Fields{
		"url": url,
	}).Debug("Requesting information")

	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

// Query returns a list of processed Work metadata from a search query.
func (c *Client) Query(search string) ([]*Work, error) {
	js, err := c.QueryJSON(search)
	works := []*Work{}
	if err != nil {
		return works, err
	}

	var data map[string]interface{}
	json.Unmarshal(js, &data)

	items, _ := data["message"].(map[string]interface{})["items"].([]interface{})

	for _, v := range items {
		content := v.(map[string]interface{})
		log.Debug(pretty.Sprint(content))
		w := Work{}
		w.populate(content)

		works = append(works, &w)
	}

	return works, nil
}

func getS(i interface{}) string {
	s, _ := i.(string)
	return strings.TrimSpace(s)
}

func getI(i interface{}) int {
	n, _ := i.(float64)
	return int(n)
}

// NewClient returns a new crossref client with an attached mailto and app
// name.
func NewClient(appname, malito string) *Client {
	c := Client{
		mailto:  malito,
		appname: appname,
	}

	return &c
}

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

type Client struct {
	mailto  string
	appname string
}

const api string = "https://api.crossref.org"

func Debug() {
	log.SetLevel(log.DebugLevel)
	log.Info("Changed to debugging mode")
}

func (c *Client) String() string {
	return fmt.Sprintf("App: %s, MailTo: %s", c.appname, c.mailto)
}

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

type Contributor struct {
	First string
	Last  string
}

type Work struct {
	Type       string
	Titles     []string
	Title      string
	BookTitles []string
	BookTitle  string
	Authors    []Contributor
	Date       string
	DOI        string
	Issue      string
	Volume     string
	Pages      string
}

func (w *Work) populate(content map[string]interface{}) {
	for _, v := range content["title"].([]interface{}) {
		w.Titles = append(w.Titles, getS(v))
	}
	w.Title = w.Titles[0]

	for _, v := range content["container-title"].([]interface{}) {
		w.BookTitles = append(w.BookTitles, getS(v))
	}
	w.BookTitle = w.BookTitles[0]

	authors, _ := content["author"].([]interface{})
	for _, v := range authors {
		a := v.(map[string]interface{})
		co := Contributor{
			Last:  getS(a["family"]),
			First: getS(a["given"]),
		}

		w.Authors = append(w.Authors, co)
	}

	w.DOI = getS(content["DOI"])
	w.Type = getS(content["type"])
	w.Issue = getS(content["issue"])
	w.Volume = getS(content["volume"])
	w.Pages = getS(content["page"])

	date := content["issued"].(map[string]interface{})["date-parts"].([]interface{})[0].([]interface{})
	for _, v := range date {
		w.Date = fmt.Sprintf("%s%d-", w.Date, getI(v))
	}
	w.Date = strings.TrimRight(w.Date, "-")
}

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

func NewClient(appname, malito string) *Client {
	c := Client{
		mailto:  malito,
		appname: appname,
	}

	return &c
}

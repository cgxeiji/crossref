package crossref

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kr/pretty"
)

type Client struct {
	mailto  string
	appname string
}

const api string = "https://api.crossref.org"

func (c *Client) String() string {
	return fmt.Sprintf("App: %s, MailTo: %s", c.appname, c.mailto)
}

func (c *Client) WorksJSON(doi, sel string) ([]byte, error) {
	url := fmt.Sprintf("%s/works?filter=doi:%s", api, doi)
	if sel != "" {
		url = fmt.Sprintf("%s&select=%s", url, sel)
	}
	url = fmt.Sprintf("%s&mailto=%s", url, c.mailto)

	fmt.Println(url)
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
	BookTitles []string
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
	for _, v := range content["container-title"].([]interface{}) {
		w.BookTitles = append(w.BookTitles, getS(v))
	}
	for _, v := range content["author"].([]interface{}) {
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
	pretty.Println(content)

	w := Work{}
	w.populate(content)

	return &w, nil
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

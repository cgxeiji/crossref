package crossref

import (
	"fmt"
	"strings"
)

type queryJSON struct {
	Status  string            `json:"status"`
	Type    string            `json:"message-type"`
	Version string            `json:"message-version"`
	Message *queryMessageJSON `json:"message"`
}

type queryMessageJSON struct {
	Query        map[string]interface{} `json:"query"`
	Facets       map[string]interface{} `json:"facets"`
	TotalResults int                    `json:"total-results"`
	Items        []*Work                `json:"items"`
	ItemsPerPage int                    `json:"items-per-page"`
}

type doiJSON struct {
	Status  string `json:"status"`
	Type    string `json:"message-type"`
	Version string `json:"message-version"`
	Work    *Work  `json:"message"`
}

// Contributor saves the name of the author in First Last name.
type Contributor struct {
	// First is the given name of the contributor.
	First string `json:"given"`
	// Last is the family name of the contributor.
	Last string `json:"family"`
}

// String implements the Stringer interface.
func (c *Contributor) String() string {
	return c.Last + ", " + c.First
}

// Work is the processed JSON for direct access to the information.
type Work struct {
	Type string
	DOI  string `json:"DOI"`

	Titles []string `json:"title"`
	Title  string

	BookTitles []string `json:"container-title"`
	BookTitle  string

	Authors []*Contributor `json:"author"`
	Editors []*Contributor `json:"editor"`

	Issued *dateParts `json:"issued"`
	Date   string

	Publisher string `json:"publisher"`

	Issue    string   `json:"issue"`
	Volume   string   `json:"volume"`
	Pages    string   `json:"pages"`
	ISSNs    []string `json:"ISSN"`
	ISSN     string
	ISBNs    []string `json:"ISBN"`
	ISBN     string
	Abstract string `json:"abstract"`
}

func (w *Work) parse() {
	if w.Titles != nil {
		if len(w.Titles) > 0 {
			w.Title = w.Titles[0]
		}
	}
	if w.BookTitles != nil {
		if len(w.BookTitles) > 0 {
			w.BookTitle = w.BookTitles[0]
		}
	}
	if w.ISBNs != nil {
		if len(w.ISBNs) > 0 {
			w.ISBN = w.ISBNs[0]
		}
	}
	if w.ISSNs != nil {
		if len(w.ISSNs) > 0 {
			w.ISSN = w.ISSNs[0]
		}
	}
	for _, v := range w.Issued.Parts[0] {
		w.Date = fmt.Sprintf("%s%d-", w.Date, v)
	}
	w.Date = strings.TrimRight(w.Date, "-")
}

type dateParts struct {
	Parts [][]int `json:"date-parts"`
}

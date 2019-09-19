# crossref
Access to Crossref API with Go

Full documentation of this packages is [here](https://godoc.org/github.com/cgxeiji/crossref).

For a detailed explanation of the JSON fields, go to [Crossref API
Documentation](https://github.com/Crossref/rest-api-doc/blob/master/api_format.md).

## Example Code
``` Go
package main

import (
	"fmt"

	"github.com/cgxeiji/crossref"
)

var client = crossref.NewClient("Crossref Go", "mail@example.com")

func main() {
	// If you want to debug the library, uncomment the following line:
	// crossref.Debug()

	// Search for a publication by doing a query. This returns up to 10 works
	// that match the query terms.
	search := "Slow Robots for Unobtrusive Posture Correction"
	works, err := client.Query(search)
	switch err {
	case crossref.ErrZeroWorks:
		fmt.Println("No works found")
	case crossref.ErrEmptyQuery:
		fmt.Println("An empty query was requested")
	case nil:
	default:
		panic(err)
	}

	fmt.Printf("Found %d article(s) for query:\n > %q\n", len(works), search)

	// Retrieve information directly from a DOI
	doi := works[0].DOI
	work, err := client.DOI(doi)
	switch err {
	case crossref.ErrZeroWorks:
		fmt.Println("No works found")
	case crossref.ErrEmptyQuery:
		fmt.Println("An empty query was requested")
	case nil:
	default:
		panic(err)
	}

	fmt.Printf("For DOI: %q found:\n", doi)
	fmt.Printf(" > %q, (%v) %q\n", work.Title, work.Date, work.Authors)

	// Output:
	// Found 10 article(s) for query:
	//  > "Slow Robots for Unobtrusive Posture Correction"
	// For DOI: "10.1145/3290605.3300843" found:
	//  > "Slow Robots for Unobtrusive Posture Correction", (2019) ["Shin, Joon-Gi" "Onchi, Eiji" "Reyes, Maria Jose" "Song, Junbong" "Lee, Uichin" "Lee, Seung-Hee" "Saakes, Daniel"]
}
```

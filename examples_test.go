package crossref_test

import (
	"fmt"

	"github.com/cgxeiji/crossref"
)

func ExampleClient_Query() {
	// Search for a work on Crossref.
	works, err := client.Query("Slow Robots for Unobtrusive Posture Correction")
	if err != nil {
		panic(err)
	}

	work := works[0]
	fmt.Println(work.DOI)

	// If no works can be found, Query returns with ErrZeroWorks error.
	_, err = client.Query("jtfiejfrlsadaksljablkjoifajebwoijffal")
	switch err {
	case crossref.ErrZeroWorks:
		fmt.Println("No works found")
	case crossref.ErrEmptyQuery:
		fmt.Println("An empty query was requested")
	case nil:
	default:
		panic(err)
	}
	// Output:
	// 10.1145/3290605.3300843
	// No works found
}

func ExampleClient_DOI() {
	// Search for a DOI metadata on Crossref.
	work, err := client.DOI("10.1145/3290605.3300843")
	if err != nil {
		panic(err)
	}

	fmt.Println(work.Title)

	// If no works can be found, DOI returns with ErrZeroWorks error.
	_, err = client.DOI("10.1145/3290605.330084") // deleted the last digit
	switch err {
	case crossref.ErrZeroWorks:
		fmt.Println("No work found")
	case crossref.ErrEmptyQuery:
		fmt.Println("An empty DOI was requested")
	case nil:
	default:
		panic(err)
	}
	// Output:
	// Slow Robots for Unobtrusive Posture Correction
	// No work found
}

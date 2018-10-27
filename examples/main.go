package main

import (
	"fmt"

	"github.com/cgxeiji/crossref"
	"github.com/kr/pretty"
)

func main() {
	crossref.Debug()
	client := crossref.NewClient("CrossRefGo", "mail@example.com")
	fmt.Println(client)
	// work, _ := client.Works("http://dx.doi.org/10.1016/0004-3702(89)90008-8")
	// pretty.Println(work)
	// work, _ = client.Works("http://dx.doi.org/10.1117/12.969296")
	// pretty.Println(work)

	q, _ := client.Query("Sarpi, a solution for artificial intelligence ")
	pretty.Println(q)

	fmt.Println()
	fmt.Println("Found", len(q), "article(s):")
	for _, v := range q {
		fmt.Println(" >", v.Titles[0])
		fmt.Println("    ", v.Authors, "#", v.Date)
	}
}

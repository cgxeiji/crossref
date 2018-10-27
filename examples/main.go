package main

import (
	"fmt"

	"github.com/cgxeiji/crossref"
	"github.com/kr/pretty"
)

func main() {
	client := crossref.NewClient("CrossRefGo", "mail@example.com")
	fmt.Println(client)
	work, _ := client.Works("http://dx.doi.org/10.1016/0004-3702(89)90008-8")
	pretty.Println(work)
	// work, _ = client.Works("http://dx.doi.org/10.1117/12.969296")
	// pretty.Println(work)
}

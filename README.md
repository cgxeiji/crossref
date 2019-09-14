# crossref
Access to Crossref API with Go

For a detailed explanation of the JSON fields, go to [Crossref API
Documentation](https://github.com/Crossref/rest-api-doc/blob/master/api_format.md).

## Example Code
``` Go
func main() {
    // If you want to debug the library, uncomment the following line:
    // crossref.Debug()
    client := crossref.NewClient("Crossref Go", "mail@example.com")

    // Retrieve information directly from a DOI
    work, err := client.DOI("10.1145/3290605.3300843")
    if err != nil {
        panic(err)
    }

    fmt.Println("Found:")
    fmt.Println(" >", work.Title, work.Authors, work.Date)

    // Search for a publication by doing a query. This returns up to 10 works
    // that match the query terms.
    works, err := client.Query("Slow Robots for Unobtrusive Posture Correction")
    if err != nil {
        panic(err)
    }

    fmt.Println("Found", len(works), "article(s):")
    for _, work := range works {
        fmt.Println(" >", work.Title, work.Authors, work.Date)
    }
}
```

package crossref

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	want := Client{
		mailto:  "mail@example.com",
		appname: "Crossref Go",
		Client:  http.Client{Timeout: 5 * time.Second},
	}

	got := NewClient("Crossref Go", "mail@example.com")

	if got.mailto != want.mailto {
		t.Errorf("\ngot: \n\t%v\nwant: \n\t%v", got.mailto, want.mailto)
	}
	if got.appname != want.appname {
		t.Errorf("\ngot: \n\t%v\nwant: \n\t%v", got.appname, want.appname)
	}
	if got.Client.Timeout != want.Client.Timeout {
		t.Errorf("\ngot: \n\t%v\nwant: \n\t%v", got.Client.Timeout, want.Client.Timeout)
	}
}

func TestString(t *testing.T) {
	client := NewClient("Crossref Go", "mail@example.com")

	want := "App: Crossref Go, MailTo: mail@example.com"
	got := client.String()

	if got != want {
		t.Errorf("\ngot: \n\t%v\nwant: \n\t%v", got, want)
	}
}

package crossref

import (
	"encoding/json"
	"testing"
)

func TestClient_DOIJSON(t *testing.T) {
	client := NewClient("Crossref Go", "mail@example.com")

	t.Run("Existing Entry", func(t *testing.T) {
		want := "Slow Robots for Unobtrusive Posture Correction"

		raw, err := client.DOIJSON("10.1145/3290605.3300843")
		if err != nil {
			t.Fatal(err)
		}

		var data doiJSON
		json.Unmarshal(raw, &data)

		if data.Work == nil {
			t.Fatal("could not unmarshal json to work")
		}
		got := data.Work.Titles[0]

		if want != got {
			t.Errorf("\ngot: \n\t%v\nwant: \n\t%v", got, want)
		}
	})

	t.Run("Non-existing Entry", func(t *testing.T) {
		raw, err := client.DOIJSON("10.1145/3290605.330084") // deleted the last character from DOI
		if err != nil {
			t.Fatal(err)
		}

		var data doiJSON
		json.Unmarshal(raw, &data)

		if data.Status == "ok" {
			t.Error("unexpected entry found for non-existing entry")
		}
	})

	t.Run("Empty Query", func(t *testing.T) {
		_, err := client.DOIJSON("")
		if err != ErrEmptyQuery {
			t.Error("failed to return an empty query error:", err)
		}
	})
}

func TestClient_DOI(t *testing.T) {
	client := NewClient("Crossref Go", "mail@example.com")

	t.Run("Existing Entry", func(t *testing.T) {
		want := "Slow Robots for Unobtrusive Posture Correction"

		work, err := client.DOI("10.1145/3290605.3300843")
		if err != nil {
			t.Fatal(err)
		}

		got := work.Title

		if want != got {
			t.Errorf("\ngot: \n\t%v\nwant: \n\t%v", got, want)
		}
	})

	t.Run("Non-existing Entry", func(t *testing.T) {
		_, err := client.DOI("10.1145/3290605.330084") // deleted the last character from DOI
		if err != ErrZeroWorks {
			t.Error("returned an error other than ErrZeroWorks for a non-existing entry:", err)
		}
	})

	t.Run("Empty Query", func(t *testing.T) {
		_, err := client.DOI("")
		if err != ErrEmptyQuery {
			t.Error("failed to return an empty query error:", err)
		}
	})
}

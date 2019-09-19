package crossref

import (
	"encoding/json"
	"sync"
	"testing"
)

func TestClient_QueryJSON(t *testing.T) {
	client := NewClient("Crossref Go", "mail@example.com")

	// Test a known existing entry
	t.Run("Existing Entry", func(t *testing.T) {
		want := "Slow Robots for Unobtrusive Posture Correction"

		raw, err := client.QueryJSON(want)
		if err != nil {
			t.Error(err)
		}

		var data queryJSON
		json.Unmarshal(raw, &data)

		t.Log("\n\ttotal results =", data.Message.TotalResults)
		if data.Message.TotalResults == 0 {
			t.Fatal("known existing entry could not be found")
		}

		got := data.Message.Items[0].Titles[0]

		if want != got {
			t.Errorf("\ngot: \n\t%v\nwant: \n\t%v", got, want)
		}
	})

	// Test a known non-existing entry
	t.Run("Non-existing Entry", func(t *testing.T) {
		search := "jtfiejfrlsadaksljablkjoifajebwoijffal"

		raw, err := client.QueryJSON(search)
		if err != nil {
			t.Error(err)
		}

		var data queryJSON
		json.Unmarshal(raw, &data)

		got := data.Message.TotalResults
		want := 0

		if want != got {
			t.Errorf("\ngot: \n\t%v\nwant: \n\t%v", got, want)
		}
	})

	t.Run("Empty Query", func(t *testing.T) {
		_, err := client.QueryJSON("")
		if err != ErrEmptyQuery {
			t.Error("failed to return an empty query error:", err)
		}
	})

	t.Run("Concurrent", func(t *testing.T) {
		queries := []string{
			"Slow Robots for Unobtrusive Posture Correction",
			"jtfiejfrlsadaksljablkjoifajebwoijffal",
			"",
		}
		var wg sync.WaitGroup
		for _, q := range queries {
			wg.Add(1)
			go func(q string) {
				defer wg.Done()
				raw, err := client.QueryJSON(q)
				t.Logf("LOG: query: %q, got: %d bytes, error: %v", q, len(raw), err)
			}(q)
		}
		wg.Wait()
	})
}

func TestClient_Query(t *testing.T) {
	client := NewClient("Crossref Go", "mail@example.com")

	t.Run("Existing Entry", func(t *testing.T) {
		want := "Slow Robots for Unobtrusive Posture Correction"

		works, err := client.Query(want)
		if err != nil {
			t.Fatal(err)
		}

		got := works[0].Title

		if want != got {
			t.Errorf("\ngot: \n\t%v\nwant: \n\t%v", got, want)
		}
	})

	t.Run("Non-existing Entry", func(t *testing.T) {
		_, err := client.Query("jtfiejfrlsadaksljablkjoifajebwoijffal")
		if err != ErrZeroWorks {
			t.Error("returned an error other than ErrZeroWorks for a non-existing entry:", err)
		}
	})

	t.Run("Empty Query", func(t *testing.T) {
		_, err := client.Query("")
		if err != ErrEmptyQuery {
			t.Error("failed to return an empty query error:", err)
		}
	})

	t.Run("Concurrent", func(t *testing.T) {
		queries := []string{
			"Slow Robots for Unobtrusive Posture Correction",
			"jtfiejfrlsadaksljablkjoifajebwoijffal",
			"",
		}
		var wg sync.WaitGroup
		for _, q := range queries {
			wg.Add(1)
			go func(q string) {
				defer wg.Done()
				works, err := client.Query(q)
				t.Logf("LOG: query: %q, got: %d works, error: %v", q, len(works), err)
			}(q)
		}
		wg.Wait()
	})
}

package eudic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIDString(t *testing.T) {
	cases := []struct {
		in   any
		want string
	}{
		{nil, "0"},
		{"", "0"},
		{"132", "132"},
		{float64(42), "42"},
		{0, "0"},
	}
	for _, c := range cases {
		if got := IDString(c.in); got != c.want {
			t.Fatalf("IDString(%v)=%q want %q", c.in, got, c.want)
		}
	}
}

func TestFindCategoryID(t *testing.T) {
	response := map[string]any{
		"data": []any{
			map[string]any{"id": float64(17), "name": "general"},
			map[string]any{"id": "42", "name": "english-coach"},
		},
	}
	id, ok := findCategoryID(response, "english-coach")
	if !ok || id != "42" {
		t.Fatalf("findCategoryID()=(%q, %v), want (%q, true)", id, ok, "42")
	}
}

func TestSyncEntriesUsesExistingCategoryAndWritesEachEntry(t *testing.T) {
	var requests []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := map[string]any{"method": r.Method, "path": r.URL.Path}
		if r.Body != nil {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
				request["body"] = body
			}
		}
		requests = append(requests, request)

		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet && r.URL.Path == "/studylist/category" {
			_, _ = w.Write([]byte(`{"data":[{"id":"134","name":"english-coach"}]}`))
			return
		}
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	client := &Client{Token: "test", HTTPClient: server.Client(), BaseURL: server.URL}
	entries := []SyncEntry{
		{Word: "documentation", ContextLine: "Split it into two files.", Note: "note one", Star: 3},
		{Word: "split into", ContextLine: "Split it into two files.", Note: "note two", Star: 3},
	}
	synced, err := client.SyncEntries(context.Background(), "english-coach", "en", entries)
	if err != nil {
		t.Fatalf("SyncEntries() error: %v", err)
	}
	if synced != 2 {
		t.Fatalf("SyncEntries()=%d, want 2", synced)
	}
	if len(requests) != 5 {
		t.Fatalf("got %d requests, want 5: %#v", len(requests), requests)
	}
	if requests[0]["method"] != http.MethodGet || requests[0]["path"] != "/studylist/category" {
		t.Fatalf("first request = %#v", requests[0])
	}
	for _, index := range []int{1, 3} {
		body := requests[index]["body"].(map[string]any)
		ids := body["category_ids"].([]any)
		if len(ids) != 1 || ids[0] != "134" {
			t.Fatalf("request %d category_ids = %#v, want [134]", index, ids)
		}
	}
}

func TestSyncEntriesCreatesMissingCategory(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/studylist/category":
			_, _ = w.Write([]byte(`{"data":[]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/studylist/category":
			_, _ = w.Write([]byte(`{"data":{"categoryId":99}}`))
		default:
			_, _ = w.Write([]byte(`{"success":true}`))
		}
	}))
	defer server.Close()

	client := &Client{Token: "test", HTTPClient: server.Client(), BaseURL: server.URL}
	synced, err := client.SyncEntries(context.Background(), "english-coach", "en", []SyncEntry{{
		Word: "refined", Note: "note", Star: 3,
	}})
	if err != nil {
		t.Fatalf("SyncEntries() error: %v", err)
	}
	if synced != 1 {
		t.Fatalf("SyncEntries()=%d, want 1", synced)
	}
	if got, want := len(paths), 4; got != want {
		t.Fatalf("request count=%d, want %d: %v", got, want, paths)
	}
}

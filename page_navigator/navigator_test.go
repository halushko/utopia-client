package page_navigator

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETRejectsLoginHTML(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<form action="/login"><input name="password" type="password"></form>`))
	}))
	defer s.Close()
	old := c
	c = &Client{navigateClient: s.Client()}
	defer func() { c = old }()
	if _, err := GET(s.URL); err == nil {
		t.Error("accepted login as torrent response")
	}
}

func TestGETStatusError(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(429) }))
	defer s.Close()
	old := c
	c = &Client{navigateClient: s.Client()}
	defer func() { c = old }()
	if _, err := GET(s.URL); err == nil {
		t.Error("accepted HTTP failure")
	}
}

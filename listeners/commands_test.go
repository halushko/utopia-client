package listeners

import (
	"strings"
	"testing"
	"utopia-client/entities"
)

func TestInvalidIDsDoNotFetchOrExit(t *testing.T) {
	for _, id := range []string{"bad", "0", "-1", "999999999999999999999999"} {
		if _, ok := getTorrentInfo(id); ok {
			t.Errorf("accepted %q", id)
		}
		if getTorrentFile(id) {
			t.Errorf("download accepted %q", id)
		}
	}
}

func TestFormattingMissingMediaAndHTML(t *testing.T) {
	got := prepareTorrentInfoMessage(&entities.Torrent{ID: 1, Title: "A < B & C", Category: "TV"})
	if !strings.Contains(got, "A &lt; B &amp; C") {
		t.Errorf("unescaped title: %s", got)
	}
	if strings.Contains(got, "/utp_siml_") {
		t.Errorf("empty similar command: %s", got)
	}
}

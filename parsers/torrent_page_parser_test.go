package parsers

import (
	"os"
	"strings"
	"testing"
)

func fixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestObservedEpisode(t *testing.T) {
	got, err := ParseTorrentPage(101, fixture(t, "episode.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got.Title, "S02E09") {
		t.Errorf("release title = %q", got.Title)
	}
	if got.FileCount != 1 || got.Files[0].Name != "Example.Show.S02E09.mkv" {
		t.Errorf("files = %#v", got.Files)
	}
	if got.PrimaryLanguage != "en" {
		t.Errorf("primary language = %q", got.PrimaryLanguage)
	}
	if got.File.Media.Resolution != "2160p" {
		t.Errorf("quality = %q", got.File.Media.Resolution)
	}
	if len(got.File.Audio) != 2 || got.File.Audio[0].Studio != "KyivstarTV" || got.File.Audio[1].Studio != "" {
		t.Errorf("audio = %#v", got.File.Audio)
	}
}

func TestObservedPack(t *testing.T) {
	got, err := ParseTorrentPage(102, fixture(t, "pack.html"))
	if err != nil {
		t.Fatal(err)
	}
	if got.FileCount != 3 {
		t.Fatalf("pack files = %#v", got.Files)
	}
	if got.File.Name != "Example.S01E01.mkv" {
		t.Errorf("representative media filename = %q", got.File.Name)
	}
	if got.Files[1].Media != nil {
		t.Error("must not copy representative MediaInfo to all pack files")
	}
	if len(got.File.Audio) != 3 || got.File.Audio[0].Studio != "1+1 Production" || got.File.Audio[1].Studio != "HDREZKA" {
		t.Errorf("audio = %#v", got.File.Audio)
	}
}

func TestRejectNonTorrentPage(t *testing.T) {
	for _, page := range []string{"", `<form action="/login"><input name="password"></form>`, `<h1>Just a moment...</h1>`} {
		if _, err := ParseTorrentPage(1, page); err == nil {
			t.Errorf("accepted non-torrent HTML %q", page)
		}
	}
}

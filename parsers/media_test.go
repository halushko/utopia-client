package parsers

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestVideoDetailsArePreserved(t *testing.T) {
	got, err := ParseTorrentPage(101, fixture(t, "episode.html"))
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(got.File.Media)
	for _, part := range []string{`"video_codec":"HEVC"`, `"width":"3 832 pixels"`, `"height":"1 600 pixels"`, `"video_bit_rate":"24.7 Mb/s"`, `Dolby Vision`} {
		if !strings.Contains(string(b), part) {
			t.Errorf("missing %s in %s", part, b)
		}
	}
}

func TestExactMediaKeys(t *testing.T) {
	page := strings.Replace(fixture(t, "episode.html"), "Overall bit rate : 25.7 Mb/s", "Overall bit rate mode : Variable\nOverall bit rate : 25.7 Mb/s", 1)
	got, _ := ParseTorrentPage(101, page)
	if got.File.Media.OverallBitRate != "25.7 Mb/s" || got.File.Media.OverallBitRateMode != "Variable" {
		t.Errorf("media = %#v", got.File.Media)
	}
}

package page_navigator

import "testing"

func TestTorrentResponseRejectsErrorPages(t *testing.T) {
	for _, body := range []string{"", "<html>login</html>", `{"message":"Unauthorized"}`, "forbidden"} {
		if err := checkTorrentResponse([]byte(body)); err == nil {
			t.Errorf("accepted %q", body)
		}
	}
	if err := checkTorrentResponse([]byte("d4:infod4:name1:xee")); err != nil {
		t.Fatal(err)
	}
}

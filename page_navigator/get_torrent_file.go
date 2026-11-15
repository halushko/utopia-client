package page_navigator

import (
	"fmt"
)

func GetTorrentFile(id int) ([]byte, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid torrent ID")
	}
	pageURL := fmt.Sprintf("https://utp.to/torrents/download/%d", id)
	body, err := GET(pageURL)
	if err != nil {
		return nil, err
	}
	if err := checkTorrentResponse(body); err != nil {
		return nil, err
	}
	return body, nil
}

// Reject HTML/JSON error responses before writing to the watched directory.
// This is an envelope check, not a complete torrent metadata validator.
func checkTorrentResponse(body []byte) error {
	if len(body) < 2 || body[0] != 'd' || body[len(body)-1] != 'e' {
		return fmt.Errorf("UTOPIA returned a non-bencode torrent response")
	}
	return nil
}

package page_navigator

import (
	"fmt"
)

func GetTorrentFile(id int) ([]byte, error) {
	pageURL := fmt.Sprintf("https://utp.to/torrents/download/%d", id)
	body, err := GET(pageURL)
	if err != nil {
		return nil, err
	}
	return body, nil
}

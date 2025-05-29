// get_torrent_info.go
package page_navigator

import (
	"fmt"
)

func GetTorrentPage(id int) (string, error) {
	pageURL := fmt.Sprintf("https://utp.to/torrents/%d", id)
	body, err := GET(pageURL)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

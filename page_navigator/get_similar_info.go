// get_torrent_info.go
package page_navigator

import (
	"fmt"
	"regexp"
)

func GetSimilarPage(id int) (string, error) {
	return GetSimilarPageByID(fmt.Sprint(id))
}

// Keep the old int entry point, while accepting real category.TMDB identifiers.
func GetSimilarPageByID(id string) (string, error) {
	if !regexp.MustCompile(`^[1-9][0-9]*(\.[1-9][0-9]*)?$`).MatchString(id) {
		return "", fmt.Errorf("invalid similar ID")
	}
	pageURL := fmt.Sprintf("https://utp.to/torrents/similar/%s", id)
	body, err := GET(pageURL)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

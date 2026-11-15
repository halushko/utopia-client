package parsers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"utopia-client/entities"
)

func ParseSimilar(id int, html string) (*[]entities.SimilarTorrent, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	if doc.Find("h1.meta__title").Length() == 0 {
		return nil, fmt.Errorf("similar page metadata missing")
	}
	result := []entities.SimilarTorrent{}
	seen := map[int]bool{}
	doc.Find("table").Each(func(_ int, table *goquery.Selection) {
		kind := ""
		table.Find("tbody > tr").Each(func(_ int, row *goquery.Selection) {
			if v := strings.TrimSpace(row.Find(".similar-torrents__type").Text()); v != "" {
				kind = v
			}
			a := row.Find(".torrent-search--grouped__name a").First()
			href, ok := a.Attr("href")
			if !ok {
				return
			}
			u, err := url.Parse(href)
			if err != nil {
				return
			}
			if u.IsAbs() && (u.Scheme != "https" || u.Host != "utp.to") {
				return
			}
			match := regexp.MustCompile(`^/torrents/([1-9][0-9]*)$`).FindStringSubmatch(u.Path)
			if match == nil {
				return
			}
			torrentID, err := strconv.Atoi(match[1])
			if err != nil || seen[torrentID] {
				return
			}
			seen[torrentID] = true
			item := entities.SimilarTorrent{ID: torrentID, Title: strings.TrimSpace(a.Text()), URL: "https://utp.to" + u.Path, Type: kind}
			row.ParentsFiltered("details").Each(func(_ int, group *goquery.Selection) {
				label := strings.TrimSpace(group.ChildrenFiltered("summary").First().Text())
				parts := strings.Fields(label)
				if len(parts) != 2 {
					return
				}
				n, err := strconv.Atoi(parts[1])
				if err != nil {
					return
				}
				if parts[0] == "Season" {
					item.Season = &n
				}
				if parts[0] == "Episode" {
					item.Episode = &n
				}
			})
			size, _ := row.Find(".torrent-search--grouped__size span").Attr("title")
			if fields := strings.Fields(size); len(fields) > 0 {
				item.SizeBytes, _ = strconv.ParseInt(fields[0], 10, 64)
			}
			item.Seeders, _ = strconv.Atoi(strings.TrimSpace(row.Find(".torrent-search--grouped__seeders").Text()))
			item.Leechers, _ = strconv.Atoi(strings.TrimSpace(row.Find(".torrent-search--grouped__leechers").Text()))
			item.Completed, _ = strconv.Atoi(strings.TrimSpace(row.Find(".torrent-search--grouped__completed").Text()))
			date, _ := row.Find("time").Attr("datetime")
			item.UploadedAt, _ = time.Parse("2006-01-02 15:04:05", date)
			item.PersonalRelease = row.Find(".torrent-icons__personal-release").Length() > 0
			item.FreeLeech = row.Find(".torrent-icons__freeleech").Length() > 0
			row.Find(".torrent-icons [title]").Each(func(_ int, icon *goquery.Selection) {
				v, _ := icon.Attr("title")
				item.Attributes = append(item.Attributes, v)
			})
			result = append(result, item)
		})
	})
	return &result, nil
}

package parsers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"utopia-client/entities"

	"github.com/PuerkitoBio/goquery"
)

func ParseTorrentPage(id int, html string) (*entities.Torrent, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	t := &entities.Torrent{ID: id}
	t.Title = getTitle(doc)
	if t.Title == "" {
		return nil, fmt.Errorf("torrent release title missing: unexpected or unauthenticated page")
	}
	t.Description = getDescription(doc)
	//t.Size = getSize(doc)
	t.Category = getCategory(doc)
	t.Genres = getGenres(doc)
	t.PrimaryLanguage = getPrimaryLanguage(doc)
	fileName, files, filesSizes := getFiles(doc)
	t.File = entities.File{
		Name: fileName,
		Size: getSize(doc),
	}

	filesArray := make([]entities.File, 0)

	for i := 0; i < len(files); i++ {
		filesArray = append(
			filesArray,
			entities.File{
				Name:  files[i],
				Size:  filesSizes[i],
				Media: nil,
				Audio: nil,
			},
		)
	}

	t.Files = filesArray
	t.FileCount = len(files)
	t.SimilarId = getSimilar(doc)

	t.File.Media, t.File.Audio = getMediaAudioInfo(doc)
	if quality := strings.TrimSpace(doc.Find(".torrent__resolution-link").First().Text()); quality != "" {
		t.File.Media.Resolution = quality
	}

	return t, nil
}

func getTitle(doc *goquery.Document) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(doc.Find("h1.torrent__name").First().Text()), " ")
}

func getDescription(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("p.meta__description").Text())
}

func getSize(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("span.torrent__size-link").Text())
}

func getCategory(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("a.torrent__category-link").Text())
}

func getGenres(doc *goquery.Document) []string {
	genres := strings.TrimSpace(doc.Find("article.meta__genres").Find("h3.meta-chip__value").Text())

	parts := strings.Split(genres, "/")
	res := make([]string, 0)
	for _, p := range parts {
		g := strings.TrimSpace(p)
		if g != "" {
			res = append(res, g)
		}
	}
	return res
}

func getPrimaryLanguage(doc *goquery.Document) string {
	if value := strings.TrimSpace(doc.Find(".work__language-link").First().Text()); value != "" {
		return value
	}
	return strings.TrimSpace(doc.Find("article.meta__language h3.meta-chip__value").First().Text())
}

func getFiles(doc *goquery.Document) (string, []string, []string) {
	files, sizes := []string{}, []string{}
	// The flat list preserves full paths and excludes folder summary rows.
	doc.Find(`dialog [data-tab="list"] tbody tr`).Each(func(_ int, row *goquery.Selection) {
		cells := row.ChildrenFiltered("td")
		if cells.Length() != 3 {
			return
		}
		name := strings.TrimSpace(cells.Eq(1).Text())
		if name == "" {
			return
		}
		files = append(files, name)
		sizes = append(sizes, strings.TrimSpace(cells.Eq(2).Text()))
	})
	name := ""
	for _, line := range strings.Split(doc.Find(`code[x-ref="mediainfo"]`).First().Text(), "\n") {
		if ok, value := getValue(strings.TrimSpace(line), "Complete name"); ok {
			name = value
			break
		}
	}
	if name == "" && len(files) == 1 {
		name = files[0]
	}
	return name, files, sizes
}

func getSimilar(doc *goquery.Document) string {
	sel := doc.Find("a.meta__title-link")
	id := ""
	href, ok := sel.Attr("href")
	if ok {
		parts := strings.Split(strings.TrimRight(href, "/"), "/")
		id = parts[len(parts)-1]
	}
	return id
}

func getMediaAudioInfo(doc *goquery.Document) (*entities.Media, []entities.Audio) {
	audio := make([]entities.Audio, 0)
	media := entities.Media{}

	row := doc.Find("code[x-ref=\"mediainfo\"]").First().Text()
	GENERAL := 1
	VIDEO := 2
	AUDIO := 3
	section := GENERAL
	for _, line := range strings.Split(row, "\n") {
		tmp := strings.TrimSpace(line)
		if tmp != "" {
			if strings.HasPrefix(tmp, "Text") || tmp == "Menu" || tmp == "Image" || tmp == "Other" {
				break
			} else if strings.HasPrefix(tmp, "Video") {
				section = VIDEO
			} else if strings.HasPrefix(tmp, "Audio") {
				section = AUDIO
			}

			switch section {
			case GENERAL:
				fillGeneralValues(&media, tmp)
			case VIDEO:
				fillVideoValues(&media, tmp)
			default:
				if isNewAudio, _ := getValue(tmp, "ID"); isNewAudio {
					audio = append(audio, entities.Audio{})
				}
				if len(audio) > 0 {
					fillAudioValues(&audio[len(audio)-1], tmp)
				}
			}
		}
	}
	return &media, audio
}

func fillGeneralValues(media *entities.Media, line string) {
	isPresent := false
	value := ""
	if isPresent, value = getValue(line, "Format version"); isPresent {
		media.FormatVersion = value
	} else if isPresent, value = getValue(line, "Format"); isPresent {
		media.Format = value
	} else if isPresent, value = getValue(line, "Overall bit rate mode"); isPresent {
		media.OverallBitRateMode = value
	} else if isPresent, value = getValue(line, "Overall bit rate"); isPresent {
		media.OverallBitRate = value
	} else if isPresent, value = getValue(line, "Frame rate"); isPresent {
		media.FrameRate = value
	}
}

func fillVideoValues(media *entities.Media, line string) {
	if ok, value := getValue(line, "Format"); ok {
		media.VideoCodec = value
	}
	if ok, value := getValue(line, "Width"); ok {
		media.Width = value
	}
	if ok, value := getValue(line, "Height"); ok {
		media.Height = value
	}
	if ok, value := getValue(line, "Bit rate"); ok {
		media.VideoBitRate = value
	}
	if ok, value := getValue(line, "HDR format"); ok {
		media.HDRFormat = value
	}
	isPresent := false
	value := ""
	if isPresent, value = getValue(line, "Height"); isPresent {
		parts := strings.Split(value, " ")
		resolution := ""
		for _, part := range parts {
			if _, err := strconv.Atoi(part); err == nil {
				resolution = resolution + part
			}
		}
		media.Resolution = resolution + "p"
	}
}

func fillAudioValues(audio *entities.Audio, line string) {
	isPresent := false
	value := ""

	if isPresent, value = getValue(line, "Format/Info"); isPresent {
		audio.FormatInfo = value
	} else if isPresent, value = getValue(line, "Format"); isPresent {
		audio.Format = value
	} else if isPresent, value = getValue(line, "Commercial name"); isPresent {
		audio.CommercialName = value
	} else if isPresent, value = getValue(line, "Codec ID"); isPresent {
		audio.CodecID = value
	} else if isPresent, value = getValue(line, "Bit rate mode"); isPresent {
		audio.BitRateMode = value
	} else if isPresent, value = getValue(line, "Bit rate"); isPresent {
		audio.BitRate = value
	} else if isPresent, value = getValue(line, "Channel(s)"); isPresent {
		audio.Channels = value
	} else if isPresent, value = getValue(line, "Channel layout"); isPresent {
		audio.ChannelLayout = value
	} else if isPresent, value = getValue(line, "Sampling rate"); isPresent {
		audio.SamplingRate = value
	} else if isPresent, value = getValue(line, "Frame rate"); isPresent {
		audio.FrameRate = value
	} else if isPresent, value = getValue(line, "Compression mode"); isPresent {
		audio.CompressionMode = value
	} else if isPresent, value = getValue(line, "Stream size"); isPresent {
		audio.StreamSize = value
	} else if isPresent, value = getValue(line, "Language"); isPresent {
		audio.Language = value
	} else if isPresent, value = getValue(line, "Default"); isPresent {
		if value == "Yes" {
			audio.Default = true
		} else {
			audio.Default = false
		}
	} else if isPresent, value = getValue(line, "Forced"); isPresent {
		if value == "Yes" {
			audio.Forced = true
		} else {
			audio.Forced = false
		}
	} else if isPresent, value = getValue(line, "Title"); isPresent {
		audio.Title = value
		audio.Studio = studioFromTitle(value)
	}
}

func getValue(line string, key string) (bool, string) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
		return true, strings.TrimSpace(parts[1])
	}
	return false, ""
}

// A studio is accepted only in an explicit voice-over title. Generic audio
// labels and bitrates are not studio names; unknown remains empty.
func studioFromTitle(title string) string {
	parts := strings.Split(title, "|")
	if len(parts) < 2 {
		return ""
	}
	hasVoiceType := false
	for _, p := range parts[:len(parts)-1] {
		switch strings.ToUpper(strings.TrimSpace(p)) {
		case "MVO", "DVO", "VO", "DUB", "DUBBED", "AVO":
			hasVoiceType = true
		}
	}
	studio := strings.TrimSpace(parts[len(parts)-1])
	if !hasVoiceType || studio == "" || regexp.MustCompile(`(?i)(bps|kb/s|mb/s|channels?|^\d+(?:\.\d+)?$)`).MatchString(studio) {
		return ""
	}
	return studio
}

package parsers

import (
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

	return t, nil
}

func getTitle(doc *goquery.Document) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(doc.Find("h1.meta__title").Text()), " ")
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
	return strings.TrimSpace(doc.Find("article.meta__language").Find("h3.meta-chip__value").Text())
}

func getFiles(doc *goquery.Document) (string, []string, []string) {
	files := make([]string, 0)
	sizes := make([]string, 0)
	mainFileName := ""
	i := -1
	isName := true
	fromPageSrc := doc.Find("li.form__group").Find("div.dialog__form").Text()

	for _, line := range strings.Split(fromPageSrc, "\n") {
		tmp := strings.TrimSpace(line)
		if tmp != "" {
			i++
			if i > 2 {
				switch {
				case isName:
					files = append(files, tmp)
					isName = false
				default:
					sizes = append(sizes, tmp)
					isName = true
				}
			} else if isName && mainFileName == "" {
				mainFileName = tmp
			}
		}
	}
	return mainFileName, files, sizes
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

	row := doc.Find("code[x-ref=\"mediainfo\"]").Text()
	GENERAL := 1
	VIDEO := 2
	AUDIO := 3
	section := GENERAL
	for _, line := range strings.Split(row, "\n") {
		tmp := strings.TrimSpace(line)
		if tmp != "" {
			if strings.HasPrefix(tmp, "Text") {
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
		parts := strings.Split(line, "|")
		if len(parts) > 0 {
			studio := strings.TrimSpace(parts[len(parts)-1])
			if !strings.HasSuffix(studio, "bps") {
				audio.Studio = studio
			}
		}
	}
}

func getValue(line string, prefix string) (bool, string) {
	if strings.HasPrefix(line, prefix) {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			return true, strings.TrimSpace(parts[1])
		} else {
			return true, "Не вказано"
		}
	}
	return false, ""
}

package listeners

import (
	"fmt"
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"html"
	"log"
	"strconv"
	"strings"
	"utopia-client/entities"
	"utopia-client/page_navigator"
	"utopia-client/parsers"
)

func GetTorrentFromUtopia() {
	processor := func(data []byte) {
		log.Printf("[GetTorrentFromUtopia] Отримано повідомлення з NATS: %s", string(data))
		chatId, args, err := parseCommand(data)
		if err != nil {
			log.Printf("[GetTorrentFromUtopia] Проблема при парсингу повідомлення: %s", err)
			return
		}
		log.Printf("[GetTorrentFromUtopia] Парсинг повідомлення: chatID = %d, message = %s", chatId, args)

		if len(args) != 1 {
			log.Printf("[GetTorrentFromUtopia] Проблема при спробі отримання аргументів. Їх кількість не відповідає вимогам: має бути 1 а не %d", len(args))
			return
		}

		if torrent, noErrors := getTorrentInfo(args[0]); noErrors {
			messageToUser := prepareTorrentInfoMessage(torrent)
			nats_helper.SendMessageToUser(chatId, messageToUser)
		}
	}

	listener := &nats_helper.NatsListenerHandler{
		Function: processor,
	}

	err := nats_helper.StartNatsListener("UTOPIA_GET_TORRENT_INFO", listener)
	if err != nil {
		log.Printf("[GetTorrentFromUtopia] Проблема при спробі стартувати лісинер UTOPIA_GET_TORRENT_INFO")
	}
}

func getTorrentInfo(torrentIdStr string) (*entities.Torrent, bool) {
	torrentId, err := strconv.Atoi(torrentIdStr)
	if err != nil || torrentId <= 0 {
		log.Printf("[GetTorrentFromUtopia] Invalid torrent ID")
		return nil, false
	}
	page, err := page_navigator.GetTorrentPage(torrentId)
	if err != nil {
		log.Printf("[GetTorrentFromUtopia] Проблема при спробі отримання HTML з Utopia")
		return nil, false
	}

	torrent, err := parsers.ParseTorrentPage(torrentId, page)
	if err != nil {
		log.Printf("[GetTorrentFromUtopia] Проблема парсингу сторінки торента з Utopia")
		return nil, false
	}
	return torrent, true
}

func prepareTorrentInfoMessage(torrent *entities.Torrent) string {
	var line strings.Builder
	line.WriteString(fmt.Sprintf("[%s] <b>%s</b>\n\n", html.EscapeString(torrent.Category), html.EscapeString(torrent.Title)))
	if torrent.FileCount > 0 {
		line.WriteString(fmt.Sprintf("Файлів на закачку: %d. ", torrent.FileCount))
	}
	resolution, bitrate := "Невідомо", "Невідомо"
	if torrent.File.Media != nil {
		resolution = torrent.File.Media.Resolution
		bitrate = torrent.File.Media.OverallBitRate
	}
	line.WriteString(fmt.Sprintf("%s -- %s -- %s\n", html.EscapeString(torrent.File.Size), html.EscapeString(resolution), html.EscapeString(bitrate)))
	line.WriteString("Переклад: ")
	translate := ""
	for _, audio := range torrent.File.Audio {
		if audio.Language == "Ukrainian" && audio.Studio != "" {
			if strings.Contains(translate, audio.Studio) {
				continue
			}
			if translate != "" {
				translate = translate + " | "
			}
			if audio.Default {
				translate = translate + "✓"
			}
			translate = translate + html.EscapeString(audio.Studio)
		}
	}
	if translate == "" {
		translate = "Невідомо"
	}
	line.WriteString(fmt.Sprintf("%s\n\n", translate))
	if torrent.SimilarId != "" {
		line.WriteString(fmt.Sprintf("Схожі торенти: /utp_siml_%s\n", strings.ReplaceAll(torrent.SimilarId, ".", "_")))
	}
	line.WriteString(fmt.Sprintf("Скачати: /utp_dwnl_%d\n", torrent.ID))

	return line.String()
}

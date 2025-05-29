package listeners

import (
	"fmt"
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"utopia-client/page_navigator"
)

func DownloadUtopiaTorrent() {
	processor := func(data []byte) {
		log.Printf("[DownloadUtopiaTorrent] Отримано повідомлення з NATS: %s", string(data))
		chatId, args, err := nats_helper.ParseNatsBotCommand(data)
		if err != nil {
			log.Printf("[DownloadUtopiaTorrent] Проблема при парсингу повідомлення: %s", err)
			return
		}
		log.Printf("[DownloadUtopiaTorrent] Парсинг повідомлення: chatID = %d, message = %s", chatId, args)

		if len(args) != 1 {
			log.Printf("[DownloadUtopiaTorrent] Проблема при спробі отримання аргументів. Їх кількість не відповідає вимогам: має бути 1 а не %d", len(args))
			return
		}

		if noErrors := getTorrentFile(args[0]); noErrors {
			nats_helper.SendMessageToUser(chatId, "Торент файл скачано. Починаэмо його обробку для подальшого завантаження")
		}
	}

	listener := &nats_helper.NatsListenerHandler{
		Function: processor,
	}

	if err := nats_helper.StartNatsListener("UTOPIA_GET_TORRENT_FILE", listener); err != nil {
		log.Printf("[DownloadUtopiaTorrent] Проблема при спробі стартувати лісинер UTOPIA_GET_TORRENT_INFO")
	}
}

func getTorrentFile(torrentIdStr string) bool {
	torrentId, err := strconv.Atoi(torrentIdStr)
	if err != nil {
		log.Printf("[DownloadUtopiaTorrent] Invalid torrent ID: %v", err)
	}
	data, err := page_navigator.GetTorrentFile(torrentId)
	if err != nil {
		log.Printf("[DownloadUtopiaTorrent] Проблема при спробі отримання HTML з Utopia")
		return false
	}

	//dstDir := "/foo/bar"
	dstDir := "."
	filePath := filepath.Join(dstDir, fmt.Sprintf("%d.torrent", torrentId))
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("[DownloadUtopiaTorrent] Не вдалось створити файл %s: %w", filePath, err)
		return false
	}
	defer outFile.Close()

	if err := os.WriteFile(filePath, data, 0777); err != nil {
		log.Printf("[DownloadUtopiaTorrent] Не вдалось створити файл %s: %v", filePath, err)
		return false
	}

	return true
}

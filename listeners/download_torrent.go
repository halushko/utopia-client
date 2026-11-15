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
		chatId, args, err := parseCommand(data)
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
	if err != nil || torrentId <= 0 {
		log.Printf("[DownloadUtopiaTorrent] Invalid torrent ID")
		return false
	}
	data, err := page_navigator.GetTorrentFile(torrentId)
	if err != nil {
		log.Printf("[DownloadUtopiaTorrent] Проблема при спробі отримання HTML з Utopia")
		return false
	}

	if err := writeTorrentFile("/root/torrent_files", torrentId, data); err != nil {
		log.Printf("[DownloadUtopiaTorrent] Не вдалось зберегти торрент: %v", err)
		return false
	}
	return true
}

// Publish to the watched directory only after the complete file is closed.
func writeTorrentFile(dir string, id int, data []byte) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	out, err := os.CreateTemp(dir, ".utopia-*.tmp")
	if err != nil {
		return err
	}
	name := out.Name()
	defer os.Remove(name)
	if _, err = out.Write(data); err != nil {
		out.Close()
		return err
	}
	// Match the existing shared-directory read access before publication.
	if err = out.Chmod(0644); err != nil {
		out.Close()
		return err
	}
	if err = out.Sync(); err != nil {
		out.Close()
		return err
	}
	if err = out.Close(); err != nil {
		return err
	}
	return os.Rename(name, filepath.Join(dir, fmt.Sprintf("%d.torrent", id)))
}

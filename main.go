package main

import (
	"github.com/halushko/kino-cat-core-go/logger_helper"
	"utopia-client/listeners"
)

func main() {
	logger_helper.SoftPrepareLogFile()

	listeners.GetTorrentFromUtopia()
	listeners.DownloadUtopiaTorrent()

	select {}
}

package main

import (
	downloader "ytdl-tgbot"
)

func main() {
	log := downloader.NewSlogLogger()

	downloader.StartBot(log)
}

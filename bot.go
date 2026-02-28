package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
)

func StartBot(log Logger) {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal("couldn't create new bot API", err)
	}

	log.Info(fmt.Sprintf("authorized on account %q", bot.Self.UserName))

	catchUpdates(log, bot)
}

func catchUpdates(log Logger, bot *tgbotapi.BotAPI) {

	log.Info("waiting for links from now on")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Fatal("couldn't get updates channel", err)
	}

	for update := range updates {
		if update.Message != nil {

			switch {
			case update.Message.Text == "/start":
				handleStart(log, bot, &update)
			default:
				handleAudioDownload(log, bot, &update)
			}
		}
	}
}

func handleStart(log Logger, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi! Send me the link to YT video and I'll download its audio")
	_, err := bot.Send(msg)
	if err != nil {
		log.Error(fmt.Sprintf("error sending /start message: %v", err))
	}
}

func handleAudioDownload(log Logger, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	videoLink := update.Message.Text
	filepath, err := DownloadAudioFromVideo(log, videoLink)

	if err != nil {
		respondWithErr(log, bot, update, err)
		return
	}

	go func() {
		err := deleteFile(filepath, log)
		if err != nil {
			log.Error(fmt.Sprintf("error during audio file deletion: %v", err))
		}
	}()

	audio := tgbotapi.NewAudioUpload(update.Message.Chat.ID, filepath)

	_, err = bot.Send(audio)
	if err != nil {
		log.Error(fmt.Sprintf("error during audio sending: %v", err))
		respondWithErr(log, bot, update, fmt.Errorf("error during audio sending. try again"))
	}
}

func respondWithErr(log Logger, bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
	msg.ReplyToMessageID = update.Message.MessageID

	_, sendErr := bot.Send(msg)

	if sendErr != nil {
		log.Error(fmt.Sprintf("%v", sendErr))
	}
}

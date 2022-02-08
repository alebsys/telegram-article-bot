package main

import (
	"log"
	"os"

	"github.com/alebsys/telegram-article-bot/internal/devto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	descp = "`Request example::\n/article go 10 5\nгде:\n* go - topic (tag);\n* 10 - search period in days;\n* 5 - number of posts.`"
)

func main() {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic("getting TELEGRAM_APITOKEN: ", err)
	}
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.EditedMessage != nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "markdown"
		msg.DisableWebPagePreview = true

		input := update.Message.Text

		log.Printf("[%s] %s", update.Message.From.UserName, input)

		switch update.Message.Command() {
		case "help":
			msg.Text = "`Hello! I can find articles of interest to you on DEV.TO\n\n`" + descp
		case "article":
			note := "`Enter the correct command!\n\n`" + descp

			b := devto.ParseInput(input)
			if !b {
				log.Print(err)
				msg.Text = note
				break
			}

			query := devto.NewQuery(input)
			articles, err := devto.GetArticles(query.Tag, query.Freshness)
			if err != nil {
				log.Panic(err)
			}

			msg.Text = articles.WriteArticles(query.Limit)
		default:
			msg.Text = "`I don't know this command. Enter /help`"
		}

		bot.Send(msg)
	}

}

package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type user struct {
	tgid          int64
	tg_username   string
	dialog_status int64
}

// main database for dialogs, key (int64) is telegram user id
var userDatabase = make(map[int64]user)

var msgTemplates = make(map[string]string)

func main() {

	msgTemplates["hello"] = "Hey, this bot is doing something, you're added to the database."
	msgTemplates["case0"] = "You are in case 0"
	msgTemplates["case1"] = "You are in case 1"

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {

		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message != nil {

			// check user in database
			if _, ok := userDatabase[update.Message.From.ID]; !ok {

				userDatabase[update.Message.From.ID] = user{update.Message.Chat.ID, update.Message.Chat.UserName, 0}
				msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, msgTemplates["hello"])
				bot.Send(msg)

			} else {

				switch userDatabase[update.Message.From.ID].dialog_status {
				case 0:

					if updateDB, ok := userDatabase[update.Message.From.ID]; ok {

						id := userDatabase[update.Message.From.ID].tgid
						username := userDatabase[update.Message.From.ID].tg_username

						infoMsg := fmt.Sprintf("%v, %v with id = %v.", msgTemplates["case0"], username, id)

						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, infoMsg)
						bot.Send(msg)

						msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "There is magic happening.")
						bot.Send(msg)

						updateDB.dialog_status = 1

						userDatabase[update.Message.From.ID] = updateDB

					}

				case 1:

					if updateDB, ok := userDatabase[update.Message.From.ID]; ok {

						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, msgTemplates["case1"])
						bot.Send(msg)

						updateDB.dialog_status = 0

						userDatabase[update.Message.From.ID] = updateDB

					}

				} // end of switch dialog_Status

			} // end of "check user in database"

		}

	}

} // end of main

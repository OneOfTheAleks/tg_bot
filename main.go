package main

import (
	"TgBot/sq"
	"errors"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	dbs := sq.InitDB("data.sqlite")
	dbs.CreateTable()

	/* 	if err != nil {
		fmt.Println(err)
		return
	} */
	defer dbs.Close()

	bot, err := tgbotapi.NewBotAPI("MyToken")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			//	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			id := update.Message.MessageID
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			if !getUserName(update.Message.Text) {
/* 				
                msg.ReplyToMessageID = update.Message.MessageID
				//update.Message.Text =
				msg.Text = "А вот и не правильно! Недостаточно параметров: надо '@at_save_bot сохрани/покажи ИмяТэга сообщение'"
				bot.Send(msg) */

				continue
			}

			//fmt.Println(update.Message.Text)
			command, tagStr, msgStr, err := splitStr(update.Message.Text)
			if err != nil {
				msg.ReplyToMessageID = update.Message.MessageID
				//update.Message.Text =
				msg.Text = err.Error()
				bot.Send(msg)
			}
			if err == nil {

				ms, err := dbs.ManageMessage(command, tagStr, msgStr, id)
				if err != nil {
					msg.Text = "Что-то пошло не так: " + err.Error()
				} else {
					if command == "сохрани" {
						msg.Text = "Спасибо, сохранил! "
					}
					if command == "покажи" {
						msg.Text = ms
					}
				}
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}

			//fmt.Println(msg)
		}
	}

}

func splitStr(sourceStr string) (command string, tagStr string, msgStr string, err error) {
	words := strings.Fields(sourceStr)
	/* for idx, word := range words {
		fmt.Printf("Word %d is: %s\n", idx, word)
	} */
	if len(words) < 3 {
		err = errors.New("А вот и не правильно! Недостаточно параметров: надо '@at_save_bot сохрани/покажи ИмяТэга сообщение'")
		return
	}
	command = words[1]
	tagStr = words[2]
	mArray := words[3:]

	if command != "сохрани" && command != "покажи" {
		err = errors.New("А вот и не правильно! Не те параметры: надо '@at_save_bot сохрани/покажи ИмяТэга сообщение'")
		return
	}

	for _, ms := range mArray {
		msgStr = msgStr + " " + ms
	}

	return

}

func getUserName(msgStr string) bool {
	words := strings.Fields(msgStr)
	if len(words) < 3 {
		return false
	}
	tbName := words[0]
	return tbName == "@at_save_bot"

}

package main

import (
	"fmt"
	"log"
	"time"

	tb "github.com/tucnak/telebot"
)

func sendMessage(bot *tb.Bot, msg *message) error {
	maxTries := 5
	// message ID=xxx, DATE=2006-01-02 15:04:05, TYPE=Album/Text\n
	fmt.Printf("message ID=%s, DATE=%s, ", msg.ID, date(msg.Time))
	if msg.isAlbum() {
		fmt.Printf("TYPE=Album\n") // End line
	} else {
		fmt.Printf("TYPE=Text\n") // End line
	}

	for i := 0; i < maxTries; i++ {
		if i > 0 {
			wt := 2 << i
			log.Println("Sleep", wt, "minutes | sendMessage")
			time.Sleep(time.Duration(wt) * time.Minute)
		}

		err := sendAlbumOrText(bot, msg)
		if err != nil {
			continue
		}

		return nil
	}

	return fmt.Errorf("Failed to send message %s at %d tries", msg.ID, maxTries)
}

func sendAlbumOrText(bot *tb.Bot, m *message) error {
	var msg *tb.Message
	if m.isAlbum() {
		ml, err := bot.SendAlbum(ch, m.Album())
		if err != nil {
			return err
		}

		msg = &ml[0]
	} else {
		m1, err := bot.Send(ch, m.mainText())
		if err != nil {
			return err
		}

		msg = m1
	}
	for _, s := range m.otherTexts() {
		err := sendSlice(bot, msg, s)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func sendSlice(bot *tb.Bot, m *tb.Message, text string) error {
	maxTries := 5
	for i := 0; i < maxTries; i++ {
		if i > 0 {
			wt := 2 << i
			log.Println("Sleep", wt, "minutes | sendSlice")
			time.Sleep(time.Duration(wt) * time.Minute)
		}

		_, err := bot.Reply(m, text)
		if err != nil {
			continue
		}

		return nil
	}

	return fmt.Errorf("ERROR: Failed to send slice after %d tries", maxTries)
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"time"

	tb "github.com/tucnak/telebot"
)

var (
	ch *channel
)

func main() {
	cfg := loadConfig("./config.json")
	ch = &channel{Username: cfg.ChannelUsername}

	bot, _ := tb.NewBot(tb.Settings{
		Token:  cfg.BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	msgs, _ := collectAndSort(cfg.DataDIR)
	max := len(msgs)

	started := false
	if cfg.NotBefore == "" {
		started = true
	}

	notice, _ := bot.Send(ch, "Go! #临时消息")
	bot.Pin(notice)
	for i, m := range msgs {
		if m.ID == cfg.NotBefore {
			started = true
		}
		if !started {
			continue
		}

		info := fmt.Sprintf("%.3d of %.3d, %s\n", i, max, m.ID)
		log.Println(info)

		jobInfo := mkMesage(i+1, len(msgs), m.ID, nil)
		bot.Edit(notice, jobInfo, &tb.SendOptions{
			DisableWebPagePreview: true,
		})

		err := sendMessage(bot, &m)
		if err != nil {
			jobInfo = mkMesage(i+1, len(msgs), m.ID, err)
			bot.Edit(notice, jobInfo, &tb.SendOptions{
				DisableWebPagePreview: true,
			})

			log.Printf("failed: %v\n", err)
			time.Sleep(time.Duration(5 * 60))
			err2 := sendMessage(bot, &m)
			if err2 != nil {
				jobInfo = mkMesage(i+1, len(msgs), m.ID, err2)
				bot.Edit(notice, jobInfo, &tb.SendOptions{
					DisableWebPagePreview: true,
				})

				log.Fatalln(err)
			}
		}
	}

	fmt.Println("done")
	return
}

func mkMesage(num, total int, id string, err error) string {
	return fmt.Sprintf(`#临时消息 搬运中

当前第（%d/%d）条 ID: %s

链接：%s

%s`,
		num, total, id,
		fmt.Sprintf("https://web.okjike.com/originalPost/%s", id),
		optionalErrorMessage(err),
	)
}

func optionalErrorMessage(err error) string {
	if err != nil {
		return fmt.Sprintf("\n遇到一个问题，正在重试：%s\n", err.Error())
	}

	return ""
}

func downloadFile(dstPath, remoteURL string) error {
	fmt.Printf("\nDownloading... %s << %s", dstPath, remoteURL)
	resp, err := http.Get(remoteURL)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = io.Copy(dstFile, resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func checkAndDownloadPictures(dirpath string, msg message) []string {
	picURLs := []string{}
	for _, pic := range msg.Pictures {
		filename := pic.filename()
		dstPath := path.Join(dirpath, filename)
		picURLs = append(picURLs, dstPath)
		if fileNotExists(dstPath) {
			downloadFile(dstPath, pic.PicURL)
		}
	}
	return picURLs
}

func collectAndSort(profileDir string) ([]message, error) {
	messages := []message{}

	dirs := listDirs(profileDir)
	for _, d := range dirs {
		msg, err := workOnDir(d)
		if err != nil {
			return messages, err
		}
		msg.parseTime()
		localPics := checkAndDownloadPictures(d, msg)
		msg.localPics = localPics

		messages = append(messages, msg)
	}

	ms := messageSorter{
		messages: messages,
	}
	sort.Sort(ms)

	return ms.messages, nil
}

func listDirs(root string) []string {
	items := []string{}

	folder, err := os.Open(root)
	if err != nil {
		log.Println(err)
	}
	defer folder.Close()

	subFolders, err := folder.Readdir(-1)
	if err != nil {
		log.Println(err)
	}

	for _, subPath := range subFolders {
		if subPath.IsDir() {
			items = append(items, path.Join(root, subPath.Name()))
		}
	}

	return items
}

func workOnDir(dirpath string) (message, error) {
	msg := message{}

	dataPath := path.Join(dirpath, "data.json")
	dataFile, err := os.Open(dataPath)
	if err != nil {
		log.Println(err)
		return msg, err
	}
	defer dataFile.Close()

	doc := json.NewDecoder(dataFile)
	err = doc.Decode(&msg)
	if err != nil {
		log.Println(err)
		return msg, err
	}

	return msg, nil
}

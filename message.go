package main

import (
	"log"
	"math"
	"net/url"
	"path"
	"strings"
	"time"

	tb "github.com/tucnak/telebot"
)

type message struct {
	ID         string           `json:"id"`
	ActionTime string           `json:"actionTime"` //  "2018-06-21T02:25:27.000Z"
	CreatedAt  string           `json:"createdAt"`  // "2018-06-21T02:25:27.628Z"
	Time       time.Time        `json:"-"`
	Type       string           `json:"type"` // ORIGINAL_POST
	Content    string           `json:"content"`
	Status     string           `json:"status"` // NORMAL
	Topic      messageTopic     `json:"topic"`
	Pictures   []messagePicture `json:"pictures"`

	localDir  string
	localPics []string

	likeCount    int64 // 6
	commentCount int64 // 0
	repostCount  int64 // 0
	shareCount   int64 // 0
}

func (m *message) isAlbum() bool {
	return len(m.localPics) > 0
}

func (m *message) Album() []tb.InputMedia {
	album := tb.Album{}
	for ix, localPic := range m.localPics {
		tbFile := tb.FromDisk(localPic)
		photo := &tb.Photo{
			File: tbFile,
		}
		if ix == 0 {
			photo.Caption = m.mainText()
		}
		album = append(album, photo)
	}

	return album
}

func (m *message) parseTime() {
	ta, _ := time.Parse(time.RFC3339, m.ActionTime)
	tc, _ := time.Parse(time.RFC3339, m.CreatedAt)
	if math.Abs(float64(ta.Unix()-tc.Unix())) > 1000 {
		log.Println("TIME_NOT_EQUAL", m.ActionTime, m.CreatedAt)
	}

	t, err := time.Parse(time.RFC3339, m.ActionTime)
	if err != nil {
		log.Println(err, "Failed to format date", m.ActionTime, "as time.RFC3339")
		return
	}

	m.Time = t
}

func (m *message) Tags() []string {
	tags := []string{
		m.ID,
		m.Time.Format("2006-01-02 15:04:05"),
	}
	if len(m.Topic.ID) > 0 {
		tags = append(tags, "#"+m.Topic.Content)
	}

	return tags
}

func (m *message) Text() string {
	lines := []string{
		m.Content,
	}
	tags := strings.Join(m.Tags(), " ")
	if len(tags) > 3 {
		lines = append(lines, tags)
	}

	return strings.Join(lines, "\n\n")
}

func (m *message) mainText() string {
	return splitTexts(m.Text(), 1600)[0]
}

func (m *message) otherTexts() []string {
	return splitTexts(m.Text(), 1600)[1:]
}

type messageTopic struct {
	ID      string `json:"id"`
	Content string `json:"content"` // "随手拍张照"
	Type    string `json:"topic"`   // TOPIC
}

type messagePicture struct {
	PicURL string `json:"picUrl"`
	Format string `json:"format"` // "jpeg"
	Width  int64  `json:"width"`
	Height int64  `json:"height"`

	thumbnailURL    string
	smallPicURL     string
	middlePicURL    string
	watermarkPicURL string
}

func (mp *messagePicture) filename() string {
	wURL, err := url.Parse(mp.PicURL)
	if err != nil {
		log.Panicln(err)
	}

	return path.Base(wURL.Path)
}

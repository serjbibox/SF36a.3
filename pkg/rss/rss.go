package rss

import (
	"bytes"
	"encoding/xml"
	"errors"
	"log"
	"strings"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/serjbibox/SF36a.3/pkg/models"
	"github.com/serjbibox/SF36a.3/pkg/storage"
)

// Обработчик RSS запросов сервера GoNews
type Rss struct {
	Channel Channel `xml:"channel"`
	Storage *storage.Storage
	Link    string
	Hash    models.Hash
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Items       []Item `xml:"item"`
}

//Конструктор объекта RSS
func New(s *storage.Storage, link string) (*Rss, error) {
	if s == nil {
		return nil, errors.New("storage is nil")
	}

	return &Rss{
		Storage: s,
		Link:    link,
		Hash: models.Hash{
			Link: link,
		},
	}, nil
}

//Чтение источника и запись в канал
func (r *Rss) ParseRss(period int, news chan<- []models.Post, errs chan<- error) {
	for {
		posts, err := r.parse()
		switch {
		case err != nil:
			errs <- err
		default:
			news <- posts
		}
		time.Sleep(time.Duration(period) * time.Second)
	}
}

// Parse читает rss-поток; возвращет
// массив раскодированных новостей и ошибку.
func (r *Rss) parse() ([]models.Post, error) {
	XMLdata, err := readRssBody(r.Link)
	if err != nil {
		return nil, err
	}
	ok, err := r.isHashEqual(XMLdata)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			err := r.hashInit(XMLdata)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}

	}
	log.Println("Parse hash check", ok, "url:", r.Link)
	if ok {
		return nil, errors.New("XML file not changed, no new posts")
	}
	buffer := bytes.NewBuffer(XMLdata)
	decoded := xml.NewDecoder(buffer)
	err = decoded.Decode(r)
	if err != nil {
		return nil, err
	}
	log.Println("Channel PubDate:", r.Channel.PubDate)
	var data []models.Post
	for _, item := range r.Channel.Items {
		var p models.Post
		p.Title = item.Title
		p.Content = item.Description
		p.Content = strip.StripTags(p.Content)
		p.Content = strings.ReplaceAll(p.Content, "&rarr;", "")
		p.Link = item.Link
		item.PubDate = strings.ReplaceAll(item.PubDate, ",", "")
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", item.PubDate)
		}
		if err == nil {
			p.PubTime = t.Unix()
		}
		data = append(data, p)
	}
	r.Channel.PubDate = strings.ReplaceAll(r.Channel.PubDate, ",", "")
	t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", r.Channel.PubDate)
	if err != nil {
		t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", r.Channel.PubDate)
	}
	if err == nil {
		r.Hash.PubTime = t.Unix()
	}
	if r.Hash.PubTime == 0 {
		r.Hash.PubTime = time.Now().Unix()
	}
	err = r.hashUpdate()
	if err != nil {
		return nil, err
	}
	log.Println("Hash updated", r.Link)

	return data, nil
}

//Создаёт первоначальные записи хэшей XML
func (r *Rss) hashInit(data []byte) error {
	var err error
	r.Hash.NewsHash, err = getMd5Hash(data)
	if err != nil {
		return errors.New("HashInit().getMd5Hash error: " + err.Error())
	}
	err = r.Storage.Hash.Create(r.Hash)
	if err != nil && !strings.Contains(err.Error(), "SQLSTATE 23505") {
		return errors.New("HashInit().Create error: " + err.Error())
	}
	return nil
}

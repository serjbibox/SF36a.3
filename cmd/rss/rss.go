package rss

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/serjbibox/SF36a.3/pkg/models"
	"github.com/serjbibox/SF36a.3/pkg/storage"
)

// Обработчик RSS запросов сервера GoNews
type Rss struct {
	Channel Channel `xml:"channel"`
	storage *storage.Storage
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
	Items       []Item `xml:"item"`
}

//Конструктор объекта RSS
func New(s *storage.Storage) (*Rss, error) {
	if s == nil {
		return nil, errors.New("storage is nil")
	}
	return &Rss{storage: s}, nil
}

func (r *Rss) Store(news []models.Post) error {
	err := r.storage.Create(news)
	if err != nil && !strings.Contains(err.Error(), "SQLSTATE 23505") {
		return err
	}
	return nil
}

// Parse читает rss-поток и возвращет
// массив раскодированных новостей.
//"https://habr.com/ru/rss/hub/go/all/?fl=ru"
func (r *Rss) Parse(url string) ([]models.Post, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	XMLdata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(XMLdata)
	decoded := xml.NewDecoder(buffer)
	err = decoded.Decode(r)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Title : %s\n", r.Channel.Title)
	fmt.Printf("Description : %s\n", r.Channel.Description)
	fmt.Printf("Link : %s\n", r.Channel.Link)
	var data []models.Post
	for _, item := range r.Channel.Items {
		var p models.Post
		p.Title = item.Title
		p.Content = item.Description
		p.Content = strip.StripTags(p.Content)
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
		//fmt.Printf("[%d] item title : %s\n", i, item.Title)
		//fmt.Printf("[%d] item description : %s\n", i, strip.StripTags(item.Description))
		//fmt.Printf("[%d] item link : %s\n\n", i, item.Link)
		//fmt.Printf("[%d] item pubDate : %s\n\n", i, item.PubDate)
	}
	return data, nil
	/*
		fmt.Printf("Title : %s\n", rss.Channel.Title)
		fmt.Printf("Description : %s\n", rss.Channel.Description)
		fmt.Printf("Link : %s\n", rss.Channel.Link)

		total := len(rss.Channel.Items)

		fmt.Printf("Total items : %v\n", total)

		for i := 0; i < total; i++ {
			fmt.Printf("[%d] item title : %s\n", i, rss.Channel.Items[i].Title)
			fmt.Printf("[%d] item description : %s\n", i, rss.Channel.Items[i].Description)
			fmt.Printf("[%d] item link : %s\n\n", i, rss.Channel.Items[i].Link)
		}
	*/
}

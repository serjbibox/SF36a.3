package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/serjbibox/SF36a.3/pkg/handler"
	"github.com/serjbibox/SF36a.3/pkg/models"
	"github.com/serjbibox/SF36a.3/pkg/rss"
	"github.com/serjbibox/SF36a.3/pkg/server"
	"github.com/serjbibox/SF36a.3/pkg/storage"
	"github.com/serjbibox/SF36a.3/pkg/storage/postgresql"
)

// Конфигурация приложения
type config struct {
	Resurses       []string                  `json:"rss"`
	Period         int                       `json:"request_period"`
	PostgresConfig postgresql.PostgresConfig `json:"postgres_settings"`
}

var ctx = context.Background()

func main() {
	c, err := readConfig("./cmd/news/config.json")
	if err != nil {
		log.Fatal(err)
	}
	connString, err := postgresql.GetConnectionString(c.PostgresConfig)
	if err != nil {
		log.Fatal(err)
	}
	db, err := postgresql.New(connString)
	if err != nil {
		log.Fatal(err)
	}
	s, err := storage.NewStoragePostgres(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	news, errs := parseRss(c, s)
	go storeNews(news, s)
	go errHandler(errs)
	handlers, err := handler.New(s)
	if err != nil {
		log.Fatal(err)
	}
	srv := new(server.Server)
	log.Fatal(srv.Run(server.HTTP_PORT, handlers.InitRoutes()))
}

//Сохраняет полученные публикации из канала в БД
func storeNews(news <-chan []models.Post, s *storage.Storage) {
	for n := range news {
		err := s.Post.Create(n)
		if err != nil {
			log.Println("s.Post.Create() error:", err)
		}
	}
}

//Обработчик ошибок
func errHandler(errs <-chan error) {
	for err := range errs {
		log.Println("ParseRss() error:", err)
	}
}

//Чтение публикаций из RSS рассылок в отдельных потоках для каждого ресурса
func parseRss(c *config, s *storage.Storage) (<-chan []models.Post, <-chan error) {
	news := make(chan []models.Post)
	errs := make(chan error)
	for _, url := range c.Resurses {
		r, err := rss.New(s, url)
		if err != nil {
			log.Fatal(err)
		}
		go r.ParseRss(c.Period, news, errs)
	}
	return news, errs
}

//Чтение JSON файла конфигурации
func readConfig(path string) (*config, error) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config config
	err = json.Unmarshal(c, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/serjbibox/SF36a.3/cmd/rss"
	"github.com/serjbibox/SF36a.3/pkg/handler"
	"github.com/serjbibox/SF36a.3/pkg/models"
	"github.com/serjbibox/SF36a.3/pkg/storage"
	"github.com/serjbibox/SF36a.3/pkg/storage/postgresql"
)

const (
	HTTP_PORT = "8080"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

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

	news := make(chan []models.Post)
	errs := make(chan error)
	var rssArray []*rss.Rss
	for idx, url := range c.Resurses {
		r, err := rss.New(s, url)
		if err != nil {
			log.Fatal(err)
		}
		rssArray = append(rssArray, r)
		go rssArray[idx].ParseRss(c.Period, news, errs)
	}
	go func() {
		for n := range news {
			err := s.Post.Create(n)
			if err != nil {
				log.Println("s.Post.Create() error:", err)
			}
		}
	}()

	go func() {
		for err := range errs {
			log.Println("ParseRss() error:", err)
		}
	}()
	handlers, err := handler.New(s)
	if err != nil {
		log.Fatal(err)
	}
	srv := new(Server)
	log.Fatal(srv.Run(HTTP_PORT, handlers.InitRoutes()))
}

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

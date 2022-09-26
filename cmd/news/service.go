package main

import (
	"log"
	"net/http"
	"time"

	"github.com/serjbibox/SF36a.3/pkg/handler"
	"github.com/serjbibox/SF36a.3/pkg/storage"
	"github.com/serjbibox/SF36a.3/pkg/storage/memdb"
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

func main() {
	db, err := memdb.New()
	if err != nil {
		log.Fatal(err)
	}
	s, err := storage.NewStorageMemDb(db)
	if err != nil {
		log.Fatal(err)
	}
	handlers, err := handler.New(s)
	if err != nil {
		log.Fatal(err)
	}
	srv := new(Server)
	log.Fatal(srv.Run(HTTP_PORT, handlers.InitRoutes()))
}

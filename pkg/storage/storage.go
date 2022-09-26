package storage

import (
	"errors"
	"log"
	"os"

	"github.com/serjbibox/SF36a.3/pkg/models"
	"github.com/serjbibox/SF36a.3/pkg/storage/memdb"
)

var elog = log.New(os.Stderr, "Storage error\t", log.Ldate|log.Ltime|log.Lshortfile)
var ilog = log.New(os.Stdout, "Storage info\t", log.Ldate|log.Ltime)

//задаёт контракт на работу с таблицей публикаций БД.
type Post interface {
	GetAll() ([]models.Post, error)             // получение всех публикаций
	GetByQuantity(n int) ([]models.Post, error) // получение публикаций по заданному количеству
	Create(models.Post) (id int, err error)     // создание новой публикации
	Update(models.Post) error                   // обновление публикации
	Delete(id string) error                     // удаление публикации по ID
}

// Хранилище данных.
type Storage struct {
	Post
}

// Конструктор объекта хранилища для БД MemDb.
func NewStorageMemDb(db memdb.DB) (*Storage, error) {
	if db == nil {
		elog.Println("db is nil")
		return nil, errors.New("db is nil")
	}
	return &Storage{
		Post: newPostMemDb(db),
	}, nil
}

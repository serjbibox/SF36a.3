package storage

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/serjbibox/SF36a.3/pkg/models"
	"github.com/serjbibox/SF36a.3/pkg/storage/memdb"
)

var elog = log.New(os.Stderr, "Storage error\t", log.Ldate|log.Ltime|log.Lshortfile)
var ilog = log.New(os.Stdout, "Storage info\t", log.Ldate|log.Ltime)

//задаёт контракт на работу с таблицей публикаций БД.
type Post interface {
	GetAll() ([]models.Post, error)             // получение всех публикаций
	GetByQuantity(n int) ([]models.Post, error) // получение публикаций по заданному количеству
	Create([]models.Post) error                 // создание новой публикации
	Update(models.Post) error                   // обновление публикации
	Delete(id string) error                     // удаление публикации по ID
}

// Хранилище данных.
type Storage struct {
	Post
}

// Конструктор объекта хранилища для БД PostgreSQL.
func NewStoragePostgres(ctx context.Context, db *pgxpool.Pool) (*Storage, error) {
	if ctx == nil {
		elog.Println("context is nil")
		return nil, errors.New("context is nil")
	}
	if db == nil {
		elog.Println("db is nil")
		return nil, errors.New("db is nil")
	}
	return &Storage{
		Post: newPostPostgres(ctx, db),
	}, nil
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

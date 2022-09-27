package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/serjbibox/SF36a.3/pkg/models"
	"github.com/serjbibox/SF36a.3/pkg/storage/memdb"
)

//Интерфейс для работы с таблицей хэшей БД.
type Hash interface {
	GetByLink(string) (string, error) // получение хэша по url
	Create(models.Hash) error         // создание новой записи
	Update(models.Hash) error         // обновление записи
}

//Интерфейс для работы с таблицей публикаций БД.
type Post interface {
	GetByQuantity(n int) ([]models.Post, error) // получение публикаций по заданному количеству
	Create([]models.Post) error                 // создание новой публикации
}

// Хранилище данных.
type Storage struct {
	Post
	Hash
}

// Конструктор объекта хранилища для БД PostgreSQL.
func NewStoragePostgres(ctx context.Context, db *pgxpool.Pool) (*Storage, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &Storage{
		Post: newPostPostgres(ctx, db),
		Hash: newHashPostgres(ctx, db),
	}, nil
}

// Конструктор объекта хранилища для БД MemDb.
func NewStorageMemDb(db memdb.DB) (*Storage, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &Storage{
		Post: newPostMemDb(db),
	}, nil
}

package storage

import (
	"errors"

	"github.com/serjbibox/SF36a.3/pkg/models"
	"github.com/serjbibox/SF36a.3/pkg/storage/memdb"
)

//Объект, реализующий интерфейс работы с публикациями в памяти.
type PostMemdb struct {
	db memdb.DB
}

//Конструктор PostMemdb
func newPostMemDb(db memdb.DB) Post {
	return &PostMemdb{db: db}
}

// Получение публикаций по заданному количеству
func (s *PostMemdb) GetByQuantity(n int) ([]models.Post, error) {
	if n >= len(s.db) {
		return s.db, nil
	}
	if n < 0 {
		return nil, errors.New("n is not positive value")
	}
	id := len(s.db) - n
	out := make([]models.Post, 0)

	out = append(out, s.db[id:]...)
	return out, nil

}

// создание новой публикации
func (s *PostMemdb) Create(p []models.Post) error {
	s.db = append(s.db, p...)
	return nil
}

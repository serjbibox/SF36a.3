package storage

import (
	"errors"
	"strconv"

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

// получение всех публикаций
func (s *PostMemdb) GetAll() ([]models.Post, error) {
	return s.db, nil
}

// получение всех публикаций
func (s *PostMemdb) GetByQuantity(n int) ([]models.Post, error) {
	if n >= len(s.db) {
		return s.db, nil
	}
	id := len(s.db) - n
	out := make([]models.Post, 0)

	out = append(out, s.db[id:]...)
	return out, nil

}

// создание новой публикации
func (s *PostMemdb) Create(p models.Post) (int, error) {
	id := s.db[len(s.db)-1].ID
	p.ID = id + 1
	s.db = append(s.db, p)
	return s.db[len(s.db)-1].ID, nil
}

// обновление публикации
func (s *PostMemdb) Update(p models.Post) error {
	id := p.ID
	if id >= len(s.db) || id == 0 {
		return errors.New("wrong post id")
	}
	s.db[id-1] = p
	return nil
}

// удаление публикации по ID
func (s *PostMemdb) Delete(id string) error {
	delId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	if delId >= len(s.db) || delId == 0 {
		return errors.New("wrong post id")
	}
	s.db[delId-1] = models.Post{}
	return nil
}

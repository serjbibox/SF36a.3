package storage

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/serjbibox/SF36a.3/pkg/models"
)

//Объект, реализующий интерфейс работы с таблицей news_hash PostgreSQL.
type HashPostgres struct {
	db  *pgxpool.Pool
	ctx context.Context
}

//Конструктор PostPostgres
func newHashPostgres(ctx context.Context, db *pgxpool.Pool) Hash {
	return &HashPostgres{
		db:  db,
		ctx: ctx,
	}
}

func (s *HashPostgres) GetByLink(l string) (string, error) {
	hash := ""
	query := `
	SELECT 
		news_hash 
	FROM news_hash
	WHERE link = $1;`
	err := s.db.QueryRow(s.ctx, query, l).Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (s *HashPostgres) Create(h models.Hash) error {
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO news_hash(link, news_hash, pub_time)
		VALUES ($1, $2, $3)`,
		h.Link,
		h.NewsHash,
		h.PubTime,
	)
	if err != nil {
		return err
	}
	log.Println("hash stored")
	return nil
}

func (s *HashPostgres) Update(h models.Hash) error {
	id := 0
	query := `
	UPDATE news_hash
	SET 
		news_hash = $1,
		pub_time = $2
	WHERE link = $3
	RETURNING ID;`
	err := s.db.QueryRow(s.ctx, query, h.NewsHash, h.PubTime, h.Link).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (s *HashPostgres) Delete(id int) error {
	res, err := s.db.Exec(s.ctx, `
	DELETE FROM news_hash 
	WHERE id = $1	
	`,
		id,
	)
	if err != nil {
		return err
	}
	if res.Delete() {
		if res.String() == "DELETE 0" {
			return errors.New("no news_hash to delete")
		}
	}
	return nil
}

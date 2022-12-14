package storage

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/serjbibox/SF36a.3/pkg/models"
)

//Объект, реализующий интерфейс работы с таблицей posts PostgreSQL.
type PostPostgres struct {
	db  *pgxpool.Pool
	ctx context.Context
}

//Конструктор PostPostgres
func newPostPostgres(ctx context.Context, db *pgxpool.Pool) Post {
	return &PostPostgres{
		db:  db,
		ctx: ctx,
	}
}

// Получение публикаций по заданному количеству
func (s *PostPostgres) GetByQuantity(n int) ([]models.Post, error) {
	query := `
	SELECT 
		news.id, 
		news.title, 
		news.content, 
		news.pub_time,
		news.link
	FROM news
	ORDER BY pub_time DESC
	LIMIT $1`
	rows, err := s.db.Query(s.ctx, query, n)
	if err != nil {
		return nil, err
	}
	var news []models.Post
	for rows.Next() {
		var p models.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		news = append(news, p)
	}
	return news, rows.Err()
}

// Создание нового списка публикаций
func (s *PostPostgres) Create(p []models.Post) error {
	for _, post := range p {
		_, err := s.db.Exec(context.Background(), `
		INSERT INTO news(title, content, pub_time, link)
		VALUES ($1, $2, $3, $4)`,
			post.Title,
			post.Content,
			post.PubTime,
			post.Link,
		)
		if err != nil && !strings.Contains(err.Error(), "SQLSTATE 23505") {
			return err
		}
	}
	return nil
}

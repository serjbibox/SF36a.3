package storage

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/serjbibox/SF36a.3/pkg/models"
	"github.com/serjbibox/SF36a.3/pkg/storage/postgresql"
)

func TestPostPostgres_GetByQuantity(t *testing.T) {
	pwd := os.Getenv("DbPass")
	connString := "postgres://serj:" + pwd + "@localhost:5432/gonews?sslmode=disable"
	db, err := postgresql.New(connString)
	if err != nil {
		t.Fatal(err)
	}
	s := newPostPostgres(context.Background(), db)
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		s       Post
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "test 1",
			s:       s,
			args:    args{n: 40},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetByQuantity(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostPostgres.GetByQuantity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != (len(got) > 0) {
				t.Errorf("PostPostgres.GetByQuantity() = %v, want %v", got, tt.want)
			}
			fmt.Println("len of got", len(got))
		})
	}
}

func TestPostPostgres_Create(t *testing.T) {
	pwd := os.Getenv("DbPass")
	connString := "postgres://serj:" + pwd + "@localhost:5432/gonews?sslmode=disable"
	db, err := postgresql.New(connString)
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewStoragePostgres(context.Background(), db)
	if err != nil {
		t.Fatal(err)
	}
	//s := newPostPostgres(context.Background(), db)
	type args struct {
		p []models.Post
	}
	tests := []struct {
		name    string
		s       Post
		args    args
		wantErr bool
	}{
		{
			name: "test 1",
			s:    s.Post,
			args: args{
				[]models.Post{
					{
						ID:      4,
						Title:   "The Go Memory Model 3",
						Content: "The Go memory model specifies the conditions under which reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in a different goroutine.",
						Link:    strconv.Itoa(rand.Intn(1_000_000_000)),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Create(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("PostPostgres.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

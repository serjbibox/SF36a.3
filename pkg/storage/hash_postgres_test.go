package storage

import (
	"context"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/serjbibox/SF36a.3/pkg/models"
	"github.com/serjbibox/SF36a.3/pkg/storage/postgresql"
)

func TestHashPostgres_GetByLink(t *testing.T) {
	pwd := os.Getenv("DbPass")
	connString := "postgres://serj1:" + pwd + "@0.0.0.0:5438/gonews1?sslmode=disable"
	db, err := postgresql.New(connString)
	if err != nil {
		t.Fatal(err)
	}
	s := newHashPostgres(context.Background(), db)
	type args struct {
		l string
	}
	tests := []struct {
		name    string
		s       Hash
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test 1",
			s:    s,
			args: args{l: "https://habr.com/ru/rss/best/daily/?fl=ru"},
			want: "a80dcaf55478ae250d53ec6ac01eb562",
		},
		{
			name:    "test 2",
			s:       s,
			args:    args{l: "failLink"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetByLink(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPostgres.GetByLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HashPostgres.GetByLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPostgres_Create(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	pwd := os.Getenv("DbPass")
	connString := "postgres://serj1:" + pwd + "@0.0.0.0:5438/gonews1?sslmode=disable"
	db, err := postgresql.New(connString)
	if err != nil {
		t.Fatal(err)
	}
	s := newHashPostgres(context.Background(), db)
	type args struct {
		h models.Hash
	}
	tests := []struct {
		name    string
		s       Hash
		args    args
		wantErr bool
	}{
		{
			name: "test 1",
			s:    s,
			args: args{
				h: models.Hash{
					Link:     strconv.Itoa(rand.Intn(1_000_000_000)),
					NewsHash: strconv.Itoa(rand.Intn(1_000_000_000)),
				},
			},
		},
		{
			name: "test 2",
			s:    s,
			args: args{
				h: models.Hash{
					Link:     "https://blog.jetbrains.com/go/feed/",
					NewsHash: "9c2964fa7f5d51cc7ff462fba140624d",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Create(tt.args.h); (err != nil) != tt.wantErr {
				t.Errorf("HashPostgres.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashPostgres_Update(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	pwd := os.Getenv("DbPass")
	connString := "postgres://serj1:" + pwd + "@0.0.0.0:5438/gonews1?sslmode=disable"
	db, err := postgresql.New(connString)
	if err != nil {
		t.Fatal(err)
	}
	s := newHashPostgres(context.Background(), db)
	type args struct {
		h models.Hash
	}
	tests := []struct {
		name    string
		s       Hash
		args    args
		wantErr bool
	}{
		{
			name: "test 1",
			s:    s,
			args: args{
				h: models.Hash{
					Link:     "https://blog.jetbrains.com/go/feed/",
					NewsHash: "newhash",
				},
			},
			wantErr: false,
		},
		{
			name: "test 2",
			s:    s,
			args: args{
				h: models.Hash{
					Link:     "?",
					NewsHash: "newhash",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Update(tt.args.h); (err != nil) != tt.wantErr {
				t.Errorf("HashPostgres.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

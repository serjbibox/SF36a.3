package storage

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/serjbibox/SF36a.3/pkg/models"
	"github.com/serjbibox/SF36a.3/pkg/storage/postgresql"
)

var testDb *pgxpool.Pool

var testConnString = "postgres://news_service:qwerty@db_postgres:5432/testdb?sslmode=disable"

//var testConnString = "postgres://serj:123456@192.168.0.109:5432/gonews?sslmode=disable"

func TestMain(m *testing.M) {
	// Write code here to run before tests
	var err error
	testDb, err = postgresql.New(testConnString)
	if err != nil {
		log.Fatal(err)
	}
	// Run tests
	exitVal := m.Run()

	// Write code here to run after tests

	// Exit with exit value from tests
	os.Exit(exitVal)
}

func TestHashPostgres_GetByLink(t *testing.T) {
	s := newHashPostgres(context.Background(), testDb)
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
			args: args{l: "test link 1"},
			want: "testhash1",
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
	s := newHashPostgres(context.Background(), testDb)
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
					Link:     "test link 1",
					NewsHash: "testhash1",
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
	s := newHashPostgres(context.Background(), testDb)
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
					Link:     "test link 2",
					NewsHash: "testhash2",
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
					NewsHash: "testhash2",
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

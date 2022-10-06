package rss

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/serjbibox/SF36a.3/pkg/models"
	"github.com/serjbibox/SF36a.3/pkg/storage"
	"github.com/serjbibox/SF36a.3/pkg/storage/postgresql"
)

var testDb *pgxpool.Pool
var testConnString = "postgres://serj:123456@192.168.0.109:5432/gonews?sslmode=disable"

//var testConnString = "postgres://news_service:qwerty@db_postgres:5432/testdb?sslmode=disable"

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

func TestRss_parse(t *testing.T) {
	s, err := storage.NewStoragePostgres(context.Background(), testDb)
	if err != nil {
		t.Fatal(err)
	}
	r, err := New(s, "https://habr.com/ru/rss/best/daily/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	want := "storage is nil"
	_, err = New(nil, "https://habr.com/ru/rss/best/daily/?fl=ru")
	if err.Error() != want {
		t.Errorf("код неверен: получили %v, а хотели %v", err, "storage is nil")
	}
	tests := []struct {
		name    string
		r       *Rss
		want    []models.Post
		wantErr bool
	}{
		{
			name: "test 1",
			r:    r,
		},
		{
			name: "test 1",
			r: &Rss{
				Storage: s,
				Link:    "https://failLink",
				Hash: models.Hash{
					Link: "https://failLink",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			XMLdata, err := readRssBody(tt.r.Link)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rss.parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := tt.r.parse(XMLdata)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rss.parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (len(got) == 0) != tt.wantErr {
				t.Fatal("данные не раскодированы")
			}
			if len(got) != 0 {
				t.Logf("получено %d новостей\n%+v", len(got), got[0])
			}

		})
	}
}

func TestRss_ParseRss(t *testing.T) {
	s, err := storage.NewStoragePostgres(context.Background(), testDb)
	if err != nil {
		t.Fatal(err)
	}
	r, err := New(s, "https://habr.com/ru/rss/best/daily/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	news := make(chan []models.Post)
	errs := make(chan error)
	type args struct {
		period int
		news   chan []models.Post
		errs   chan error
	}
	tests := []struct {
		name string
		r    *Rss
		args args
	}{
		{
			name: "test 1",
			r:    r,
			args: args{
				period: 5,
				news:   news,
				errs:   errs,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.r.ParseRss(tt.args.period, tt.args.news, tt.args.errs)
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				for n := range tt.args.news {
					if len(n) != 0 {
						fmt.Println(n[0])
					}
					close(news)
					close(errs)
				}
				wg.Done()
			}()
			go func() {
				for n := range tt.args.errs {
					t.Log(n)
					close(news)
					close(errs)
				}
				wg.Done()
			}()
			wg.Wait()
		})
	}
}

func TestRss_hashInit(t *testing.T) {
	s, err := storage.NewStoragePostgres(context.Background(), testDb)
	if err != nil {
		t.Fatal(err)
	}
	r, err := New(s, "https://habr.com/ru/rss/best/daily/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	XMLdata, err := readRssBody(r.Link)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		r       *Rss
		args    args
		wantErr bool
	}{
		{
			name: "test 1",
			r:    r,
			args: args{
				data: XMLdata,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.hashInit(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Rss.hashInit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

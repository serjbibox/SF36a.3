package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/serjbibox/SF36a.3/pkg/models"
	"github.com/serjbibox/SF36a.3/pkg/storage"
	"github.com/serjbibox/SF36a.3/pkg/storage/memdb"
)

func TestHandler_getNews(t *testing.T) {
	db, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}
	s, err := storage.NewStorageMemDb(db)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/news/40", nil)
	rr := httptest.NewRecorder()
	handlers, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	handlers.InitRoutes().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	b, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	var data []models.Post
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/news/notnumber", nil)
	rr = httptest.NewRecorder()
	handlers.InitRoutes().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusInternalServerError) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusInternalServerError)
	}

	req = httptest.NewRequest(http.MethodGet, "/news/-6", nil)
	rr = httptest.NewRecorder()
	handlers.InitRoutes().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusInternalServerError) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusInternalServerError)
	}

	want := "storage is nil"
	_, err = New(nil)
	if err.Error() != want {
		t.Errorf("код неверен: получили %v, а хотели %v", err, "db is nil")
	}
}

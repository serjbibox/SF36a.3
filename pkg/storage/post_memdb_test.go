package storage

import (
	"reflect"
	"testing"

	"github.com/serjbibox/SF36a.3/pkg/models"
	"github.com/serjbibox/SF36a.3/pkg/storage/memdb"
)

var posts = []models.Post{
	{
		ID:      1,
		Title:   "Effective Go",
		Content: "Go is a new language. Although it borrows ideas from existing languages, it has unusual properties that make effective Go programs different in character from programs written in its relatives. A straightforward translation of a C++ or Java program into Go is unlikely to produce a satisfactory result—Java programs are written in Java, not Go. On the other hand, thinking about the problem from a Go perspective could produce a successful but quite different program. In other words, to write Go well, it's important to understand its properties and idioms. It's also important to know the established conventions for programming in Go, such as naming, formatting, program construction, and so on, so that programs you write will be easy for other Go programmers to understand.",
	},
	{
		ID:      2,
		Title:   "The Go Memory Model 1",
		Content: "The Go memory model specifies the conditions under which reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in a different goroutine.",
	},
	{
		ID:      3,
		Title:   "The Go Memory Model 2",
		Content: "The Go memory model specifies the conditions under which reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in a different goroutine.",
	},
	{
		ID:      4,
		Title:   "The Go Memory Model 3",
		Content: "The Go memory model specifies the conditions under which reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in a different goroutine.",
	},
}

func TestPostMemdb_GetByQuantity(t *testing.T) {
	db, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}
	s := newPostMemDb(db)
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		s       Post
		args    args
		want    []models.Post
		wantErr bool
	}{
		{
			name: "test 1",
			s:    s,
			args: args{n: 1},
			want: []models.Post{{
				ID:      4,
				Title:   "The Go Memory Model 3",
				Content: "The Go memory model specifies the conditions under which reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in a different goroutine.",
			}},
		},
		{
			name:    "test 2",
			s:       s,
			args:    args{n: -6},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test 3",
			s:       s,
			args:    args{n: 40},
			want:    posts,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetByQuantity(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostMemdb.GetByQuantity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostMemdb.GetByQuantity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostMemdb_Create(t *testing.T) {
	db, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}
	s := newPostMemDb(db)
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
			s:    s,
			args: args{
				[]models.Post{
					{
						ID:      4,
						Title:   "The Go Memory Model 3",
						Content: "The Go memory model specifies the conditions under which reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in a different goroutine.",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Create(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("PostMemdb.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

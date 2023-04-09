package ngamux

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang-must/must"
)

func TestSetPrefix(t *testing.T) {
	f := &fileServer{}

	f.Prefix("/something/")
	must.Equal(t, "/something/", f.prefix)

	f.Prefix("/templates/")
	must.Equal(t, "/templates/", f.prefix)
}

func TestServeDir(t *testing.T) {

	_file, err := os.CreateTemp("", "test.html")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(_file.Name())

	_, err = _file.Write([]byte("<html><body><h1>Test Page</h1></body></html>"))
	if err != nil {
		t.Fatal(err)
	}

	// Membuat http recorder untuk merekam responsenya
	rr := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/static/test.html", nil)
	if err != nil {
		t.Fatal(err)
	}

	f := &fileServer{
		w: rr,
		r: req,
	}

	f.Dir("/static")

	// Memeriksa tipe konten
	if ct := rr.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Errorf("handler returned wrong content type: got %v want text/plain; charset=utf-8", ct)
	}

	// Memeriksa isi konten
	expectedBody := "<html><body><h1>Test Page</h1></body></html>"
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

package ngamux

import "net/http"

type (
	fileServer struct {
		w      http.ResponseWriter
		r      *http.Request
		prefix string
	}
)

func FileServer(w http.ResponseWriter, r *http.Request) *fileServer {
	return &fileServer{w: w, r: r}
}

func (f *fileServer) Prefix(prefix string) *fileServer {
	f.prefix = prefix
	return f
}

func (f *fileServer) Dir(dir string) error {

	if f.prefix == "" {
		f.prefix = "/static/"
	}

	_fileServer := http.StripPrefix(f.prefix, http.FileServer(http.Dir(dir)))

	_fileServer.ServeHTTP(f.w, f.r)
	return nil
}

package ngamux

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-must/must"
)

func TestNewNgamux(t *testing.T) {
	must := must.New(t)
	result := New(
		WithLogLevel(LogLevelQuiet),
	)
	expected := &Ngamux{
		routes:            routeMap{},
		routesParam:       routeMap{},
		config:            NewConfig(),
		regexpParamFinded: paramsFinder,
	}

	must.Equal(expected.routes, result.routes)
	must.Equal(expected.routesParam, result.routesParam)
	must.Equal(expected.config.RemoveTrailingSlash, result.config.RemoveTrailingSlash)
	must.Equal(expected.regexpParamFinded, result.regexpParamFinded)
}

func TestUse(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	middleware := func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) error {
			return nil
		}
	}
	mux.Use(middleware)
	mux.Use(middleware)
	mux.Use(middleware)

	result := len(mux.middlewares)
	expected := 3

	must.Equal(expected, result)
}

func TestGet(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) error {
		return Res(rw).String("ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestPost(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	mux.Post("/", func(rw http.ResponseWriter, r *http.Request) error {
		return Res(rw).String("ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestPut(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	mux.Put("/", func(rw http.ResponseWriter, r *http.Request) error {
		return Res(rw).String("ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestPatch(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	mux.Patch("/", func(rw http.ResponseWriter, r *http.Request) error {
		return Res(rw).String("ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestDelete(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	mux.Delete("/", func(rw http.ResponseWriter, r *http.Request) error {
		return Res(rw).String("ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestAll(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	mux.All("/", func(rw http.ResponseWriter, r *http.Request) error {
		return Res(rw).String("ok")
	})

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(method, "/", nil)
		mux.ServeHTTP(rec, req)

		result := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expected := "ok"
		must.Equal(expected, result)
	}
}

func TestErrorResponse(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	mux.Get("/error-method", func(rw http.ResponseWriter, r *http.Request) error {
		return errors.New("something bad")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/error-method", nil)

	mux.ServeHTTP(rec, req)

	result := rec.Result()
	resultBody := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expectedBody := "something bad"
	must.Equal(expectedBody, resultBody)

	resultStatus := result.StatusCode
	expectedStatus := 500
	must.Equal(expectedStatus, resultStatus)
}

func TestWith(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	mux = mux.With(func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) error {
			return next(rw, r)
		}
	})
	must.NotNil(mux)
	must.NotNil(mux.parent)
}

func BenchmarkNgamux(b *testing.B) {
	h1 := func(w http.ResponseWriter, r *http.Request) error { return nil }
	h2 := func(w http.ResponseWriter, r *http.Request) error { return nil }
	h3 := func(w http.ResponseWriter, r *http.Request) error { return nil }
	h4 := func(w http.ResponseWriter, r *http.Request) error { return nil }
	h5 := func(w http.ResponseWriter, r *http.Request) error { return nil }
	h6 := func(w http.ResponseWriter, r *http.Request) error { return nil }

	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	mux.Get("/", h1)
	mux.Get("/hi", h2)
	mux.Get("/sup/:id/and/:this", h3)

	mux1 := mux.Group("/sharing/:x/:hash")
	mux1.Get("/", h4)          // subrouter-1
	mux1.Get("/{network}", h5) // subrouter-1
	mux1.Get("/twitter", h5)

	mux2 := mux.Group("/direct")
	mux2.Get("/", h6) // subrouter-2
	mux2.Get("/download", h6)

	routes := []string{
		"/",
		"/hi",
		"/sup/123/and/this",
		"/sup/123/foo/this",
		"/sharing/z/aBc",                 // subrouter-1
		"/sharing/z/aBc/twitter",         // subrouter-1
		"/sharing/z/aBc/direct",          // subrouter-2
		"/sharing/z/aBc/direct/download", // subrouter-2
	}

	for _, path := range routes {
		b.Run("route:"+path, func(b *testing.B) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", path, nil)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				mux.ServeHTTP(w, r)
			}
		})
	}
}

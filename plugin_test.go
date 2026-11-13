package markdown_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	markdown "github.com/juliens/plugin-markdown-demo"
)

func newHandler(t *testing.T, cfg *markdown.Config, status int, body string) http.Handler {
	t.Helper()

	next := http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Set("Content-Type", "text/markdown")
		rw.Header().Set("X-Backend", "files")
		rw.WriteHeader(status)
		_, _ = rw.Write([]byte(body))
	})

	handler, err := markdown.New(context.Background(), next, cfg, "markdown")
	if err != nil {
		t.Fatal(err)
	}

	return handler
}

func TestMarkdown_render(t *testing.T) {
	sample, err := os.ReadFile("testdata/sample.md")
	if err != nil {
		t.Fatal(err)
	}

	handler := newHandler(t, markdown.CreateConfig(), http.StatusOK, string(sample))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/docs/README.md", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	if got := rec.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Errorf("Content-Type: got %q", got)
	}

	if got := rec.Header().Get("X-Backend"); got != "files" {
		t.Errorf("X-Backend: got %q, want the header of the backend", got)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"<title>README.md</title>",
		`<h1 id="yaegi-v2-demo">Yaegi v2 demo</h1>`,
		"<table>",
		`<input checked="" disabled="" type="checkbox"`,
		// A highlighted Go keyword.
		`<span style="color:#cf222e">func</span>`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("the page does not contain %q", want)
		}
	}
}

func TestMarkdown_otherPaths(t *testing.T) {
	handler := newHandler(t, markdown.CreateConfig(), http.StatusOK, "# Title")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/index.html", nil))

	if got := rec.Body.String(); got != "# Title" {
		t.Errorf("body: got %q, want the response of the backend", got)
	}
}

func TestMarkdown_errorStatus(t *testing.T) {
	handler := newHandler(t, markdown.CreateConfig(), http.StatusNotFound, "not found")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/missing.md", nil))

	if rec.Code != http.StatusNotFound || rec.Body.String() != "not found" {
		t.Errorf("got %d %q, want the response of the backend", rec.Code, rec.Body.String())
	}
}

func TestMarkdown_cacheHeaders(t *testing.T) {
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("If-None-Match") != "" || req.Header.Get("If-Modified-Since") != "" {
			rw.WriteHeader(http.StatusNotModified)
			return
		}

		rw.Header().Set("ETag", `"file-etag"`)
		rw.Header().Set("Last-Modified", "Mon, 02 Jan 2006 15:04:05 GMT")
		_, _ = rw.Write([]byte("# Title"))
	})

	handler, err := markdown.New(context.Background(), next, markdown.CreateConfig(), "markdown")
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/a.md", nil)
	req.Header.Set("If-None-Match", `"file-etag"`)
	req.Header.Set("If-Modified-Since", "Mon, 02 Jan 2006 15:04:05 GMT")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d: the page must not use the cache of the Markdown file", rec.Code, http.StatusOK)
	}

	for _, name := range []string{"ETag", "Last-Modified"} {
		if v := rec.Header().Get(name); v != "" {
			t.Errorf("header %s: got %q, want it removed", name, v)
		}
	}
}

func TestNew_unknownStyle(t *testing.T) {
	cfg := markdown.CreateConfig()
	cfg.Style = "does-not-exist"

	_, err := markdown.New(context.Background(), http.NotFoundHandler(), cfg, "markdown")
	if err == nil {
		t.Fatal("expected an error")
	}
}

// Package markdown is a Traefik middleware plugin, run by Yaegi v2.
// It renders the Markdown files served by the backend as HTML pages,
// with syntax highlighting of the code blocks.
//
// The appengine tag selects the files of goldmark that do not use unsafe:
// Yaegi v2 refuses unsafe in a plugin that does not set useUnsafe.
//
// yaegi:tags appengine
package markdown

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// Config is the configuration of the plugin.
type Config struct {
	// Style is the name of the chroma style for the code blocks.
	Style string `json:"style,omitempty"`
	// LineNumbers adds line numbers to the code blocks.
	LineNumbers bool `json:"lineNumbers,omitempty"`
}

// CreateConfig creates the default configuration of the plugin.
func CreateConfig() *Config {
	return &Config{Style: "github"}
}

// Markdown is the middleware.
type Markdown struct {
	next     http.Handler
	name     string
	renderer goldmark.Markdown
}

// New creates a new instance of the middleware.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if styles.Get(config.Style) == styles.Fallback && config.Style != "swapoff" {
		return nil, fmt.Errorf("unknown style: %s", config.Style)
	}

	renderer := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle(config.Style),
				highlighting.WithFormatOptions(chromahtml.WithLineNumbers(config.LineNumbers)),
			),
		),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)

	return &Markdown{next: next, name: name, renderer: renderer}, nil
}

func (m *Markdown) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if !strings.EqualFold(path.Ext(req.URL.Path), ".md") {
		m.next.ServeHTTP(rw, req)
		return
	}

	// The page is built from the whole file: the backend must not answer
	// 304 Not Modified or a part of the file.
	req = req.Clone(req.Context())
	for _, name := range conditionalHeaders {
		req.Header.Del(name)
	}

	rec := &recorder{header: http.Header{}, status: http.StatusOK}
	m.next.ServeHTTP(rec, req)

	if rec.status != http.StatusOK {
		rec.copyTo(rw)
		return
	}

	var body bytes.Buffer
	if err := m.renderer.Convert(rec.body.Bytes(), &body); err != nil {
		http.Error(rw, "cannot render the Markdown file", http.StatusInternalServerError)
		return
	}

	var page bytes.Buffer
	err := pageTemplate.Execute(&page, struct {
		Title string
		Body  template.HTML
	}{
		Title: path.Base(req.URL.Path),
		Body:  template.HTML(body.String()), //nolint:gosec // The HTML comes from goldmark, which escapes raw HTML by default.
	})
	if err != nil {
		http.Error(rw, "cannot render the page", http.StatusInternalServerError)
		return
	}

	for key, values := range rec.header {
		rw.Header()[key] = values
	}

	// These headers describe the Markdown file, not the HTML page.
	for _, name := range fileHeaders {
		rw.Header().Del(name)
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.Header().Set("Content-Length", strconv.Itoa(page.Len()))
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(page.Bytes())
}

var (
	conditionalHeaders = []string{"If-Modified-Since", "If-None-Match", "If-Range", "Range", "Accept-Encoding"}
	fileHeaders        = []string{"ETag", "Last-Modified", "Accept-Ranges", "Content-Range", "Content-Encoding", "Content-Length"}
)

// recorder keeps the response of the backend.
type recorder struct {
	header http.Header
	status int
	body   bytes.Buffer
}

func (r *recorder) Header() http.Header {
	return r.header
}

func (r *recorder) WriteHeader(status int) {
	r.status = status
}

func (r *recorder) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func (r *recorder) copyTo(rw http.ResponseWriter) {
	for key, values := range r.header {
		rw.Header()[key] = values
	}

	rw.WriteHeader(r.status)
	_, _ = rw.Write(r.body.Bytes())
}

var pageTemplate = template.Must(template.New("page").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{ .Title }}</title>
<style>
html { background: #fff; color-scheme: light; }
body { max-width: 860px; margin: 2rem auto; padding: 0 1rem; font-family: system-ui, sans-serif; line-height: 1.6; color: #1f2328; }
pre { padding: 1rem; overflow-x: auto; border-radius: 6px; }
code { font-family: ui-monospace, monospace; font-size: 0.9em; }
table { border-collapse: collapse; }
th, td { border: 1px solid #d0d7de; padding: 0.3rem 0.8rem; }
footer { margin-top: 3rem; font-size: 0.8em; color: #656d76; }
</style>
</head>
<body>
{{ .Body }}
<footer>Rendered by a Traefik plugin run by Yaegi v2.</footer>
</body>
</html>
`))

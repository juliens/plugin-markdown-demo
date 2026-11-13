# How the plugin works

[Back to the index](index.md)

The middleware only looks at the paths that end with `.md`.

1. It calls the backend and keeps its response in memory.
2. If the status is not `200`, it sends the response as it is.
3. Otherwise, goldmark converts the Markdown to HTML, and chroma highlights the code blocks.
4. The HTML goes into a small page template, with the `text/html` content type.

## Creation of the renderer

The renderer is created once, in `New`, when Traefik loads the configuration:

```go
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
```

## For each request

```go
var body bytes.Buffer
if err := m.renderer.Convert(rec.body.Bytes(), &body); err != nil {
	http.Error(rw, "cannot render the Markdown file", http.StatusInternalServerError)
	return
}
```

## A missing page

The backend answers `404` for [a page that does not exist](missing.md): the plugin does not change this response.

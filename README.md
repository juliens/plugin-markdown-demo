# Markdown Demo

A demo of a Traefik plugin that runs large libraries with [Yaegi v2](https://github.com/traefik/yaegi/tree/v2).

Yaegi v2 is experimental, and Traefik supports it only on a development branch for now.

The middleware renders the Markdown files served by the backend as HTML pages:

- GitHub Flavored Markdown (tables, task lists, strikethrough, links);
- syntax highlighting of the code blocks, for more than 250 languages;
- the other responses are not changed.

With a file server behind Traefik, a folder of Markdown files becomes a small documentation site.

## Why it is a good demo of Yaegi v2

The plugin is small, but it runs a lot of code:

- [goldmark](https://github.com/yuin/goldmark): the Markdown parser and renderer;
- [chroma](https://github.com/alecthomas/chroma): the syntax highlighter, with its lexers;
- [regexp2](https://github.com/dlclark/regexp2): the regular expression engine of chroma.

The vendor directory has about 45,000 lines of Go, all run by the interpreter.

Yaegi v1 cannot load this plugin: it fails on the recursive types of regexp2.
Yaegi v2 loads it and gives the same HTML as compiled Go.

Measures on a laptop (Apple M-series, Go 1.27.3), for the page in [`testdata/sample.md`](testdata/sample.md):

| | Compiled Go | Yaegi v2 |
|---|---|---|
| Load of the plugin | | 0.4 s |
| First request | 2.3 ms | 35 ms |
| Next requests | 0.5 ms | 15 ms |

## No unsafe

Yaegi v2 refuses the `unsafe` package in a plugin that does not set `useUnsafe`. Two dependencies use it:

- goldmark uses `unsafe` for fast string conversions, but has a safe version of these files for the `appengine` build tag.
  The plugin selects it with a comment that Go ignores and that Yaegi reads:

  ```go
  // yaegi:tags appengine
  package markdown
  ```

- chroma uses `github.com/dlclark/regexp2/v2` from v2.25.0, which uses `unsafe` without a safe version.
  The plugin uses chroma v2.24.1, the last version with regexp2 v1.

## Ask for Yaegi v2

The manifest [`.traefik.yml`](.traefik.yml) has `runtimeVersion: v2`:

```yaml
displayName: Markdown Demo
type: middleware
runtimeVersion: v2
import: github.com/juliens/plugin-markdown-demo
basePkg: markdown
```

Without this line, Traefik runs the plugin with Yaegi v1, and the plugin does not load.

## Run the demo

The [`demo`](demo) directory runs Traefik with this plugin from the Plugin Catalog, in front of an nginx that serves a few Markdown files.

```bash
cd demo && docker compose up --build
```

Then open http://localhost:8000.

- `/` renders the pages with the `github` style.
- `/dark/` renders the same pages with the `monokai` style and line numbers.
- Edit a file in `demo/docs`, then reload the page.
- The Traefik dashboard is on http://localhost:8081.

Traefik v3.8.0-rc.3 is the first version with Yaegi v2 support.
[`demo/Dockerfile.traefik`](demo/Dockerfile.traefik) builds it from its tag, with the web UI, as `make binary` does (about 4 minutes the first time).

## Configuration

Static configuration:

```yaml
experimental:
  plugins:
    markdown:
      moduleName: github.com/juliens/plugin-markdown-demo
      version: v0.1.0
```

Dynamic configuration:

```yaml
http:
  routers:
    docs:
      rule: Host(`docs.localhost`)
      service: files
      middlewares:
        - markdown

  middlewares:
    markdown:
      plugin:
        markdown:
          style: github     # a chroma style: github, monokai, dracula...
          lineNumbers: false

  services:
    files:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:8080
```

For example, serve a folder with `python3 -m http.server 8080`, then open `http://docs.localhost/README.md`.

## CI

[`.github/workflows/main.yml`](.github/workflows/main.yml) has two jobs:

- `Lint and test`: `golangci-lint`, `go test` (compiled Go), and a check that `go.mod` and `vendor` are up to date.
- `Run with Yaegi v2`: the checks of the Traefik Plugin Catalog, a load of the plugin with Yaegi v2 (with the options of Traefik), and the rendering of `testdata/sample.md`.

The second job runs [`.github/plugincheck`](.github/plugincheck), a separate Go module. Run it locally:

```bash
cd .github/plugincheck && go run . -dir ../.. -request /docs/sample.md -body ../../testdata/sample.md -contains '<table>'
```

The Go version must match the AOT copies of the `yaegi/v2` version in `.github/plugincheck/go.mod` (today Go 1.27.3).

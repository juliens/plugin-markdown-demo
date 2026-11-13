# Yaegi v2

[Back to the index](index.md)

Traefik runs its Go plugins with Yaegi, a Go interpreter: no compiler, no restart.

Yaegi v2 is a new engine. It uses the official Go tools (`go/parser`, `go/types`, `x/tools/go/ssa`) to read and check the code, and runs it with a bytecode VM.

## What it changes for this plugin

This plugin uses three libraries, about 45,000 lines of Go in its `vendor` directory:

| Library | Role |
|---|---|
| goldmark | Markdown parser and HTML renderer |
| chroma | syntax highlighter, with more than 250 lexers |
| regexp2 | regular expression engine used by chroma |

- **Yaegi v1** cannot load it: it fails on the recursive types of regexp2.
- **Yaegi v2** loads it, and the HTML is the same as with compiled Go.

## Ask for Yaegi v2

One line in `.traefik.yml`:

```yaml
runtimeVersion: v2
```

## No unsafe

A plugin cannot use the `unsafe` package unless the user allows it. goldmark has a safe version of its files for the `appengine` build tag, and the plugin selects it with a comment that Go ignores:

```go
// yaegi:tags appengine
package markdown
```

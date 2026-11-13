# Markdown Demo

This site is a folder of Markdown files, served by a plain nginx.
nginx knows nothing about Markdown: a **Traefik plugin run by Yaegi v2** renders each page on the fly.

## Pages

- [Yaegi v2](yaegi-v2.md): what changes for plugins.
- [Languages](languages.md): syntax highlighting for many languages.
- [How the plugin works](plugin.md): the code of the middleware.
- [Markdown features](features.md): tables, task lists, and more.

## Two styles

The same files are served on two routes, with two configurations of the same plugin:

| Route | Style | Line numbers |
|---|---|---|
| [`/`](/index.md) | `github` | no |
| [`/dark/`](/dark/index.md) | `monokai` | yes |

```yaml
middlewares:
  markdown-dark:
    plugin:
      markdown:
        style: monokai
        lineNumbers: true
```

## Try it

Edit a file in `demo/docs`, then reload the page: there is no build step.

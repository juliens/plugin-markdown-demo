# Markdown features

[Back to the index](index.md)

The plugin uses GitHub Flavored Markdown.

## Text

Some **bold**, some *italic*, some ~~strikethrough~~, and some `inline code`.

Links become clickable: https://github.com/traefik/yaegi

> A quote, on
> two lines.

## Task lists

- [x] Load the plugin with Yaegi v2
- [x] Render Markdown
- [x] Highlight code
- [ ] Your next plugin

## Tables

| Engine | Load | Request |
|:---|---:|---:|
| Compiled Go | | 0.5 ms |
| Yaegi v2 | 0.4 s | 15 ms |

## Lists

1. First
2. Second
   - a nested item
   - another one
3. Third

## Headings get an id

Each heading has an `id`, so a link can point to it: [go to Tables](#tables).

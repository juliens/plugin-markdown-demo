# Yaegi v2 demo

This page is **Markdown**, rendered by a Traefik plugin.

| Engine | Language |
|---|---|
| Yaegi v2 | Go |

- [x] GFM task list
- [ ] Another item

```go
package main

import "fmt"

func main() {
	for i := range 3 {
		fmt.Println("hello", i)
	}
}
```

```yaml
http:
  middlewares:
    markdown:
      plugin:
        markdown:
          style: github
```

```javascript
const add = (a, b) => a + b;
console.log(add(1, 2));
```

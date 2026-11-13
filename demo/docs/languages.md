# Languages

[Back to the index](index.md)

Each block below is highlighted by chroma, inside the plugin.

## Go

```go
package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	stars := map[string]int{"traefik": 50000, "yaegi": 7000}
	for _, name := range slices.Sorted(maps.Keys(stars)) {
		fmt.Printf("%s: %d\n", name, stars[name])
	}
}
```

## Rust

```rust
fn main() {
    let names = vec!["traefik", "yaegi"];
    for (i, name) in names.iter().enumerate() {
        println!("{i}: {name}");
    }
}
```

## Python

```python
from dataclasses import dataclass

@dataclass
class Route:
    rule: str
    service: str

print(Route(rule="PathPrefix(`/`)", service="docs"))
```

## TypeScript

```typescript
type Middleware = (req: Request, next: () => Response) => Response;

const addHeader: Middleware = (req, next) => {
  const res = next();
  res.headers.set("X-Demo", "yaegi-v2");
  return res;
};
```

## SQL

```sql
SELECT name, COUNT(*) AS requests
FROM access_log
WHERE status >= 500
GROUP BY name
ORDER BY requests DESC;
```

## Shell

```bash
docker compose up --build
curl -s http://localhost:8000/languages.md | head
```

## Dockerfile

```dockerfile
FROM nginx:1.29-alpine
COPY docs /usr/share/nginx/html
```

## TOML

```toml
[experimental.localPlugins.markdown]
  moduleName = "github.com/juliens/plugin-markdown-demo"
```

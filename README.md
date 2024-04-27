# toukibo-parser-worker

## Development
### Commands

```
make dev     # run dev server
make build   # build Go Wasm binary
make deploy # deploy worker
```

### Call API

prod
```
curl -X POST -H "Content-Type: application/pdf" --data-binary "@sample.pdf" https://go-worker.a2sin2a2ko1115.workers.dev/parse
```

dev

```
curl -X POST -H "Content-Type: application/pdf" --data-binary "@sample.pdf"  http://localhost:8787/parse
```

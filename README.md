# toukibo-parser-worker
[商業登記簿パーサー](https://github.com/tychy/toukibo-parser)を組み込んだAPI

商業登記簿PDFをPostすると、解析結果を返します。

[デモページ](https://toukibo-parser-demo.tychy.jp/?index)で試すことができます。

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

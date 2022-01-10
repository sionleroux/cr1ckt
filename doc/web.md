get WASM exec:

```bash
cp $(go env GOROOT)/misc/wasm/wasm_exec.js .
```

build:

```bash
GOOS=js GOARCH=wasm go build -ldflags "-w -s" -o cr1ckt.wasm
```

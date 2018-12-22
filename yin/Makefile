NODEWASM="`go env GOROOT`/misc/wasm/go_js_wasm_exec"

all: build test

build:
	cp "`go env GOROOT`/misc/wasm/wasm_exec.js" ./static
	GOOS=js GOARCH=wasm go build -o yin.wasm 

test:
	GOOS=js GOARCH=wasm go test -v -exec=$(NODEWASM) ./...

run:
	cd ../yang && make build 
	../yang/yang
	google-chrome http://localhost:3000
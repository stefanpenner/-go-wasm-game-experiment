
install:
	# wget https://raw.githubusercontent.com/tinygo-org/tinygo/52983794d702af7a00833ae12b0d2e7175e46017/targets/wasm_exec.js
	# wget https://raw.githubusercontent.com/golang/go/refs/heads/master/lib/wasm/wasm_exec.js
main:
	# cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
	# GOOS=js GOARCH=wasm go build o main.wasm main.go
	GOOS=js GOARCH=wasm go build -ldflags="-s -w=0" -o main.wasm

	python3 -m http.server &
	open http://localhost:8000
	fg

test_native:
	GOOS=darwin GOARCH=arm64 go test ./rect.go ./rect_test.go

dlv_test:
	GOOS=darwin GOARCH=arm64 dlv test ./player.go ./player_test.go ./rect.go ./rect_test.go

test_node:
	GOOS=js GOARCH=wasm go test -c -o test.wasm
	node test.mjs

test: test_node test_native

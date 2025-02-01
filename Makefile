
install:
	#https://tinygo.org/getting-started/install/macos/
	brew tap tinygo-org/tools
	brew install tinygo
	# wget https://raw.githubusercontent.com/tinygo-org/tinygo/52983794d702af7a00833ae12b0d2e7175e46017/targets/wasm_exec.js
	# wget https://raw.githubusercontent.com/golang/go/refs/heads/master/lib/wasm/wasm_exec.js
main:
	# cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
	# GOOS=js GOARCH=wasm go build o main.wasm main.go
	GOOS=js GOARCH=wasm go build -ldflags="-s -w=0" -o main.wasm

	python3 -m http.server &
	open http://localhost:8000
	fg

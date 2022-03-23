.PHONY: build

build:
	@echo "[*] Building helper app 'htb'"
	@mkdir -p build
	@go build -o build/ ./...
	@echo "[OK] Done building"

install:
	@echo "[*] Installing it to gopath '$$GOPATH/bin'"
	@go install ./...
	@echo "[OK] Done installing it"
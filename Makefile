init:
	$(MAKE) clone-tree-sitter
	$(MAKE) build-parser

clone-tree-sitter:
	[ ! -d "tree-sitter-c3" ] && git clone git@github.com:cbuttner/tree-sitter-c3.git assets/tree-sitter-c3 || true

build-parser:
	cd assets/tree-sitter-c3 && tree-sitter generate
	rm -rf server/lsp/cst/tree_sitter
	rm -f server/lsp/cst/parser.c
	cp -r assets/tree-sitter-c3/src/tree_sitter server/lsp/cst
	cp assets/tree-sitter-c3/src/parser.c server/lsp/cst/parser.c
	cp assets/tree-sitter-c3/src/scanner.c server/lsp/cst/scanner.c

index-c3-std:
	cd server/stdlib_indexer && go run main.go blurp.go
	cp server/stdlib_indexer/stdlib/*.go server/lsp/language/stdlib

build:
	go build -C server -o bin/c3-lsp

build-dev:
	go build -C server -gcflags="all=-N -l" -o bin/c3-lsp

build-all:
# Build darwin-amd64
	echo "Building darwin-amd64"
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -C server -o bin/c3-lsp
	cd server/bin && zip ./darwin-amd64-c3lsp.zip c3-lsp
	echo "darwin-amd64 built"

# Build linux
	echo "Building linux-amd64"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-linux-musl-gcc" go build -C server -o bin/c3-lsp
	cd server/bin && zip ./linux-amd64-c3lsp.zip c3-lsp
	echo "linux-amd64 built"


#attach-process:
#	dlv attach --headless --listen=:2345 $(pgrep c3-lsp) ./server/c3-lsp --log

test:
	cd server && go test ./...


## VS Code extension
build-vscode:
	cd client/vscode && npm run vscode:prepublish
	cd client/vscode && vsce package
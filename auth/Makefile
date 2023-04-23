BIN=./bin
FILE=./bin/server
SRC=./cmd/server/server.go
.PHONY: all default prereqs
default: prereqs
prereqs: clean ensure_deps test fix_ordering lint fmt
all: prereqs build
clean:
	$(info ******** Cleaning the bin folder ********)
	@./scripts/clean.sh $(BIN)
test:
	$(info ******** Testing the app ********)
	@./scripts/test.sh
build:
	$(info ******** Building the app ********)
	@./scripts/build.sh $(SRC) $(FILE)
lint:
	$(info ******** Linting the app ********)
	golangci-lint run ./...
fmt:
	$(info ******** Formatting the app ********)
	gofmt -w -s -l .
fix_ordering:
	$(info ******** Fixing the struct ordering ********)
	fieldalignment -fix ./...
ensure_deps:
	go mod tidy

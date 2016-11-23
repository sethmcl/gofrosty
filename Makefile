BIN_NAME=frosty

.PHONY: clean install darwin test

build:
	mkdir -p build/darwin.amd64
	go build -v -o build/$(BIN_NAME) cmd/gofrosty/main.go
	GOOS=darwin GOARCH=amd64 go build -v -o build/darwin.amd64/$(BIN_NAME) cmd/gofrosty/main.go

clean:
	rm -rf build

install: clean build
	cp build/$(BIN_NAME) $(GOPATH)/bin/

test:
	cd test && go test -v

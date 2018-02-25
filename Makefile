BIN=bin
SRC=./src
APP=megaedx7-run

build:	$(BIN)/$(APP)

$(BIN)/$(APP): $(SRC)/$(APP)/main.go
	go build -o $@ $(SRC)/$(APP)

install:
	go install $(SRC)/$(APP)

.PHONY: build install

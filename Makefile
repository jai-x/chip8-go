SRC = $(wildcard *.go)

all: $(SRC)
	go build

test:
	go test -v

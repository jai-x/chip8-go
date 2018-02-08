SRC = $(wildcard *.go)

all: $(SRC)
	go build

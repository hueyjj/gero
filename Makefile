GOPATH=$(shell go env GOPATH)

all: phony

phony:
	echo "gero program"

build:
	go build -o bin/gero

clean:
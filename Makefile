.PHONY: build run
build:
	go build

run: build
	./record-updater

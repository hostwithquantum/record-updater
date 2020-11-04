.PHONY: build run test
build:
	goreleaser --snapshot --skip-publish --rm-dist

run: build
	./dist/record-updater_darwin_amd64/record-updater

test:
	$(info Nothing yet)

all: clean build

clean:
	rm -r dist/ || true

build:
	$$GOPATH/bin/goreleaser --config=.github/goreleaser.yml --snapshot

run:
	go run ./cmd/...

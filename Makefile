.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o infra/bin/query query/main.go

clean:
	rm -rf ./infra/bin

deploy: clean build
	sls deploy --verbose

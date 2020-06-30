.PHONY: build build_windows build_mac test

build:
	go build -mod=vendor -o ukfast -ldflags "-s -X 'main.VERSION=$$(git describe --tags)' -X 'main.BUILDDATE=$$(date +'%Y-%m-%dT%H:%M:%S')'"
	chmod +x ukfast

build_windows:
	GOOS=windows go build -mod=vendor -o ukfast.exe -ldflags "-s -X 'main.VERSION=$$(git describe --tags)' -X 'main.BUILDDATE=$$(date +'%Y-%m-%dT%H:%M:%S')'"

build_mac:
	GOOS=darwin go build -mod=vendor -o ukfast -ldflags "-s -X 'main.VERSION=$$(git describe --tags)' -X 'main.BUILDDATE=$$(date +'%Y-%m-%dT%H:%M:%S')'"
	chmod +x ukfast

test:
	go test -mod=vendor -v -cover ./...

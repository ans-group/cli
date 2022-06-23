.PHONY: build build_windows build_mac test

build:
	go build -o ans -ldflags "-s -X 'main.VERSION=$$(git describe --tags)' -X 'main.BUILDDATE=$$(date +'%Y-%m-%dT%H:%M:%S')'"

clean:
	go clean --modcache 
	rm ans

install: 
	sudo cp ans /usr/local/bin/    

build_windows:
	GOOS=windows go build -o ans.exe -ldflags "-s -X 'main.VERSION=$$(git describe --tags)' -X 'main.BUILDDATE=$$(date +'%Y-%m-%dT%H:%M:%S')'"

build_mac:
	GOOS=darwin go build -o ans -ldflags "-s -X 'main.VERSION=$$(git describe --tags)' -X 'main.BUILDDATE=$$(date +'%Y-%m-%dT%H:%M:%S')'"

test:
	go test -v -cover ./...


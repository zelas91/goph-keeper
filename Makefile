ARCH = amd64 arm64 arm 386
PLATFORMS = linux darwin windows
NAME_SERVER = server
NAME_CLIENT = goph-keeper-client

.SILENT:
.PHONY:
encrypt:
	if [ ! -f "secretKey.txt" ]; then \
		openssl rand -hex 8 > secretKey.txt; \
	fi

clean:
	if [ -d "build" ]; then \
    	rm -r build; \
    fi

server-build:encrypt
	go build -o build/server/$(NAME_SERVER) \
 -ldflags "-X main.buildCommit=$$(git rev-parse --short HEAD)\
     -X main.buildDate=$$(date +'%Y-%m-%d_%H:%M')\
     -X main.secretKey=$(shell cat secretKey.txt)"\
     cmd/server/*.go
client-build:
	go build -o build/client/$(NAME_CLIENT)-$(GOOS)-$(GOARCH) \
	-ldflags "-X main.buildCommit=$$(git rev-parse --short HEAD)\
    -X main.buildDate=$$(date +'%Y-%m-%d_%H:%M')" cmd/client/*.go

run-docker:
	make server-build
	docker-compose up

build-all:clean encrypt
	$(foreach GOOS,$(PLATFORMS),\
		$(foreach GOARCH,$(ARCH),\
			GOOS=$(GOOS) GOARCH=$(GOARCH) make client-build;))
			make server-build;

test-all:
	go generate ./...
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out

fmt:
	gofmt -s -w .
	goimports -w .
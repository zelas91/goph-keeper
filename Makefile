.PHONY: encrypt

encrypt:
	if [ ! -f "secretKey.txt" ]; then \
		openssl rand -hex 8 > secretKey.txt; \
	fi


ARCH = amd64 arm64 arm 386
PLATFORMS = linux darwin windows
NAME_SERVER = server
NAME_CLIENT = goph-keeper-client
SECRET_KEY := $(shell cat secretKey.txt)
clean:
	if [ -d "build" ]; then \
    	rm -r build; \
    fi

server-build:
	go build -o build/server/$(NAME_SERVER) cmd/server/*.go
client-build:
	go build -o build/client/$(NAME_CLIENT)-$(GOOS)-$(GOARCH) \
	-ldflags "-X main.buildCommit=$$(git rev-parse --short HEAD)\
    -X main.buildDate=$$(date +'%Y-%m-%d_%H:%M')\
    -X main.secretKey=$(SECRET_KEY)" cmd/client/*.go

run-docker:
	make server-build
	docker-compose up

build-all-client:
	$(foreach GOOS,$(PLATFORMS),\
		$(foreach GOARCH,$(ARCH),\
			GOOS=$(GOOS) GOARCH=$(GOARCH) make client-build;))

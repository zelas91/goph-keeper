ARCH = amd64 arm64 arm 386
PLATFORMS = linux darwin windows
NAME_SERVER = server
NAME_CLIENT = goph-keeper-client

clean:
	if [ -d "build" ]; then \
    	rm -r build; \
    fi

server-build:
	go build -o build/server/$(NAME_SERVER) cmd/server/*.go

run-docker: clean
	make server-build
	docker-compose up

build-all-client: clean
	$(foreach GOOS,$(PLATFORMS),\
		$(foreach GOARCH,$(ARCH),\
			GOOS=$(GOOS) GOARCH=$(GOARCH) make client-build;))
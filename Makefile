server-build:
	go build   -o build/server/server  cmd/server/*.go

run-docker:
	make server-build
	docker-compose up
FROM golang:1.20
LABEL authors="zelas"

ENV GOPATH=/app
WORKDIR /app
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.12.1/wait /wait

#RUN go mod download
RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash
RUN apt-get update
RUN apt-get install migrate
RUN chmod +x /wait
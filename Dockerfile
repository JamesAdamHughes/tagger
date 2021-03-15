FROM golang:1.16-alpine
RUN apk add build-base

COPY air.toml /etc/

ENV BASE_URL="http://localhost"
ENV REDIS_DOMAIN="redis"

RUN go get -u -v github.com/cosmtrek/air

WORKDIR /go/src/tagger
COPY go.mod .
COPY go.sum .

RUN go mod tidy

RUN go mod download
CMD /go/bin/air -c /etc/air.toml
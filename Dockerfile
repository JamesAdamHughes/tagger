FROM golang:1.16

COPY air.toml /etc/

RUN go get -u -v github.com/cosmtrek/air

WORKDIR /go/src/tagger
COPY go.mod .
COPY go.sum .

RUN go mod tidy

RUN go mod download
CMD /go/bin/air -c /etc/air.toml
FROM golang:latest

MAINTAINER "Ove Stavnås"

WORKDIR /go/src/app

COPY *.go ./

RUN go get && go build -o app

EXPOSE 8080

CMD ["./app"]

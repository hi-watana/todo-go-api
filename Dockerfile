FROM golang:latest

WORKDIR /app
ADD *.go go.mod /app/
RUN go get .
CMD go run .

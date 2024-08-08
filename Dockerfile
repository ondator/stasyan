FROM golang:1.22-alpine3.20 as build

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o stasyan

FROM alpine:3.20 as run

WORKDIR /

COPY --from=build stasyan stasyan

ENTRYPOINT ["/stasyan"]
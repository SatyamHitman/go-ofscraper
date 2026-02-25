FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates ffmpeg

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /build/bin/gofscraper ./cmd/gofscraper

FROM alpine:3.19
RUN apk add --no-cache ca-certificates ffmpeg tzdata
COPY --from=builder /build/bin/gofscraper /usr/local/bin/gofscraper

ENTRYPOINT ["gofscraper"]

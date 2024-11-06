FROM golang:1.23-alpine AS builder

COPY . /app

WORKDIR /app

# Add gcc
RUN apk add build-base tzdata

RUN go mod download && \
    go env -w GOFLAGS=-mod=mod && \
    go get . && \
    go build -v -o backend .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/backend ./backend

CMD [ "./backend" ]

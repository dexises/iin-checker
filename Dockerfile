FROM golang:1.20-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o iinservice internal/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/iinservice .

EXPOSE 8080

ENTRYPOINT ["./iinservice"]

FROM golang:1.23 AS builder 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o app .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .
COPY config.yaml .      
COPY .env .            


CMD ["./app"]

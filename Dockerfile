# syntax=docker/dockerfile:1

FROM golang:1.22.3

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /ews-backend ./cmd

CMD ["go", "run", "./cmd/main.go"]


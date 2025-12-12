FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN go build -o test-section-serve ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/test-section-serve .


EXPOSE 50052

CMD ["./test-section-serve"]
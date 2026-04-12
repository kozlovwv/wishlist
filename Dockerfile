FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/wishlist


FROM alpine:3.20

WORKDIR /service

COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./app"]
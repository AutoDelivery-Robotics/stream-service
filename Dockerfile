FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY shared/go shared/go
COPY backend/stream-service backend/stream-service
WORKDIR /app/backend/stream-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/backend/stream-service/main .
CMD ["./main"]

# Stage 1: Builder
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./main.go

# Stage 2: Final Image 
FROM alpine:latest
WORKDIR /app
COPY --from=builder /main .
COPY web ./web
EXPOSE 7540
CMD ["./main"]

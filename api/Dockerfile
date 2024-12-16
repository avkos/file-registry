FROM golang:1.23.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o api main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]

FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/api_frete main.go

FROM alpine:latest

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder /app/api_frete .


COPY --from=builder /app/migrations ./migrations

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./api_frete"]
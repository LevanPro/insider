FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/service/main.go

FROM alpine:latest

RUN apk update && apk upgrade

# Reduce image size
RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

WORKDIR /app

COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/app .


# Avoid running code as a root user
RUN adduser -D appuser
USER appuser

ENV HTTP_ADDR=:8080

EXPOSE 8080

CMD ["./app"]

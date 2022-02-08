FROM golang:1.17-alpine as builder
COPY . /app
WORKDIR /app
RUN go build -o ./telegram-article-bot ./cmd

FROM alpine:3.15.0
COPY --from=builder /app/telegram-article-bot .
ENTRYPOINT ["./telegram-article-bot"]
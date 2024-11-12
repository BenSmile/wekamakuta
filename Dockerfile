# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add curl
RUN go build -o main .
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .
COPY --from=builder /app/migrate .
COPY --from=builder /app/db/migration ./migration
COPY --from=builder /app/start.sh .
COPY --from=builder /app/wait-for.sh .
RUN chmod +x /app/start.sh
RUN chmod +x /app/wait-for.sh
EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
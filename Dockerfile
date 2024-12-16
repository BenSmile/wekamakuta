# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .
COPY --from=builder /app/db/migration ./db/migration
COPY --from=builder /app/start.sh .
COPY --from=builder /app/wait-for.sh .
RUN chmod +x /app/start.sh
RUN chmod +x /app/wait-for.sh
EXPOSE 8080
EXPOSE 9090
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
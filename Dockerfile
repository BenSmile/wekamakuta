# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .
EXPOSE 8080
CMD [ "./main" ]
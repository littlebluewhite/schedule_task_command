# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN apk update && \
    apk add --no-cache gcc && \
    rm -rf /var/lib/apt/lists/* && \
    apk add musl-dev
RUN go build -o main /app/cmd/api/main.go

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main /app/cmd/api/
COPY env env

RUN apk update
RUN apk add -U tzdata

ENV TZ=Asia/Taipei
EXPOSE 5487
CMD ["/app/cmd/api/main"]
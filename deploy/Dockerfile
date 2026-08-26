FROM golang:1.26.1-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/storyshard-app ./cmd/storyshard-app

FROM alpine:3.23

RUN apk add --no-cache ca-certificates \
    && addgroup -S storyshard \
    && adduser -S -G storyshard storyshard

WORKDIR /app
COPY --from=builder /out/storyshard-app ./storyshard-app

ENV CONFIG_PATH=/app/mnt \
    CONFIG_NAME=storyshard

USER storyshard
EXPOSE 17001

ENTRYPOINT ["./storyshard-app"]

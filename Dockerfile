# Build stage 
FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /personae-fasti

# Final stage
FROM alpine

WORKDIR /app

COPY --from=builder /personae-fasti /app/personae-fasti

EXPOSE 4121

ENV CONFIG_PATH="/app/mnt"
ENV CONFIG_NAME="fasti.json"

ENTRYPOINT ["./personae-fasti", ">>mnt/personae-fasti.out", "2>&1"]
CMD ["/bin/sh"]
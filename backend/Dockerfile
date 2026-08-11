FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o server ./cmd/server

FROM alpine:3.22

RUN apk --no-cache add ca-certificates tzdata
RUN addgroup -S -g 10001 app && adduser -S -D -H -u 10001 -G app app

WORKDIR /app

COPY --from=builder --chown=app:app /app/server .

USER app

EXPOSE 8080

CMD ["./server"]

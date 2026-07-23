FROM golang:1.26.5-alpine3.24 AS tools

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4

FROM golang:1.26.5-alpine3.24 AS builder

COPY --from=tools /go/bin/swag /usr/local/bin/swag

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN swag init -g cmd/api/main.go --parseDependency --parseInternal --md ./docs/md -o ./docs/gen

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/api

FROM alpine:3.24 AS runtime

RUN apk add --no-cache ca-certificates tzdata wget \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/server .

RUN chown -R appuser:appgroup /app

EXPOSE 8080

USER appuser

ENTRYPOINT ["./server"]
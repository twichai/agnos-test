FROM golang:1.25-alpine AS base

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

FROM base AS dev

EXPOSE 8080
CMD ["go", "run", "./cmd/main.go"]

FROM base AS builder
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/main.go

FROM alpine:3.22 AS prod

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/server /app/server

EXPOSE 8080
CMD ["/app/server"]
FROM golang:1.26-alpine AS base

WORKDIR /testy

RUN apk add --no-cache docker-cli

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Development
FROM base AS dev

RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", ".air.toml"]


# Production
FROM base AS prod

RUN go build -o /testy-server ./cmd/api

CMD ["/testy-server"]

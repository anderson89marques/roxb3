# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN apk update && apk add --no-cache make

COPY . .

RUN GOOS=linux \
    go install github.com/swaggo/swag/cmd/swag@latest \
    && go mod download \
    && swag init ./cmd/api/main.go -o docs --parseDependency --parseInternal

# Build API executable
RUN go build -o /app/cmd/api/api ./cmd/api/main.go

# Build CLI executable
RUN go build -o /app/cmd/cli/cli ./cmd/cli/main.go

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go get github.com/golang-migrate/migrate/v4/database/postgres

# Final stage
FROM alpine:latest

WORKDIR /app

RUN apk update && apk add --no-cache make

COPY --from=builder /app/cmd/api/api /app/cmd/api/api
COPY --from=builder /app/cmd/cli/cli /app/cmd/cli/cli
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY Makefile ./
COPY migrations ./migrations
COPY docs ./docs

EXPOSE 8080

RUN mkdir files

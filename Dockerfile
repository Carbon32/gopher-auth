FROM golang:1.27.1-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/auth-api ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app
COPY --from=build /out/auth-api /auth-api
COPY --from=build /out/migrate /migrate
COPY --from=build /src/cmd/migrate/migrations /app/cmd/migrate/migrations

USER app
EXPOSE 8080
ENTRYPOINT ["/auth-api"]

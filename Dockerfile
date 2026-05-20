FROM golang:1.24.4 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/mb-user-service ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=builder /out/mb-user-service /app/mb-user-service

EXPOSE 8080
ENTRYPOINT [ "/app/mb-user-service" ]
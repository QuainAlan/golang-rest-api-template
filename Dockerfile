# Align with go.mod and CI.
FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Pin this version instead of @latest.
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4

RUN swag init -g ./cmd/server/main.go -o ./docs

RUN CGO_ENABLED=1 go build -o /out/server ./cmd/server/main.go


FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN useradd --system --no-create-home --uid 10001 appuser

COPY --from=builder /out/server /app/server

USER appuser

EXPOSE 8080

CMD ["/app/server"]

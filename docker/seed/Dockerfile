# Build
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w" \
    -o seed .

# Runtime
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/seed /seed
COPY --from=builder /app/audio /audio

ENTRYPOINT ["/seed"]
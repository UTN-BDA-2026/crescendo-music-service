# Build
FROM golang:1.26.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w" \
    -o streaming .

# Runtime
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/streaming /streaming

EXPOSE 8081

USER nonroot:nonroot

ENTRYPOINT ["/streaming"]

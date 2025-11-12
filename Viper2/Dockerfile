# Multi-stage Dockerfile for Viper2
# Build stage
FROM golang:1.24-bullseye AS builder

WORKDIR /src

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the binary
ARG BINARY_NAME=viper2
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
ARG VERSION=dev
RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} \
    go build -ldflags "-s -w -X viper2/cmd.version=${VERSION}" -o /out/${BINARY_NAME} ./

# Optional server id: can be provided at build time
ARG SERVER_ID=""

# Final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=builder /out/viper2 /app/viper2
RUN chown app:app /app/viper2
USER app

EXPOSE 8080

# Expose server id as image label and environment variable if provided at build time.
# When building: docker build --build-arg SERVER_ID=your-server-id -t viper2:1.0 .
ARG SERVER_ID
ENV SERVER_ID=${SERVER_ID}
LABEL server.id=${SERVER_ID}

ENTRYPOINT ["/app/viper2"]

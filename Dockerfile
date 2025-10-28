# syntax=docker/dockerfile:1

FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.24-alpine AS build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Install swag CLI first (caches better)
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy mod files and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Generate Swagger docs (adjust main.go path if needed)
RUN swag init -g main.go -o ./docs

# Build binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-w -s" -o /server ./main.go

# ---- Runtime image ----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates

COPY --from=build /server /usr/local/bin/server
COPY --from=build /app/docs /app/docs

EXPOSE 8080
ENTRYPOINT ["server"]

FROM golang:1.25-alpine AS builder
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG TARGETVARIANT=
ARG GOARM=7
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}
ENV GOARM=${GOARM}
RUN go build -ldflags="-s -w" -o /irrigation-server .

FROM alpine:3.20 AS runtime
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /irrigation-server /app/irrigation-server
EXPOSE 8080
ENV TZ=UTC
ENTRYPOINT ["/app/irrigation-server"]

# Examples:
# Single-arch build (amd64):
# docker build --platform=linux/amd64 -t irrigation-server:amd64 .
# Single-arch build (arm64):
# docker build --platform=linux/arm64 -t irrigation-server:arm64 .
# Multi-arch build+push using buildx:
# docker buildx build --platform linux/amd64,linux/arm64 -t <repo>/irrigation-server:multi --push .
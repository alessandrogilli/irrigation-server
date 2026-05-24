FROM golang:1.25-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /irrigation-server .

FROM alpine:3.20 AS runtime
ENV TZ=Europe/Rome
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /irrigation-server /app/irrigation-server
COPY --from=builder /src/templates /app/templates
RUN mkdir -p /data
VOLUME /data
EXPOSE 8080

ENTRYPOINT ["/app/irrigation-server"]
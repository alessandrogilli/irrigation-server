FROM golang:1.25-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /irrigation-server .

FROM alpine:3.20 AS runtime
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /irrigation-server /app/irrigation-server
COPY --from=builder /src/templates /app/templates
EXPOSE 8080
ENV TZ=UTC
ENTRYPOINT ["/app/irrigation-server"]
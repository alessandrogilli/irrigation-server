FROM golang:1.25-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /irrigation-server .

FROM alpine:3.20 AS runtime
ENV TZ=Europe/Rome
RUN apk add --no-cache ca-certificates \
    mosquitto \
    mosquitto-clients \
    avahi \
    avahi-libs \
    dbus \
    dbus-libs \
    tzdata

WORKDIR /app
COPY --from=builder /irrigation-server /app/irrigation-server
COPY --from=builder /src/templates /app/templates

RUN mkdir -p /etc/avahi && \
    printf '[server]\nuse-ipv6=no\nallow-interfaces=eth0\n[publish]\npublish-workstation=yes\n' \
    > /etc/avahi/avahi-daemon.conf

RUN mkdir -p /etc/mosquitto && \
    echo "listener 1883 0.0.0.0" > /etc/mosquitto/mosquitto.conf && \
    echo "allow_anonymous true" >> /etc/mosquitto/mosquitto.conf

COPY raspberry-full-stack/start.sh /start.sh

RUN mkdir -p data
VOLUME data

EXPOSE 8080
EXPOSE 1883

RUN chmod +x /start.sh
CMD ["/start.sh"]

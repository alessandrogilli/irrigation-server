#!/bin/sh

mkdir -p /run/dbus

dbus-daemon --system --fork

avahi-daemon --daemonize

mosquitto -c /etc/mosquitto/mosquitto.conf --daemon

sleep 2

/app/irrigation-server

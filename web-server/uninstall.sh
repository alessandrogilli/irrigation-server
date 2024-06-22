#!/bin/bash

DESTINATION="$HOME/.local/irrigation-server/web-server"
SERVICE_NAME="irrigation-web-server.service"
SERVICE_PATH="/etc/systemd/system/$SERVICE_NAME"

sudo systemctl stop $SERVICE_NAME
sudo systemctl disable $SERVICE_NAME

sudo rm $SERVICE_PATH

rm -r $DESTINATION
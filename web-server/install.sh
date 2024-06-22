#!/bin/bash

set -e

DESTINATION="$HOME/.local/irrigation-server/web-server"
VENV=$DESTINATION/venv

if [ ! -d $VENV ]; then
    python3 -m venv $VENV
fi

source $VENV/bin/activate
pip install -r requirements.txt

cp server.py $DESTINATION/
cp -r templates $DESTINATION/

if [ ! -f ".env" ]; then
   echo "You must create the .env file before the installation."
   exit 1
fi
cp .env $DESTINATION/

SERVICE_NAME="irrigation-web-server.service"
SERVICE_PATH="/etc/systemd/system/$SERVICE_NAME"

SERVICE_CONTENT="[Unit]
Description=Gunicorn instance for the irrigation-server web server.
After=network.target

[Service]
Environment="PATH=$VENV/bin"
ExecStart=$VENV/bin/gunicorn -w 1 -b 0.0.0.0:10000 server:app
WorkingDirectory=$DESTINATION
User=$(whoami)
Group=$(id -gn)
Restart=always

[Install]
WantedBy=multi-user.target"

echo "$SERVICE_CONTENT" | sudo tee $SERVICE_PATH

sudo systemctl daemon-reload
sudo systemctl enable $SERVICE_NAME
sudo systemctl start $SERVICE_NAME


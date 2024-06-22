#!/bin/bash

set -e

DESTINATION="$HOME/.local/irrigation-server"
VENV=$DESTINATION/venv

if [ ! -d $VENV ]; then
    python3 -m venv $VENV
fi

source $VENV/bin/activate
pip install -r requirements.txt

cp irrigation_config.json $DESTINATION/

if [ ! -f ".env" ]; then
   echo "You must create the .env file before the installation."
   exit 1
fi
cp .env $DESTINATION/

echo "#!$VENV/bin/python" > $DESTINATION/irrigation-programmer
cat irrigation-programmer.py >> $DESTINATION/irrigation-programmer
chmod u+x $DESTINATION/irrigation-programmer

echo "#!$VENV/bin/python" > $DESTINATION/irrigation-server
cat irrigation-server.py >> $DESTINATION/irrigation-server
chmod u+x $DESTINATION/irrigation-server

if [ -f "$HOME/.bashrc" ]; then
    if ! grep -q "export PATH=\"\$PATH:$DESTINATION\"" "$HOME/.bashrc"; then
        echo "export PATH=\"\$PATH:$DESTINATION\"" >> "$HOME/.bashrc"
        echo "exporting PATH to ~/.bashrc"
    else
        echo "PATH export statement already exists in ~/.bashrc"
    fi
fi

if [ -f "$HOME/.zshrc" ]; then
    if ! grep -q "export PATH=\"\$PATH:$DESTINATION\"" "$HOME/.zshrc"; then
        echo "export PATH=\"\$PATH:$DESTINATION\"" >> "$HOME/.zshrc"
        echo "exporting PATH to ~/.zshrc"
    else
        echo "PATH export statement already exists in ~/.zshrc"
    fi
fi

export PATH=$PATH:$DESTINATION

SERVICE_NAME="irrigation-server.service"
SERVICE_PATH="/etc/systemd/system/$SERVICE_NAME"

SERVICE_CONTENT="[Unit]
Description=Service for the irrigation-server program.
After=network.target

[Service]
ExecStart=$DESTINATION/irrigation-server
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


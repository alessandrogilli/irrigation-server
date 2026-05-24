# Irrigation Server

Run locally:
```bash
export MQTT_BROKER="tcp://<your-broker-ip>:1883
go run main.go
```

Docker (lite):
```bash
docker build . -t irrigation-server:lite
docker run -d -e MQTT_BROKER="tcp://<your-broker-ip>:1883" -p 8080:8080 irrigation-server:lite
```

Docker (full):
```bash
docker build -f raspberry-full-stack/Dockerfile . -t irrigation-server:full
docker run -d --cap-add=NET_ADMIN --cap-add=NET_RAW --cap-add=NET_BROADCAST --hostname raspberrypi -e MQTT_BROKER="tcp://localhost:1883" -p 8080:8080 -p 1883:1883 irrigation-server:full
```

The **full** version replicates an old irrigation-server stack on Raspberry Pi with a built-in broker and includes Avahi to make the container respond to the mDNS name 'raspberrypi.local'.
Might be usefull if you want an all-in-one installation.
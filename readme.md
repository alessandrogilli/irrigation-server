# Irrigation Server

Run locally:
```bash
export MQTT_BROKER="tcp://<your-broker-ip>:1883
go run main.go
```

Docker:
```bash
docker build . -t irrigation-server:latest
docker run -d -e MQTT_BROKER="tcp://<your-broker-ip>:1883" -p 8080:8080
```
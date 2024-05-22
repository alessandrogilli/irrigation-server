import json
import schedule
import time
from datetime import datetime, timedelta
import paho.mqtt.publish as mqtt
from dotenv import load_dotenv
import os


def mqtt_publisher(payload):
    try:
        mqtt.single(topic=TOPIC, payload=payload, hostname=BROKER, port=PORT)
    except Exception as e:
        print(f"MQTT error: {e}")


def start_irrigate(line):
    payload = {"line": line, "cmd": "ON"}
    mqtt_publisher(json.dumps(payload))
    print(f"Starting line {line}.")


def stop_irrigate(line):
    payload = {"line": line, "cmd": "OFF"}
    mqtt_publisher(json.dumps(payload))
    print(f"Stopping line {line}.")


def load_configuration(file_path):
    try:
        with open(file_path) as f:
            return json.load(f)
    except:
        return {"programs": []}


def schedule_irrigation(config):
    for program in config["programs"]:
        line = program["line"]
        start_time = program["start_time"]
        duration = program["duration"]

        final_time = str(
            datetime.strptime(start_time, "%H:%M") + timedelta(minutes=duration)
        ).split()[1]

        schedule.every().day.at(start_time).do(start_irrigate, line=line)
        schedule.every().day.at(final_time).do(stop_irrigate, line=line)


if __name__ == "__main__":
    load_dotenv()
    BROKER = os.getenv("BROKER")
    PORT = int(os.getenv("PORT"))
    TOPIC = os.getenv("TOPIC")

    config_file_path = "irrigation_config.json"

    config = load_configuration(config_file_path)
    schedule_irrigation(config)

    while True:
        config = load_configuration(config_file_path)

        schedule.run_pending()
        try:
            time.sleep(1)
        except KeyboardInterrupt:
            print("Exit")
            break

        new_config = load_configuration(config_file_path)

        if config != new_config:
            print("Update!")
            schedule.clear()
            schedule_irrigation(new_config)

            for job in schedule.get_jobs():
                print(job)

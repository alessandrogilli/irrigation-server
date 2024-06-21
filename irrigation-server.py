import json
import schedule
import time
import logging
from logging.handlers import TimedRotatingFileHandler
from datetime import datetime, timedelta
import paho.mqtt.publish as mqtt
from dotenv import load_dotenv
import os
import argparse

# Create a logger
logger = logging.getLogger()
logger.setLevel(logging.INFO)

# Create console handler
console_handler = logging.StreamHandler()
console_handler.setLevel(logging.INFO)

# Create a TimedRotatingFileHandler
current_date = datetime.now().strftime("%Y-%m-%d")
log_file_name = f"irrigation-server-{current_date}.log"

file_handler = TimedRotatingFileHandler(
    filename=log_file_name,  # Log file name
    when="midnight",  # Roll over at midnight
    interval=1,  # Create a new log file daily
    backupCount=30,  # Keep up to 30 backup log files
)
file_handler.setLevel(logging.INFO)

# Create formatter
formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")

# Attach formatter to handlers
console_handler.setFormatter(formatter)
file_handler.setFormatter(formatter)

logger.addHandler(console_handler)
logger.addHandler(file_handler)


def mqtt_publisher(payload):
    try:
        mqtt.single(topic=TOPIC, payload=payload, hostname=BROKER, port=PORT)
    except Exception as e:
        logging.error(f"MQTT error: {e}")


def start_irrigate(line):
    payload = {"line": line, "cmd": "ON"}
    mqtt_publisher(json.dumps(payload))
    logging.info(f"Starting line {line}.")


def stop_irrigate(line):
    payload = {"line": line, "cmd": "OFF"}
    mqtt_publisher(json.dumps(payload))
    logging.info(f"Stopping line {line}.")


def load_configuration(file_path):
    try:
        with open(file_path) as f:
            return json.load(f)
    except Exception as e:
        logging.error(f"Error loading configuration: {e}")
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
    parser = argparse.ArgumentParser(description="Run the irrigation server.")
    parser.add_argument(
        "--mqtt-check",
        help="Check if the server can publish on the MQTT broker.",
        action="store_true",
    )

    load_dotenv()
    BROKER = os.getenv("BROKER")
    PORT = int(os.getenv("PORT"))
    TOPIC = os.getenv("TOPIC") + "/lines"

    args = parser.parse_args()
    if args.mqtt_check:
        payload = "MQTTCHECK"
        logging.info(f'Sending "{payload}" on {BROKER}:{PORT} at topic {TOPIC}')
        mqtt_publisher(payload)
        exit(1)

    logging.info(f"Starting irrigation-server")

    config_file_path = "irrigation_config.json"

    config = load_configuration(config_file_path)
    schedule_irrigation(config)

    while True:
        config = load_configuration(config_file_path)

        schedule.run_pending()
        try:
            time.sleep(1)
        except KeyboardInterrupt:
            logging.info("Exit")
            break

        new_config = load_configuration(config_file_path)

        if config != new_config:
            logging.info("Configuration updated!")
            schedule.clear()
            schedule_irrigation(new_config)

            for job in schedule.get_jobs():
                logging.info(f"Scheduled job: {job}")

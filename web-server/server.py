from flask import Flask, render_template, request
import paho.mqtt.client as mqtt
from dotenv import load_dotenv
import os

app = Flask(__name__)

load_dotenv()
MQTT_BROKER = os.getenv("BROKER")
MQTT_PORT = int(os.getenv("PORT"))
MQTT_TOPIC = os.getenv("TOPIC") + "/lines"

client = mqtt.Client()

def on_connect(client, userdata, flags, rc):
    print(f"Connected with result code {rc}")

client.on_connect = on_connect
client.connect(MQTT_BROKER, MQTT_PORT, 60)

@app.route('/irrigation-server/lines/<int:line_number>/<cmd>', methods=['POST'])
def irrigation_control(line_number, cmd):
    if cmd not in ["ON", "OFF"]:
        return {"error": "Invalid command"}, 400

    if line_number not in range(1, 5):
        return {"error": "Invalid line number"}, 400

    message = {"line": line_number, "cmd": cmd}
    client.publish(MQTT_TOPIC, str(message))
    return {"status": "success", "message": message}, 200

@app.route('/')
def home():
    return render_template('index.html')

if __name__ == '__main__':
    client.loop_start()
    app.run(host='0.0.0.0', port=10000)

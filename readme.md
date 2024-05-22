# irrigation-server
A simple tool to schedule irrigation programs.

## Core concepts
The application is composed primarly of two scripts: ```irrigation-programmer.py```, and ```irrigation-server.py```. The first script is used for the creation of the file ```irrigation_config.json```, the file containing all the informations about the programs being schedulated, like the line to activate, the start time and the duration. An example of a valid ```irrigation_config.json``` is the following:

```
{
    "programs": [
        {
            "program_id": 1,
            "line": 1,
            "start_time": "07:00",
            "duration": 5
        },
        {
            "program_id": 2,
            "line": 2,
            "start_time": "07:05",
            "duration": 5
        },
        {
            "program_id": 3,
            "line": 1,
            "start_time": "20:00",
            "duration": 5
        },
        {
            "program_id": 4,
            "line": 2,
            "start_time": "20:05",
            "duration": 5
        }
    ]
}
```

```irrigation-server.py```, once executed, reads the configuration file and schedules all the programs. A scheduled program consist of two scheduled action: turning on a line at a certain time, and turning it off after a specified duration. The implementation of the action at a low level is not a concern of the server, so its job is limited to sending a message on an MQTT Broker, allowing another component to subscribe and implement the ON/OFF function. The message sent on MQTT will have the following format:
```
{"line": 1, "cmd": "ON"}
```

Here's a text schema of the functioning:

```
irrigation-programmer.py -> irrigation_config.json -> irrigation-server.py -> MQTT Broker
```

## Quick start

1) Create the .env file, you can start from the .env.example:
    ```
    cp .env.example .env
    nano .env
    ```
2) Create a Python virtual environment, activate it and install the dependencies:
    ```
    python -m venv venv
    source venv/bin/activate
    pip install -r requirements.txt
    ```
3) Create some programs in the ```irrigation_config.json```. Here are the commands used to replicate the config file seen before in the guide:
    ```
    python irrigation-programmer.py --program_id 1 --line_number 1 --start_time 07:00 --duration 5
    python irrigation-programmer.py --program_id 2 --line_number 2 --start_time 07:05 --duration 5
    python irrigation-programmer.py --program_id 3 --line_number 1 --start_time 20:00 --duration 5
    python irrigation-programmer.py --program_id 4 --line_number 2 --start_time 20:05 --duration 5
    python irrigation-programmer.py --program_id 5 --line_number 3 --start_time 21:00 --duration 5

    ```
    Check if the config file has been updated successfully, you can ```cat irrigation_config.json```, or:
    ```
    python irrigation-programmer.py --get
    ```
    Oops, the program with id 5 was a mistake. Let's delete it:
    ```
    python irrigation-programmer.py --program_id 5 --line_number 3 --start_time 21:00 --duration 5 --delete
    ```
    or, more simply:

    ```
    python irrigation-programmer.py --program_id 5 --delete
    ```
    At any time, you can discover all the functionalities of this script using the -h or --help option:
    ```
    python irrigation-programmer.py --help
    ```
4) Run the ```irrigation-server.py```:
    ```
    python irrigation-server.py
    ```
5) You can also update the ```irrigation_config.json``` while this script is running. Every three seconds it checks if there are some changes in the config file, and if so, it reschedules all the programs according with the modified config file.

## Troubleshooting

### MQTT

You can use a MQTT client, like mosquitto-clients (for Ubuntu you can install via apt: ```sudo apt install mosquitto-clients```) and subscribe to the topic, port and broker specified in the .env file to check the correct configuration of the ```irrigation-server.py```.
Run the script with the flag --mqtt-check, to try publish a message in the specified broker, port and topic:
```
python irrigation-server.py --mqtt-check
```
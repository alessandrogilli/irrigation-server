import argparse
import json
import os


def load_configuration(file_path):
    with open(file_path) as f:
        return json.load(f)


def modify_configuration(config, program_id, line_number, start_time, duration):
    modified_config = config.copy()
    for program in modified_config["programs"]:
        if program["program_id"] == program_id:
            program["line"] = line_number
            program["start_time"] = start_time
            program["duration"] = duration
            break
    else:
        modified_config["programs"].append(
            {
                "program_id": program_id,
                "line": line_number,
                "start_time": start_time,
                "duration": duration,
            }
        )
    return modified_config


def delete_program(config, program_id):
    modified_config = config.copy()
    modified_config["programs"] = [
        program
        for program in modified_config["programs"]
        if program["program_id"] != program_id
    ]
    return modified_config


def save_configuration(config, file_path):
    with open(file_path, "w") as f:
        json.dump(config, f, indent=4)


def get_configuration(config):
    result = ""
    for program in config["programs"]:
        program_id = program["program_id"]
        line = program["line"]
        start_time = program["start_time"]
        duration = program["duration"]
        result += f"Program ID: {program_id} - Line {line} Starts at {start_time} duration {duration}\n"
    return result


def main():
    default_config_path = os.path.join(
        os.getenv("HOME"), ".local", "irrigation-server", "irrigation_config.json"
    )

    parser = argparse.ArgumentParser(description="Modify or delete irrigation program")
    parser.add_argument(
        "--file_path",
        help="Path to the configuration file",
        default=default_config_path,
    )
    parser.add_argument(
        "--init", action="store_true", help="Initialize the config file"
    )
    parser.add_argument("--get", action="store_true", help="Shows all the programs")
    parser.add_argument(
        "--program_id", nargs="?", type=int, help="The ID of a single program"
    )
    parser.add_argument("--line_number", nargs="?", type=int, help="Line number")
    parser.add_argument("--start_time", help="Start time")
    parser.add_argument("--duration", type=int, help="Duration")
    parser.add_argument(
        "--delete",
        action="store_true",
        help="Delete the program with matching line number",
    )
    args = parser.parse_args()

    if args.init:
        init = {"programs": []}
        save_configuration(init, args.file_path)
        print("Initialized config file")
        return

    try:
        config = load_configuration(args.file_path)
    except:
        print(f"File {args.file_path} not found. Create one with --init option")
        exit(1)

    if args.delete:
        if not args.program_id:
            parser.error("You must provide the id of the program to delete.")
        modified_config = delete_program(config, args.program_id)
        action = "deleted"
    elif args.get:
        print(get_configuration(config))
        return
    else:
        if (
            not args.program_id
            or not args.start_time
            or not args.duration
            or not args.line_number
        ):
            parser.error(
                "You must provide program id, line number, start time and duration to modify a program."
            )
        modified_config = modify_configuration(
            config, args.program_id, args.line_number, args.start_time, args.duration
        )
        action = "modified"

    save_configuration(modified_config, args.file_path)
    print(f"Irrigation program with id {args.program_id} {action} successfully.")


if __name__ == "__main__":
    main()

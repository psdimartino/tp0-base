import sys
import yaml


def generate_docker_yaml(clients):
    template = {
        "name": "tp0",
        "services": {
            "server": {
                "container_name": "server",
                "image": "server:latest",
                "entrypoint": "python3 /main.py",
                "environment": [
                    "PYTHONUNBUFFERED=1"
                ],
                "networks": ["testing_net"],
                "volumes": ["./server/config.ini:/config.ini"]
            }
        },
        "networks": {
            "testing_net": {
                "ipam": {
                    "driver": "default",
                    "config": [
                        {"subnet": "172.25.125.0/24"}
                    ]
                }
            }
        }
    }

    for i in range(1, clients + 1):
        client_name = f"client{i}"
        template["services"][client_name] = {
            "container_name": client_name,
            "image": "client:latest",
            "entrypoint": "/client",
            "environment": [
                f"CLI_ID={i}",
                "CLI_BET_PATH=./dataset.csv",
            ],
            "networks": ["testing_net"],
            "volumes": [
                "./client/config.yaml:/config.yaml",
                f"./.data/agency-{i}.csv:/dataset.csv"
            ],
            "depends_on": ["server"]
        }

    return template


def dump_yaml_to_file(filename):
    with open(filename, "w") as file:
        yaml.dump(docker_compose_template, file)


if __name__ == "__main__":
    print("Se inicio el generador de docker compose")
    if len(sys.argv) != 3:
        print("Uso: python3 docker-compose-generator.py <output_file> <num_clients>")
        sys.exit(1)

    compose_filename = sys.argv[1]
    num_clients = int(sys.argv[2])
    print("Se va crear un docker compose de nombre", compose_filename, "para", num_clients, "clientes")
    docker_compose_template = generate_docker_yaml(num_clients)
    dump_yaml_to_file(compose_filename)

# Nodevin

Nodevin allows anyone to run blockchain nodes effortlessly. It simplifies the process of setting up and managing nodes for various blockchains, ensuring they are always up-to-date with the latest software versions. With Nodevin, you can run nodes for Bitcoin, Ethereum, Litecoin, and more with ease.

## Features

- **Easy Setup:** Quickly set up blockchain nodes with a single command.
- **Automatic Updates:** The Nodevin daemon ensures your nodes are always running the latest software versions.
- **System Inspection:** Checks your system for the required Docker and Docker Compose versions, and installs them if necessary.
- **Cross-Platform Support:** Works on Linux, macOS, and Windows.
- **Detailed Logs:** Provides detailed logs for monitoring node activity and daemon status.

## Installation

To get started with Nodevin, you need to have Docker and Docker Compose installed on your system. If you don't have them installed, Nodevin can help you set them up.

### Prerequisites

- Docker 20+
- Docker Compose

### Installation

Download the latest version of Nodevin from the [releases page](https://github.com/fiftysixcrypto/nodevin/releases).

```sh
chmod +x nodevin
sudo mv nodevin /usr/local/bin/
```

## Getting Started

1. **Initialize Nodevin:**

Run the following command to initialize Nodevin and check your system capabilities:

```sh
nodevin init
```

This command will check if you have the proper versions of Docker and Docker Compose installed. If not, it will prompt you to install them.

2. **Start a Blockchain Node:**

Once Nodevin is initialized, you can start a blockchain node. For example, to start a Bitcoin node, run:

```sh
nodevin start bitcoin
```

## Nodevin Daemon

The Nodevin daemon runs in the background and ensures your blockchain nodes are always up-to-date. It checks for updates every hour and automatically updates the nodes if a new version is available.

### Starting the Daemon

To start the daemon, run:

```sh
nodevin daemon start
```

To start the daemon in detached mode (background):

```sh
nodevin daemon start -d
```

### Stopping the Daemon

To stop the daemon, run:

```sh
nodevin daemon stop
```

### Viewing Daemon Logs

To view the logs of the running daemon, run:

```sh
nodevin daemon logs
```

## Commands

Here are some of the key commands you can use with Nodevin:

- `nodevin init`: Initialize Nodevin and check system capabilities.
- `nodevin start <network>`: Start a blockchain node (e.g., `nodevin start bitcoin`).
- `nodevin stop <network>`: Stop a running blockchain node.
- `nodevin shell <network>`: Open a shell to the running node container.
- `nodevin logs <network>`: View logs for a specific node.
- `nodevin delete <volume-name-or-image-name>`: Delete a Docker volume or image.
- `nodevin request <network> --method <http-method> --params <json-data>`: Make an RPC request to a blockchain network.
- `nodevin daemon start`: Start the Nodevin daemon.
- `nodevin daemon stop`: Stop the Nodevin daemon.
- `nodevin daemon logs`: Show logs of the Nodevin daemon.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request or open an Issue on GitHub.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any questions or suggestions, feel free to contact us at [nodes@fiftysix.xyz](mailto:nodes@fiftysix.xyz).

---

Thank you for using Nodevin! Happy node running!

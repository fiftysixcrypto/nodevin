# Nodevin

Nodevin allows anyone to run blockchain nodes effortlessly. It simplifies the process of setting up and managing nodes for various blockchains, ensuring they are always up-to-date with the latest software versions. With Nodevin, you can run nodes for Bitcoin, Litecoin, and more with ease.

Our goal is to facilitate blockchain node standup and maintenance for every chain in the world.

## Features

- **Easy Setup:** Quickly set up blockchain nodes with a single command.
- **Automatic Updates:** The Nodevin daemon ensures your nodes are always running the latest software versions.
- **Maximum Customization:** Set unique ports, data storage, networking, images, or even run multiple nodes at once.
- **Cross-Platform Support:** Works on Linux, macOS, and Windows.

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

This command will check if you have the proper versions of Docker and Docker Compose installed. If not, it will install them for you.

```sh
nodevin init
```

2. **Start a Blockchain Node:**

Once Nodevin is initialized, you can start a blockchain node. For example, to start a Bitcoin node, run:

```sh
nodevin start bitcoin
```

3. **(Optional) - Advanced Features:**

Nodevin allows for full customization in node startup. View the full list of flags for configuration details [here](./docs). For example, this command runs a Bitcoin Testnet node with a specified command, docker image and tag (version), unique nodevin data directory, and more:

```sh
nodevin start bitcoin \
  --ord \
  --command="--rpcallowip=0.0.0.0/0 -rpcuser=user -rpcpassword=pass" \
  --testnet \
  --image=fiftysix/bitcoin-core \
  --version=27.0 \
  --ports="8332:8332,8333:8333,18332:18332,18333:18333" \
  --data-dir="~/Desktop" \
  --restart=always \
  --cpu-limit=2.0 \
  --mem-limit=1g \
  --cpu-reservation=1.0 \
  --mem-reservation=512m
```

## Snapshot Synchronization

Data snapshots are compressed archives of the state of a blockchain node. Using snapshot synchronization greatly speeds up the process of catching up with the network, as the node starts from downloaded data rather than trying to synchronize from the beginning of the blockchain. Running this command will use snapshot synchronization when starting up your node.

```
nodevin start litecoin --snapshot-sync
```

Snapshot synchronization can save up to **days** of node initialization.

### Adding Your Data Snapshot

Nodevin offers [a monthly subscription](https://nodevin.xyz/#/business) for networks interested in having their data snapshots integrated and universally accessible to all users.

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

Here are some key commands you can use with Nodevin:

- `nodevin init`: Initialize Nodevin and check system capabilities.
- `nodevin start <network>`: Start a blockchain node (e.g., `nodevin start bitcoin`).
- `nodevin stop <network>`: Stop a running blockchain node.

- `nodevin shell <network>`: Open a shell to the running node container.
- `nodevin logs <network>`: View logs for a specific node.
- `nodevin delete <volume-name-or-image-name>`: Delete local blockchain data.
- `nodevin request <network> --method <http-method> --params <json-data>`: Make an RPC request to a blockchain network.

More documentation [here](./docs/).

## Integrating Your Blockchain

Adding your blockchain to Nodevin requires a small one-time grant. For more information, visit [https://nodevin.xyz/#/business](https://nodevin.xyz/#/business).

### Nodevin Subscription

For more features, such as:
- Up to 90% faster sync times (`--snapshot-sync`)
- Universally available data snapshots
- Open source Docker images
- Docker image support and documentation

Visit [https://nodevin.xyz/#/business](https://nodevin.xyz/#/business).

## Uninstalling Nodevin

To uninstall Nodevin and delete all associated data and docker images, run the following commands:

```sh
nodevin cleanup
nodevin delete all
rm nodevin # wherever nodevin is located
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request or open an Issue on GitHub.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Contact

Socials:
- [Website](https://nodevin.xyz)
- [Business](https://nodevin.xyz/#/business)
- [Discord](https://discord.com/invite/XuhW2ykW3D)
- [Twitter/X](https://x.com/nodevin_)

This repository is currently maintained by [Fiftysix](https://fiftysix.xyz).

For any questions or suggestions, feel free to contact us at [nodes@fiftysix.xyz](mailto:nodes@fiftysix.xyz).

---

Thank you for using Nodevin! Happy node running!

# Nodevin

Nodevin allows anyone to run blockchain nodes effortlessly. It simplifies the process of setting up and managing nodes for various blockchains, ensuring they are always up-to-date with the latest software versions. With Nodevin, you can run nodes for Bitcoin, Litecoin, and IPFS with ease.

Our goal is to facilitate blockchain node standup and maintenance for every chain in the world.

## Features

- **Easy Setup:** Quickly set up blockchain nodes with a single command.
- **Automatic Updates:** Nodevin ensures your nodes are always running the latest software versions.
- **Maximum Customization:** Set unique ports, data storage, networking, images, or even run multiple nodes at once.
- **Cross-Platform Support:** Works on Linux, macOS, and Windows.

## Getting Started

### Setup

Stand up a blockchain node in three steps.

1. **Download Nodevin:**

Download the latest version of Nodevin from the [releases page](https://github.com/fiftysixcrypto/nodevin/releases).

2. **Initialize Nodevin and Docker:**

```sh
nodevin init
```

*For more information setting up Nodevin on Windows, read [these docs](./docs/windows-setup.md).*

3. **Start a Blockchain Node:**

Once Nodevin is initialized, you can start a blockchain node. For example, to start a Bitcoin node, run:

```sh
nodevin start bitcoin
```

*Nodevin stores blockchain data by default in `$HOME/.nodevin`.*

---

#### **(Optional) - Advanced Features:**

Nodevin allows for full customization in node startup. View the full list of flags and configuration details, including `.env` file setup [here](./docs/cli-commands.md). For example, this command runs a Bitcoin Testnet node with a specified command, docker image and tag (version), unique nodevin data directory, and more:

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

#### **(For Linux/MacOS) - Set Nodevin Permissions:**

After downloading Nodevin, you may need to set executable permissions and move it to a directory in `$PATH`.

```sh
chmod +x nodevin
sudo mv nodevin /usr/local/bin/
```

### More Documentation

Continue to learn about Nodevin:
- [Introduction to Nodevin](./docs/nodevin-intro.md) - Learn more about what Nodevin is and how it works.
- [Using Nodevin](./docs/cli-commands.md) - Documentation on how to use Nodevin. 
- [Nodevin Blog](...) - Coming soon... 

## Snapshot Synchronization

**NOTE: Snapshot Synchronization is currently disabled.**

Data snapshots are compressed archives of the state of a blockchain node. Using snapshot synchronization greatly speeds up the process of catching up with the network, as the node starts from downloaded data rather than trying to synchronize from the beginning of the blockchain. Running this command will use snapshot synchronization when starting up your node.

```
nodevin start litecoin --snapshot-sync
```

Snapshot synchronization can save up to **days** of node initialization.

## Integrating Your Blockchain

Adding your blockchain data to Nodevin requires a small one-time grant. For more information, reach out to us at [business@nodevin.xyz](mailto:business@nodevin.xyz).

### Nodevin Docker Images

Nodevin pulls images by default from [Docker Hub](https://hub.docker.com/u/fiftysix), with code located [here](https://github.com/fiftysixcrypto/node-images). The node image respository contains helpful resources including blockchain requirements and synchronization times, Docker installation steps, Docker compose files with documentation, and more.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request or open an Issue on GitHub.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Contact

Socials:
- [Website](https://nodevin.xyz)
- [Discord](https://discord.com/invite/XuhW2ykW3D)
- [Twitter/X](https://x.com/node_vin)

This repository is currently maintained by [Fiftysix](https://fiftysix.xyz).

For any questions or suggestions, feel free to contact us at [business@nodevin.xyz](mailto:business@nodevin.xyz).

---

Thank you for using Nodevin! Happy node running!

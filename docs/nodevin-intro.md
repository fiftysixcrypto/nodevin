# Introduction to Nodevin

**Nodevin** is an open-source tool that simplifies the process of running blockchain nodes. Whether you're a developer, enthusiast, or institution, Nodevin automates node setup, updates, and management—making it easy for anyone to participate in decentralized blockchain networks.

---

## Nodevin Quick Start Guide

1. **Check System Requirements**
   - **OS**: Linux, macOS, Windows (64-bit)
   - **Docker**: Version 20+ ([Get Docker](https://docs.docker.com/get-docker/))
   - **Docker Compose**: ([Install Docker Compose](https://docs.docker.com/compose/install/))
   - **CPU**: (Depends on blockchain)
   - **RAM**: (Depends on blockchain)
   - **Storage**: (Depends on blockchain)
   - **Network**: Reliable broadband with unlimited data

The `nodevin init` command can help you get set up and test system requirements.

2. **Install Nodevin**
   - Download the latest version from the [Nodevin GitHub Releases](https://github.com/fiftysixcrypto/nodevin/releases) or the [Nodevin Website](https://nodevin.xyz).
   - Ensure Docker and Docker Compose are installed (you can test this with `nodevin init`).

3. **Start a Blockchain Node**
   - Run: `nodevin list` to list all supported networks.
   - Run: `nodevin start <network>` (e.g., `nodevin start bitcoin`).

---

## What is a Node?

A **blockchain node** is a computer that participates in a blockchain network by maintaining the blockchain's ledger and validating transactions. Nodes ensure decentralization by distributing control across multiple participants.

### Types of Nodes:
- **Full Node**: Stores and validates the entire blockchain.
- **Light Node**: Stores only part of the blockchain and queries full nodes.
- **More**: [Read more](https://getblock.io/blog/blockchain-node-types/).

---

## How Nodevin Simplifies Node Setup

Nodevin removes the complexities of running blockchain nodes. Normally, setting up a node requires technical expertise and configuration. Nodevin does all this automatically by:
- Using Docker to isolate the node software.
- Handling updates with the daemon.
- Running nodes with a simple command: `nodevin start <blockchain>`.
- In depth information about each blockchain: `nodevin info`.

---

## Mining, Running a Node, and Staking: What's the Difference?

### Running a Node
- **Purpose**: Validates transactions and maintains the blockchain.
- **Rewards**: Typically no direct rewards, but essential for decentralization.
- **Hardware**: Requirements are not as restrictive, but depend on blockchain.
  
### Mining (Proof of Work)
- **Purpose**: Solves puzzles to create new blocks.
- **Rewards**: Earns block rewards and transaction fees.
- **Hardware**: Requires specialized hardware (e.g., ASICs for Bitcoin).

### Staking (Proof of Stake)
- **Purpose**: Locks cryptocurrency to validate transactions and create new blocks.
- **Rewards**: Earns rewards based on the amount staked.
- **Hardware**: Requires much less power and specialized equipment than mining.

Nodevin currently only supports running nodes to help users contribute to the network in various ways.

---

## What is Docker and Why Does Nodevin Use Docker?

**Docker** is a platform that bundles software into containers, which include everything needed to run the software. Nodevin uses Docker because:
- **Isolation**: Docker keeps your node separate from the rest of your system.
- **Standardization**: It ensures consistent performance across different systems.
- **Accessibility**: The [Docker Hub](https://hub.docker.com/u/fiftysix) allows anyone to access Nodevin built software versions and updates.

By using Docker, Nodevin makes deploying and managing blockchain nodes easy and efficient, no matter your technical background.

---

## More Resources
- [Blockchain Basics](https://www.investopedia.com/terms/b/blockchain.asp)
- [Getting Started with Docker](https://docs.docker.com/get-started/)
- [Nodevin Windows Installation](./windows-setup.md)
- [Nodevin Commands](./cli-commands.md)

# Nodevin CLI

Nodevin is a command-line interface (CLI) daemon that simplifies the setup, management, and running of blockchain nodes. Below is a comprehensive list of the available flags:

## USAGE:
```bash
nodevin command [command options] [arguments...]
```

### Example:
```bash
nodevin start bitcoin \
  --ord \
  --command="--rpcallowip=0.0.0.0/0 -rpcuser=user -rpcpassword=pass" \
  --testnet \
  --image=fiftysix/bitcoin-core \
  --version=27.0 \
  --container-name=bitcoin-node \
  --ports="8332:8332,8333:8333,18332:18332,18333:18333" \
  --volumes="bitcoin-core-data-2:/node/bitcoin-core" \
  --volume-definitions="bitcoin-core-data-2" \
  --volume-labels="nodevin.blockchain.software=bitcoin-core-testnet,remember-to-delete=yes" \
  --restart=always \
  --cpu-limit=2.0 \
  --mem-limit=1g \
  --cpu-reservation=1.0 \
  --mem-reservation=512m
```

## GLOBAL OPTIONS

### Command and Network Flags

- **`--command`**  
  *Description*: Specifies the node command with its own configuration options (e.g., RPC settings).  
  *Usage*: `--command="<command>"`  
  *Example*: `--command="--rpcallowip=0.0.0.0/0 -rpcuser=user -rpcpassword=pass"`  

- **`--testnet`**  
  *Description*: Run on an assumed test network.  
  *Default*: `false`  
  *Usage*: `--testnet`  

- **`--network`**  
  *Description*: Specify the network to run the node on (e.g., goerli, testnet3).  
  *Usage*: `--network=<network>`  
  *Example*: `--network=goerli`  

### Authentication Flags

- **`--rpc-user`**  
  *Description*: Username for JSON RPC authentication.  
  *Default*: `user`  
  *Usage*: `--rpc-user=<username>`  
  *Example*: `--rpc-user=myUser`  

- **`--rpc-pass`**  
  *Description*: Password for JSON RPC authentication.  
  *Default*: `fiftysix`  
  *Usage*: `--rpc-pass=<password>`  
  *Example*: `--rpc-pass=myPassword`  

- **`--cookie-auth`**  
  *Description*: Use node cookie file for authentication.  
  *Default*: `false`  
  *Usage*: `--cookie-auth`  

### Docker and Container Flags

- **`--image`**  
  *Description*: Docker image to use for the node (e.g., fiftysix/bitcoin-core).  
  *Usage*: `--image=<docker-image>`  
  *Example*: `--image=fiftysix/bitcoin-core`  

- **`--version`**  
  *Description*: Version of the Docker image to use (e.g., latest, 27.0).  
  *Usage*: `--version=<tag>`  
  *Example*: `--version=latest`  

- **`--container-name`**  
  *Description*: Name for the Docker container.  
  *Usage*: `--container-name=<name>`  
  *Example*: `--container-name=my-node`  

- **`--docker-networks`**  
  *Description*: Docker networks to connect to for the node.  
  *Usage*: `--docker-networks=<network1,network2,...>`  
  *Example*: `--docker-networks=network1,network2`  

- **`--network-driver`**  
  *Description*: Docker network driver to use.  
  *Usage*: `--network-driver=<driver>`  
  *Example*: `--network-driver=bridge`  

- **`--ports`**  
  *Description*: Specifies the port mappings for the node.  
  *Usage*: `--ports="<port1:port1,port2:port2,...>"`  
  *Example*: `--ports="8332:8332,8333:8333,18332:18332,18333:18333"`  

- **`--volumes`**  
  *Description*: Specifies the Docker volumes to mount.  
  *Usage*: `--volumes="<volume1:/path1,...>"`  
  *Example*: `--volumes="bitcoin-core-data-2:/node/bitcoin-core"`  

- **`--volume-definitions`**  
  *Description*: Defines the Docker volumes used in the compose file.  
  *Usage*: `--volume-definitions=<volume-name>`  
  *Example*: `--volume-definitions="bitcoin-core-data-2"`  

- **`--volume-labels`**  
  *Description*: Defines custom labels for the Docker volumes.  
  *Usage*: `--volume-labels="<label-key=value,...>"`  
  *Example*: `--volume-labels="nodevin.blockchain.software=bitcoin-core-testnet,remember-to-delete=yes"`  

- **`--restart`**  
  *Description*: Docker restart policy (e.g., always, no).  
  *Default*: `no`  
  *Usage*: `--restart=<policy>`  
  *Example*: `--restart=always`  

### Resource Management Flags

- **`--cpu-limit`**  
  *Description*: Maximum CPU limit for the container.  
  *Usage*: `--cpu-limit=<value>`  
  *Example*: `--cpu-limit=1.5`  

- **`--mem-limit`**  
  *Description*: Maximum memory limit for the container.  
  *Usage*: `--mem-limit=<value>`  
  *Example*: `--mem-limit=512m`  

- **`--cpu-reservation`**  
  *Description*: Reserved CPU for the container.  
  *Usage*: `--cpu-reservation=<value>`  
  *Example*: `--cpu-reservation=1.0`  

- **`--mem-reservation`**  
  *Description*: Reserved memory for the container.  
  *Usage*: `--mem-reservation=<value>`  
  *Example*: `--mem-reservation=256m`  

### Bitcoin and Litecoin Specific Flags

- **`--ord`**  
  *Description*: Run ordinal software `ord` alongside the Bitcoin node.  
  *Default*: `false`  
  *Usage*: `--ord`  

- **`--ord-litecoin`**  
  *Description*: Run ordinal software `ord` alongside the Litecoin node.  
  *Default*: `false`  
  *Usage*: `--ord-litecoin`  

# Nodevin CLI

Nodevin is a command-line interface (CLI) daemon that simplifies the setup, management, and running of blockchain nodes. Below is a comprehensive list of the available flags:

## USAGE:
```bash
nodevin command [command options] [arguments...]
```

## START OPTIONS

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

## REQUEST OPTIONS

### Example:
```bash
nodevin request bitcoin --method getblockcount
```

- **`nodevin request <network> --method <http-method> --params <json-data>`**  
  *Description*: Makes an RPC request to a specified blockchain network.  
  *Usage*: `nodevin request <network> --method <http-method> --params <json-data>`  
  *Example*: `nodevin request bitcoin --method getblockcount`  
  *Example with parameters*: `nodevin request bitcoin --method getblockheader --params '[\"00000000c937983704a73af28acdec37b049d214adbda81d7e2a3dd146f6ed09\"]'`  
  - **`--method`**  
    *Description*: Specifies the HTTP method to use for the request (e.g., GET, POST).  
    *Usage*: `--method=<http-method>`  
    *Example*: `--method=getblockcount`  
  - **`--params`**  
    *Description*: JSON data to send as parameters in the request body.  
    *Usage*: `--params=<json-data>`  
    *Example*: `--params='[\"param1\", \"param2\"]'`  
  - **`--header`**  
    *Description*: Adds optional extra headers to the request.  
    *Usage*: `--header=<key:value>`  
  - **`--endpoint`**  
    *Description*: Specify an optional API endpoint (default is `http://127.0.0.1`).  
    *Usage*: `--endpoint=<url>`  
  - **`--port`**  
    *Description*: Optional port to override the default RPC port for the network.  
    *Usage*: `--port=<port>`

## SHELL OPTIONS

### Example:
```bash
nodevin shell bitcoin
```

- **`nodevin shell <network>`**  
  *Description*: Opens an interactive shell in the running container for the specified blockchain network.  
  *Usage*: `nodevin shell <network>`  
  *Example*: `nodevin shell bitcoin`  
  - **`--detach`**: Run the shell in detached mode (in the background).  
    *Usage*: `nodevin shell <network> --detach`  
  - **`--docker-user=<user>`**: Specifies the username or UID to run the shell as inside the container.  
    *Usage*: `nodevin shell <network> --docker-user=root`  
  - **`--workdir=<path>`**: Set the working directory inside the container.  
    *Usage*: `nodevin shell <network> --workdir=/node`  
  - **`--env=<key=value>`**: Set environment variables.  
    *Usage*: `nodevin shell <network> --env=KEY=VALUE`  
  - **`--env-file=<file>`**: Read environment variables from a file.  
    *Usage*: `nodevin shell <network> --env-file=./env.list`  
  - **`--privileged`**: Run the shell with extended privileges.  
    *Usage*: `nodevin shell <network> --privileged`  

## STOP OPTIONS

### Example:
```bash
nodevin stop bitcoin
```

- **`nodevin stop <network>`**  
  *Description*: Stops a running blockchain node for the specified network.  
  *Usage*: `nodevin stop <network>`  
  *Example*: `nodevin stop bitcoin`  
  If no network is specified, it will list the available networks and usage examples.  
  - **`stop all`**  
    *Description*: Stops **all** Docker containers.  
    *Usage*: `nodevin delete all` 
  - **`--testnet`**: Stops the node running on a test network.  
    *Usage*: `nodevin stop <network> --testnet`  
  - **`--network=<network>`**: Specify a custom network.  
    *Usage*: `nodevin stop <network> --network="goerli"`  

## DELETE OPTIONS

### Example:
```bash
nodevin delete all
```

- **`nodevin delete <volume-name-or-image-name>`**  
  *Description*: Deletes a Docker volume and its associated images, or deletes a specific Docker image.  
  *Usage*: `nodevin delete <volume-name-or-image-name>`  
  *Example*: `nodevin delete fiftysix/bitcoin-core:27.0`  
  If no volume or image name is provided, it will list the available volumes and images.  
  - **`delete all`**  
    *Description*: Deletes all Docker volumes with the label `nodevin.blockchain.software`.  
    *Usage*: `nodevin delete all`  

## LOG OPTIONS

### Example:
```bash
nodevin logs bitcoin
```

- **`nodevin logs <network>`**  
  *Description*: Fetches logs for a running blockchain node.  
  *Usage*: `nodevin logs <network>`  
  *Example*: `nodevin logs bitcoin`  
  - **`--follow`**: Continuously stream the logs in real-time.  
    *Usage*: `nodevin logs <network> --follow`  
  - **`--tail=<value>`**: Shows a specific number of lines from the end of the logs.  
    *Usage*: `nodevin logs <network> --tail=100`  

# Nodevin Daemon

Nodevin is a command-line interface (CLI) daemon that simplifies the setup, management, and running of blockchain nodes. Below is a comprehensive list of the available flags for the built-in daemon:

## USAGE:
```bash
nodevin command [command options] [arguments...]
```

### Example:
```bash
nodevin start daemon -d
```

## DAEMON COMMANDS

### Command and Network Flags

- **`--detach` or `-d`**  
  *Description*: Run the daemon detached, in the background.  
  *Default*: `false`  
  *Usage*: `nodevin start --detach` or `nodevin start -d`  

### Daemon Control Flags

- **`nodevin stop daemon`**  
  *Description*: Stops the running Nodevin daemon.  
  *Usage*: `nodevin stop daemon`  
  This will signal the daemon process to terminate and remove the associated PID file.  
  If the daemon is not running, an error message will be displayed indicating that the daemon is not found.  

### Daemon Logs Flags

- **`nodevin daemon logs`**  
  *Description*: Manage and view the logs of the Nodevin daemon.  
  *Usage*: `nodevin daemon logs`  
  By default, this command displays the entire content of the daemon's log file.  
  Log files are located at the configured log file path, typically `nodevin.log`.  
  - **`--tail` or `-t <number>`**  
    *Description*: Display the last `<number>` of lines from the end of the log file.  
    *Usage*: `nodevin daemon logs --tail=100`  
  - **`--clear` or `-c`**  
    *Description*: Clear the current log file.  
    *Usage*: `nodevin daemon logs --clear`  
  - **`--follow` or `-f`**  
    *Description*: Follow the log file in real time (similar to `tail -f`).  
    *Usage*: `nodevin daemon logs --follow`  

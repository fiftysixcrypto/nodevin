package nodes

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/internal/utils"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View your Nodevin in a fun and artistic way!",
	Run: func(cmd *cobra.Command, args []string) {
		nodeSizes := fetchNodeSizes()
		nodeStats := fetchNodeStats()
		displayNodevinArt(nodeSizes, nodeStats)
	},
}

type NodeData struct {
	Network     string `json:"network"`
	Name        string `json:"name"`
	Uptime      int    `json:"uptime"`
	Peers       int    `json:"peers"`
	LatestBlock int    `json:"latest_block"`
}

func fetchNodeSizes() map[string]int64 {
	nodevinDataDir, err := utils.GetNodevinDataDir()
	if err != nil {
		logger.LogError("Failed to find Nodevin data directory: " + err.Error())
		return nil
	}

	networks := utils.GetAllSupportedNetworks()
	if networks == "" {
		fmt.Println("No data found.")
		return nil
	}

	nodeSizes := make(map[string]int64)
	for _, network := range strings.Split(networks, ", ") {
		containerName, exists := utils.GetDefaultLocalMappedContainerName(network)
		if !exists {
			logger.LogError("Unsupported blockchain network: " + network)
			continue
		}

		networkDir := filepath.Join(nodevinDataDir, containerName)
		size, err := getDirectorySize(networkDir)
		if err == nil {
			nodeSizes[network] = size
		}
	}
	return nodeSizes
}

func fetchNodeStats() []NodeData {
	args := []string{"ps", "--format", "{{json .}}"}
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.LogError("Failed to fetch Docker container information: " + err.Error())
		return nil
	}

	containers := strings.Split(string(output), "\n")
	var nodes []NodeData
	for _, containerJSON := range containers {
		if strings.TrimSpace(containerJSON) == "" {
			continue
		}
		var container ContainerInfo
		if err := json.Unmarshal([]byte(containerJSON), &container); err != nil {
			logger.LogError("Failed to parse container JSON: " + err.Error())
			continue
		}
		if strings.Contains(container.Names, "watchtower") {
			continue
		}

		peers := getPeers(container.Names)
		localLatestBlock, _ := getLatestBlocks(container.Names)
		nodes = append(nodes, NodeData{
			Network:     getSoftwareNetworkName(container.Names),
			Name:        container.Names,
			Uptime:      extractUptime(container.Status),
			Peers:       peers,
			LatestBlock: localLatestBlock,
		})
	}
	return nodes
}

func getNodevinName(network string) string {
	switch strings.ToLower(network) {
	case "bitcoin-core":
		return "Bitvin"
	case "litecoin-core":
		return "Litevin"
	default:
		return "Nodevin"
	}
}

func getSoftwareNetworkName(softwareName string) string {
	switch softwareName {
	case "bitcoin-core":
		return "bitcoin"
	case "litecoin-core":
		return "litecoin"
	default:
		return ""
	}
}

func extractUptime(status string) int {
	if strings.Contains(status, "days") {
		var days int
		fmt.Sscanf(status, "%d days", &days)
		return days
	} else if strings.Contains(status, "hours") {
		return 1
	}
	return 0
}

func displayNodevinArt(nodeSizes map[string]int64, nodes []NodeData) {
	if len(nodes) == 0 {
		fmt.Printf("\nVin Status:\n\n")
		fmt.Printf("   (v_v)  \n  /[   ]\\   \n   [   ]   \n Feeling Lonely... \n\n")
		return
	}

	for _, node := range nodes {
		nodeSize := nodeSizes[node.Network]
		vinName := getNodevinName(node.Name)
		fmt.Printf("\n%s Status:\n\n", vinName)

		var art string
		switch {
		case node.Uptime >= 7 && node.Peers >= 3 && node.LatestBlock >= 500000 && nodeSize > 100000000000:
			art = "   (n_v)  \n  \\[***]/  \n   [###]  \nNetwork Supremacy! "
		case node.Peers >= 3 && node.LatestBlock >= 10000 && nodeSize > 100000000000:
			art = "   (n_v)  \n  \\[**#]/  \n   [##*]  \nSupercharged! "
		case nodeSize > 600000000000:
			art = "   (n_v)  \n  \\[^^^]/  \n   [^^^]  \nRejoining Crew... "
		case nodeSize > 100000000000:
			art = "   (n_v)  \n  \\[-*-]/  \n   [-*-]  \nGetting Set... "
		case node.Uptime >= 3 && node.Peers >= 3:
			art = "   (n_v)  \n  \\[ **]/  \n   [ ##]  \nGrowing Strong! "
		case node.Uptime >= 2 && node.Peers >= 3:
			art = "   (n_v)  \n  \\[ * ]/  \n   [ # ]  \nWarming Up! "
		case node.Uptime >= 1 && node.Peers >= 3:
			art = "   (n_v)  \n  \\[ .#]/   \n   [ .#]   \nPowering On! "
		case node.Uptime >= 0 && node.Peers >= 1:
			art = "   (n_v)  \n  \\[ . ]/   \n   [ . ]   \nJust Hatched! "
		default:
			art = "   (n_v)  \n  /[   ]\\   \n   [   ]   \nGetting Situated... "
		}

		fmt.Printf("%s\n\n", art)
		displayNodeDetails(node, nodeSize)
	}
}

func displayNodeDetails(node NodeData, size int64) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "| METRIC\t VALUE")
	fmt.Fprintf(w, "| Network\t %s\n", node.Network)
	fmt.Fprintf(w, "| Uptime\t %d days\n", node.Uptime)
	fmt.Fprintf(w, "| Peers\t %d\n", node.Peers)
	fmt.Fprintf(w, "| Latest Block\t %d\n", node.LatestBlock)
	fmt.Fprintf(w, "| Size\t %s\n", utils.GetSizeDescription(size))
	w.Flush()
}

package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/curveballdaniel/nodevin/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var requestCmd = &cobra.Command{
	Use:   "request [network]",
	Short: "Make an RPC request to a blockchain network",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		network := args[0]

		method := viper.GetString("method")
		params := viper.GetString("params")
		headers := viper.GetString("header")
		endpoint := viper.GetString("endpoint")
		port := viper.GetInt("port")
		user := viper.GetString("rpc-user")
		pass := viper.GetString("rpc-pass")

		if method == "" {
			logger.LogError("HTTP method is required.")
			printUsageAndExample()
			return
		}

		if endpoint == "" {
			endpoint = "http://127.0.0.1"
		}

		if port == 0 {
			port = utils.NetworkDefaultRPCPorts()[network]
		}

		url := fmt.Sprintf("%s:%d", endpoint, port)

		if err := makeRequest(network, url, method, params, headers, user, pass); err != nil {
			logger.LogError("Failed to make request: " + err.Error())
		}
	},
}

func printUsageAndExample() {
	fmt.Println("Usage: nodevin request [network] --method <http-method> --params <json-data> --header <optional-extra-headers> --endpoint <optional-api-endpoint>")
	fmt.Println("Example: nodevin request bitcoin --method getblockheader --params '[\"00000000c937983704a73af28acdec37b049d214adbda81d7e2a3dd146f6ed09\"]'")
}

func makeRequest(network, url, method, params, headers, user, pass string) error {
	jsonData := map[string]interface{}{
		"jsonrpc": "1.0",
		"id":      "nodevin",
		"method":  method,
	}

	var jsonParams []interface{}
	if params != "" {
		if err := json.Unmarshal([]byte(params), &jsonParams); err != nil {
			return fmt.Errorf("invalid params: %w", err)
		}
	} else {
		jsonParams = []interface{}{}
	}

	jsonData["params"] = jsonParams

	jsonValue, _ := json.Marshal(jsonData)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if headers != "" {
		for _, header := range strings.Split(headers, ",") {
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}

	if user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 401 {
			return fmt.Errorf("request failed with status code %d: Unauthorized\nMaybe consider using the --user and --pass flags?", resp.StatusCode)
		}
		return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Print(string(body))

	return nil
}

func init() {
	requestCmd.Flags().StringP("method", "m", "", "HTTP method to use for the request")
	requestCmd.Flags().StringP("params", "p", "", "JSON data to send in the request body")
	requestCmd.Flags().StringP("header", "H", "", "Optional extra headers")
	requestCmd.Flags().StringP("endpoint", "e", "http://127.0.0.1", "Optional API endpoint")
	requestCmd.Flags().IntP("port", "P", 0, "Optional port to override the default")
	requestCmd.Flags().StringP("user", "u", "", "Optional user for authentication")
	requestCmd.Flags().StringP("pass", "w", "", "Optional password for authentication")

	viper.BindPFlag("method", requestCmd.Flags().Lookup("method"))
	viper.BindPFlag("params", requestCmd.Flags().Lookup("params"))
	viper.BindPFlag("header", requestCmd.Flags().Lookup("header"))
	viper.BindPFlag("endpoint", requestCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("port", requestCmd.Flags().Lookup("port"))
	viper.BindPFlag("user", requestCmd.Flags().Lookup("user"))
	viper.BindPFlag("pass", requestCmd.Flags().Lookup("pass"))
}

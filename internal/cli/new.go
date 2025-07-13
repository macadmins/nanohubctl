package cli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

type ProvisioningResponse struct {
	Name              string `json:"name"`
	EnrollmentProfile string `json:"enrollment_profile"`
	APIKey            string `json:"api_key"`
}

func newCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new nanohub instance",
		Long:  "Create a new nanohub instance with the provided token",
		RunE:  newCmdFn,
	}

	newCmd.Flags().StringP("token", "t", "", "Token for authentication")
	newCmd.MarkFlagRequired("token")

	return newCmd
}

func newCmdFn(cmd *cobra.Command, args []string) error {
	// Get home directory and create config path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}
	
	configDir := filepath.Join(homeDir, ".nanohubctl")
	configFile := filepath.Join(configDir, "config.json")
	
	// Check if config file already exists
	if _, err := os.Stat(configFile); err == nil {
		// Config exists, read and display it
		data, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to read existing config: %v", err)
		}
		
		var existingResp ProvisioningResponse
		if err := json.Unmarshal(data, &existingResp); err != nil {
			return fmt.Errorf("failed to parse existing config: %v", err)
		}
		
		fmt.Printf("You already have a nanohub instance configured:\n\n")
		printInstanceInfo(existingResp)
		return nil
	}
	
	// Config doesn't exist, proceed with creating new instance
	token, err := cmd.Flags().GetString("token")
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 60 * time.Second, // Set timeout to 60 seconds since it can take > 30s
	}

	req, err := http.NewRequest("POST", "https://provisioning.macadmins.io/new", nil)
	if err != nil {
		return err
	}

	// Set basic auth header
	username := "nanohub"
	auth := username + ":" + token
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encodedAuth)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Check if the response status is not successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	var provisioningResp ProvisioningResponse
	if err := json.Unmarshal(body, &provisioningResp); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v\nResponse body: %s", err, string(body))
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}
	
	// Save the response to config file
	configData, err := json.MarshalIndent(provisioningResp, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config data: %v", err)
	}
	
	if err := os.WriteFile(configFile, configData, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	// Print formatted output
	fmt.Printf("Successfully created new nanohub instance:\n\n")
	printInstanceInfo(provisioningResp)
	return nil
}

func printInstanceInfo(resp ProvisioningResponse) {
	fmt.Printf("Instance Name:               %s\n", resp.Name)
	fmt.Printf("API Key:                     %s\n", resp.APIKey)
	fmt.Printf("Enrollment Profile:          %s\n", resp.EnrollmentProfile)
	fmt.Printf("\nTo configure nanohubctl, run:\n\n")
	fmt.Printf("export NANOHUB_URL=https://%s\n", resp.Name)
	fmt.Printf("export NANOHUB_API_KEY=%s\n", resp.APIKey)
}

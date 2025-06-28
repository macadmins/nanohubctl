package nanocmd

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/macadmins/nanohubctl/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Example URL
// http://example.com/api/v1/nanocmd/workflow/io.micromdm.wf.devinfolog.v1/start?id=9876-5432-1012

// StartWorkflow posts a workflow to the nanocmd API endpoint
func StartWorkflow(workflowName, clientID string) (*http.Response, error) {
	baseUrl, err := utils.GetNanoCMDUrl()
	if err != nil {
		return nil, fmt.Errorf("failed to get nanocmd URL: %w", err)
	}

	// Build the path: /api/v1/nanocmd/workflow/{workflowName}/start
	baseUrl.Path = path.Join(baseUrl.Path, "workflow", workflowName, "start")

	// Add the client ID as a query parameter
	params := url.Values{}
	params.Add("id", clientID)
	baseUrl.RawQuery = params.Encode()

	var resp *http.Response
	err = postReq(baseUrl.String(), &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// WorkflowCmd creates the workflow command
func WorkflowCmd() *cobra.Command {
	workflowCmd := &cobra.Command{
		Use:     "workflow <workflow-name> [client-id]",
		Short:   "Start a workflow for a specific client",
		Long:    "Start a workflow by name for a specific client ID. If client-id is not provided, uses --client-id flag value.",
		Args:    cobra.RangeArgs(1, 2),
		PreRunE: utils.ApplyPreExecFn,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			workflowName := args[0]
			var clientID string

			if len(args) == 2 {
				clientID = args[1]
			} else {
				clientID = viper.GetString("client_id")
				if clientID == "" {
					return cmd.Help()
				}
			}

			resp, err := StartWorkflow(workflowName, clientID)
			if err != nil {
				return fmt.Errorf("failed to start workflow: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				fmt.Printf("Workflow %s started successfully for client %s\n", workflowName, clientID)
			} else {
				fmt.Printf("Failed to start workflow. Status: %s\n", resp.Status)
			}

			return nil
		},
	}

	return workflowCmd
}

package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"path"

	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func statusCmd() *cobra.Command {
	statusCmd := &cobra.Command{
		Use:     "status",
		Short:   fmt.Sprintf("This verb handles all status related operations"),
		Long:    fmt.Sprintf("This verb handles all status related operations"),
		PreRunE: applyPreExecFn,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	statusCmd.PersistentFlags().StringP("ID", "i", "", "Identifier of the client to return status for")
	statusCmd.MarkPersistentFlagRequired("ID")
	statusCmd.AddCommand(
		declarationStatusCmd(),
		errorsCmd(),
		valuesCmd(),
	)

	return statusCmd
}

// getCmd handles getting statuss on the server
func declarationStatusCmd() *cobra.Command {
	declarationStatusCmd := &cobra.Command{
		Use:     "declarations",
		Short:   fmt.Sprintf("List declarations for a specified enrollment ID"),
		Long:    fmt.Sprintf("List declarations for a specified enrollment ID"),
		PreRunE: applyPreExecFn,
		RunE:    StatusFn,
	}

	return declarationStatusCmd
}

// getCmd handles getting statuss on the server
func errorsCmd() *cobra.Command {
	errorsCmd := &cobra.Command{
		Use:     "errors",
		Short:   fmt.Sprintf("List errors for a specified enrollment ID"),
		Long:    fmt.Sprintf("List errors for a specified enrollment ID"),
		PreRunE: applyPreExecFn,
		RunE:    StatusFn,
	}

	return errorsCmd
}

// getCmd handles getting statuss on the server
func valuesCmd() *cobra.Command {
	valuesCmd := &cobra.Command{
		Use:     "values",
		Short:   fmt.Sprintf("List values for a specified enrollment ID"),
		Long:    fmt.Sprintf("List values for a specified enrollment ID"),
		PreRunE: applyPreExecFn,
		RunE:    StatusFn,
	}

	return valuesCmd
}

func StatusFn(cmd *cobra.Command, statuss []string) error {
	clientID, err := cmd.Flags().GetString("ID")
	// fmt.Println(cmd.Use)
	if err != nil {
		return err
	}
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}
	switch cmd.Use {
	case "declarations":
		ddmUrl.Path = path.Join(ddmUrl.Path, "v1/declaration-status", clientID)
	case "values":
		ddmUrl.Path = path.Join(ddmUrl.Path, "v1/status-values", clientID)
	case "errors":
		ddmUrl.Path = path.Join(ddmUrl.Path, "v1/status-errors", clientID)
	default:
		return fmt.Errorf("%s is not a valid status type", cmd.Use)
	}
	var resp *http.Response
	err = getReq(ddmUrl.String(), &resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	var jsonResponse interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return err
	}
	fmt.Println(PrettyJsonPrint(jsonResponse))
	return nil
}

package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strings"

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

	// statusCmd.PersistentFlags().StringP("ID", "i", "", "Identifier of the client to return status for")
	// statusCmd.MarkPersistentFlagRequired("ID")
	statusCmd.AddCommand(
		declarationStatusCmd(),
		errorsCmd(),
		valuesCmd(),
	)

	return statusCmd
}

// declarationStatusCmd lists declarations for a specified device ID
func declarationStatusCmd() *cobra.Command {
	declarationStatusCmd := &cobra.Command{
		Use:     "declarations [--client-id $ID]",
		Short:   fmt.Sprintf("List declarations for a specified device ID"),
		Long:    fmt.Sprintf("List declarations for a specified device ID"),
		PreRunE: applyPreExecFn,
		RunE:    StatusFn,
	}

	return declarationStatusCmd
}

// errorsCmd Lists errors for a specified device ID
func errorsCmd() *cobra.Command {
	errorsCmd := &cobra.Command{
		Use:     "errors [--client-id $ID]",
		Short:   fmt.Sprintf("List errors for a specified device ID"),
		Long:    fmt.Sprintf("List errors for a specified device ID"),
		PreRunE: applyPreExecFn,
		RunE:    StatusFn,
	}

	return errorsCmd
}

// valuesCmd lists all values for a specified device ID
func valuesCmd() *cobra.Command {
	valuesCmd := &cobra.Command{
		Use:     "values [--client-id $ID]",
		Short:   fmt.Sprintf("List values for a specified device ID"),
		Long:    fmt.Sprintf("List values for a specified device ID"),
		PreRunE: applyPreExecFn,
		RunE:    StatusFn,
	}

	return valuesCmd
}

// StatusFn handles all logic for the various status commands
func StatusFn(cmd *cobra.Command, statuss []string) error {
	clientID := viper.GetString("client_id")
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}
	cmdVerb := strings.Split(cmd.Use, " ")[0]
	switch cmdVerb {
	case "declarations":
		ddmUrl.Path = path.Join(ddmUrl.Path, "v1/declaration-status", clientID)
	case "values":
		ddmUrl.Path = path.Join(ddmUrl.Path, "v1/status-values", clientID)
	case "errors":
		ddmUrl.Path = path.Join(ddmUrl.Path, "v1/status-errors", clientID)
	default:
		return fmt.Errorf("%s is not a valid status type", cmdVerb)
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

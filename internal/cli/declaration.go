package cli

import (
	"fmt"
	"os"
	"path"

	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var jsonPayload string

func declarationCmd() *cobra.Command {
	declarationCmd := &cobra.Command{
		Use:     "declaration",
		Short:   fmt.Sprintf("This verb handles all declaration related operations"),
		Long:    fmt.Sprintf("This verb handles all declaration related operations"),
		PreRunE: applyPreExecFn,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}
	declarationCmd.AddCommand(
		createCmd(),
		// getCmd(),
		// deleteCmd(),
	)

	return declarationCmd
}

// createCmd handles creating declarations on the server
func createCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:     "create",
		Short:   fmt.Sprintf("create a declaration"),
		Long:    fmt.Sprintf("create a declaration"),
		PreRunE: applyPreExecFn,
		RunE:    createFn,
	}

	createCmd.PersistentFlags().StringVarP(&jsonPayload, "json", "j", "", "json payload to create a declaration")
	createCmd.MarkPersistentFlagRequired("json")

	return createCmd
}

func createFn(cmd *cobra.Command, declarations []string) error {
	jsonPath, err := cmd.Flags().GetString("json")
	if err != nil {
		return err
	}
	jsonBytes, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	fmt.Printf("Creating declaration using %s\n", jsonPath)
	fmt.Println(viper.GetString("url"))
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}
	ddmUrl.Path = path.Join(ddmUrl.Path + "v1/declarations")
	var resp *http.Response
	err = putReq(ddmUrl.String(), jsonBytes, &resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	return nil
}

package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
		getCmd(),
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

	createCmd.Flags().StringP("json", "j", "", "json payload to create a declaration")
	createCmd.MarkFlagRequired("json")

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
	ddmUrl.Path = path.Join(ddmUrl.Path, "v1/declarations")
	var resp *http.Response
	err = putReq(ddmUrl.String(), jsonBytes, &resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	return nil
}

// getCmd handles getting declarations on the server
func getCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:     "get",
		Short:   fmt.Sprintf("create a declaration"),
		Long:    fmt.Sprintf("create a declaration"),
		PreRunE: applyPreExecFn,
		RunE:    getFn,
	}

	getCmd.Flags().StringP("identifier", "i", "", "Identifier of the declaration to retrieve")
	getCmd.MarkFlagRequired("identifier")

	return getCmd
}

func getFn(cmd *cobra.Command, declarations []string) error {
	identifier, err := cmd.Flags().GetString("identifier")
	if err != nil {
		return err
	}
	fmt.Printf("Getting declaration for identifier %s\n", identifier)
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}
	ddmUrl.Path = path.Join(ddmUrl.Path, "v1/declarations", identifier)
	var resp *http.Response
	err = getReq(ddmUrl.String(), &resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return err
	}
	fmt.Println(PrettyJsonPrint(jsonResponse))
	return nil
}

package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		createDeclarationCmd(),
		getDeclarationCmd(),
		deleteDeclarationCmd(),
	)

	return declarationCmd
}

// createCmd handles creating declarations on the server
func createDeclarationCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:     "create",
		Short:   fmt.Sprintf("create a declaration"),
		Long:    fmt.Sprintf("create a declaration"),
		PreRunE: applyPreExecFn,
		RunE:    createDeclarationFn,
	}

	createCmd.Flags().StringP("json", "j", "", "json payload to create a declaration")
	createCmd.MarkFlagRequired("json")

	return createCmd
}

func createDeclarationFn(cmd *cobra.Command, declarations []string) error {
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
	err = putJsonReq(ddmUrl.String(), jsonBytes, &resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	return nil
}

// getCmd handles getting declarations on the server
func getDeclarationCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:     "get",
		Short:   fmt.Sprintf("get a declaration"),
		Long:    fmt.Sprintf("get a declaration"),
		PreRunE: applyPreExecFn,
		RunE:    getDeclarationFn,
	}

	getCmd.Flags().StringP("identifier", "i", "", "Identifier of the declaration to retrieve")
	getCmd.MarkFlagRequired("identifier")

	return getCmd
}

func getDeclarationFn(cmd *cobra.Command, declarations []string) error {
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

// getCmd handles getting declarations on the server
func deleteDeclarationCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:     "delete",
		Short:   fmt.Sprintf("delete a declaration"),
		Long:    fmt.Sprintf("delete a declaration"),
		PreRunE: applyPreExecFn,
		RunE:    deleteDeclarationFn,
	}

	deleteCmd.Flags().StringP("identifier", "i", "", "Identifier of the declaration to retrieve")
	deleteCmd.MarkFlagRequired("identifier")

	return deleteCmd
}

func deleteDeclarationFn(cmd *cobra.Command, declarations []string) error {
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
	err = deleteReq(ddmUrl.String(), &resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(body))
	fmt.Println(resp.Status)
	return nil
}

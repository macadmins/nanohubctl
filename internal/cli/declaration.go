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
		getSetsDeclarationCmd(),
	)

	return declarationCmd
}

// createCmd handles creating declarations on the server
func createDeclarationCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:     "create",
		Short:   fmt.Sprintf("Create a declaration"),
		Long:    fmt.Sprintf("Create a declaration"),
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
		Short:   fmt.Sprintf("Get declaration details"),
		Long:    fmt.Sprintf("Get declaration details"),
		PreRunE: applyPreExecFn,
		RunE:    getDeclarationFn,
	}

	getCmd.Flags().StringP("identifier", "i", "", "Identifier of the declaration to retrieve")
	getCmd.Flags().BoolP("all", "a", false, "Identifier of the declaration to retrieve")
	// We can only retrieve a single declaration or all of them
	getCmd.MarkFlagsMutuallyExclusive("identifier", "all")

	return getCmd
}

func getDeclarationFn(cmd *cobra.Command, declarations []string) error {
	identifier, err := cmd.Flags().GetString("identifier")
	if err != nil {
		return err
	}
	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}
	fmt.Printf("Getting declaration for identifier %s\n", identifier)
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if all {
		ddmUrl.Path = path.Join(ddmUrl.Path, "v1/declarations")
	} else {
		ddmUrl.Path = path.Join(ddmUrl.Path, "v1/declarations", identifier)
	}
	if err != nil {
		return err
	}
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
	// Could be an array of strings or a proper dictionary
	var jsonResponse interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return err
	}
	fmt.Println(PrettyJsonPrint(jsonResponse))
	return nil
}

// getCmd handles getting declarations on the server
func getSetsDeclarationCmd() *cobra.Command {
	getSetsCmd := &cobra.Command{
		Use:     "sets",
		Short:   fmt.Sprintf("List sList set membership for the specified declarationet membership for the specified declaration"),
		Long:    fmt.Sprintf("List set membership for a given declaration"),
		PreRunE: applyPreExecFn,
		RunE:    getSetsDeclarationFn,
	}

	getSetsCmd.Flags().StringP("identifier", "i", "", "Identifier of the declaration to retrieve")
	getSetsCmd.MarkFlagRequired("identifier")

	return getSetsCmd
}

func getSetsDeclarationFn(cmd *cobra.Command, declarations []string) error {
	// ToDo(natewalck) - Check to see if identifier exists before getting sets for it
	identifier, err := cmd.Flags().GetString("identifier")
	if err != nil {
		return err
	}
	fmt.Printf("Getting set membership for identifier %s\n", identifier)
	ddmUrl, err := url.Parse(viper.GetString("url"))
	ddmUrl.Path = path.Join(ddmUrl.Path, "/v1/declaration-sets", identifier)
	if err != nil {
		return err
	}
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
	// Could be an array of strings or a proper dictionary
	var jsonResponse []string
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
		Short:   fmt.Sprintf("Delete a declaration"),
		Long:    fmt.Sprintf("Delete a declaration"),
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
	return nil
}

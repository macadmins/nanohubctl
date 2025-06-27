package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"slices"

	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func declarationCmd() *cobra.Command {
	declarationCmd := &cobra.Command{
		Use:     "declaration [command]",
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

// getDeclarationCmd retrieves a declaration from the server
func getDeclarationCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:     "get com.example.declaration",
		Short:   fmt.Sprintf("Get declaration details for identifier"),
		Long:    fmt.Sprintf("Get declaration details for identifier"),
		Args:    cobra.ExactArgs(1),
		PreRunE: applyPreExecFn,
		RunE:    getDeclarationFn,
	}

	return getCmd
}

func getDeclarationFn(cmd *cobra.Command, args []string) error {
	identifier := args[0]

	fmt.Printf("Getting declaration for identifier %s\n", identifier)
	ddmUrl, err := url.Parse(viper.GetString("url"))
	ddmUrl.Path = path.Join(ddmUrl.Path, "declarations", identifier)
	if err != nil {
		return err
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
	// Could be an array of strings or a proper dictionary
	var jsonResponse interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return err
	}
	fmt.Println(PrettyJsonPrint(jsonResponse))
	return nil
}

// getSetsDeclarationCmd Lists set membership for a given declaration
func getSetsDeclarationCmd() *cobra.Command {
	getSetsCmd := &cobra.Command{
		Use:     "sets com.example.declaration",
		Short:   fmt.Sprintf("List set membership for a given declaration"),
		Long:    fmt.Sprintf("List set membership for a given declaration"),
		Args:    cobra.ExactArgs(1),
		PreRunE: applyPreExecFn,
		RunE:    getSetsDeclarationFn,
	}

	return getSetsCmd
}

func getSetsDeclarationFn(cmd *cobra.Command, args []string) error {
	identifier := args[0]

	ddmGetDeclsUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}
	ddmGetDeclsUrl.Path = path.Join(ddmGetDeclsUrl.Path, "declarations")
	allDecls, nil := getAllDeclarations(ddmGetDeclsUrl)
	if !slices.Contains(allDecls, identifier) {
		return fmt.Errorf("%s is not a valid declaration", identifier)
	}

	fmt.Printf("Getting set membership for identifier %s\n", identifier)
	ddmUrl, err := url.Parse(viper.GetString("url"))
	ddmUrl.Path = path.Join(ddmUrl.Path, "/declaration-sets", identifier)
	if err != nil {
		return err
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
	// Could be an array of strings or a proper dictionary
	var jsonResponse []string
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return err
	}
	fmt.Println(PrettyJsonPrint(jsonResponse))
	return nil
}

// createDeclarationCmd creates a new declaration based on a JSON file on disk
func createDeclarationCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:     "create /path/to/declaration.json",
		Short:   fmt.Sprintf("Create a declaration"),
		Long:    fmt.Sprintf("Create a declaration"),
		Args:    cobra.ExactArgs(1),
		PreRunE: applyPreExecFn,
		RunE:    createDeclarationFn,
	}

	return createCmd
}

func createDeclarationFn(cmd *cobra.Command, args []string) error {
	jsonPath := args[0]
	jsonBytes, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	fmt.Printf("Creating declaration using %s\n", jsonPath)
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}
	ddmUrl.Path = path.Join(ddmUrl.Path, "declarations")
	var resp *http.Response
	err = putJsonReq(ddmUrl.String(), jsonBytes, &resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	return nil
}

// deleteDeclarationCmd deletes a declaration from the server
func deleteDeclarationCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:     "delete com.example.declaration",
		Short:   fmt.Sprintf("Delete a declaration"),
		Long:    fmt.Sprintf("Delete a declaration"),
		Args:    cobra.ExactArgs(1),
		PreRunE: applyPreExecFn,
		RunE:    deleteDeclarationFn,
	}

	return deleteCmd
}

func deleteDeclarationFn(cmd *cobra.Command, args []string) error {
	identifier := args[0]
	fmt.Printf("Getting declaration for identifier %s\n", identifier)
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}
	ddmUrl.Path = path.Join(ddmUrl.Path, "declarations", identifier)
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

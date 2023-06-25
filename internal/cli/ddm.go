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

func ddmCmd() *cobra.Command {
	ddmCmd := &cobra.Command{
		Use:     "ddm",
		Short:   fmt.Sprintf("This verb handles all ddm endpoint related operations"),
		Long:    fmt.Sprintf("This verb handles all ddm endpoint related operations"),
		PreRunE: applyPreExecFn,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	ddmCmd.PersistentFlags().StringP("ID", "i", "", "Identifier of the client to return ddm for")
	ddmCmd.MarkPersistentFlagRequired("ID")
	ddmCmd.AddCommand(
		tokenDdmCmd(),
		declarationsDdmCmd(),
		getDeclarationDdmCmd(),
	)

	return ddmCmd
}

// declarationddmCmd lists declarations for a specified device ID
func tokenDdmCmd() *cobra.Command {
	tokenDdmCmd := &cobra.Command{
		Use:     "token",
		Short:   fmt.Sprintf("Show DDM token for a given device ID"),
		Long:    fmt.Sprintf("Show DDM token for a given device ID"),
		PreRunE: applyPreExecFn,
		RunE:    ddmFn,
	}

	return tokenDdmCmd
}

// errorsCmd Lists errors for a specified device ID
func declarationsDdmCmd() *cobra.Command {
	declarationsDdmCmd := &cobra.Command{
		Use:     "declarations",
		Short:   fmt.Sprintf("Show declaration items from the ddm endpoint"),
		Long:    fmt.Sprintf("Show declaration items from the ddm endpoint"),
		PreRunE: applyPreExecFn,
		RunE:    ddmFn,
	}

	return declarationsDdmCmd
}

// valuesCmd lists all values for a specified device ID
func getDeclarationDdmCmd() *cobra.Command {
	getDeclarationDdmCmd := &cobra.Command{
		Use:     "declaration",
		Short:   fmt.Sprintf("Get a specific declaration type about a given device ID from the ddm endpoint"),
		Long:    fmt.Sprintf("Get a specific declaration type about a given device ID from the ddm endpoint"),
		PreRunE: applyPreExecFn,
		RunE:    getDeclarationDdmFn,
	}

	getDeclarationDdmCmd.Flags().StringP("type", "t", "", "Type of the declaration to retrieve. (Configuration, Management, Assets, etc)")
	getDeclarationDdmCmd.Flags().StringP("identifier", "d", "", "Identifier of the declaration to retrieve.")
	getDeclarationDdmCmd.MarkFlagRequired("type")
	getDeclarationDdmCmd.MarkFlagRequired("identifier")
	getDeclarationDdmCmd.MarkFlagsRequiredTogether("type", "identifier")

	return getDeclarationDdmCmd
}

// ddmFn handles all logic for the various ddm commands
func getDeclarationDdmFn(cmd *cobra.Command, ddms []string) error {
	deviceID, err := cmd.Flags().GetString("ID")
	if err != nil {
		return err
	}
	declarationType, err := cmd.Flags().GetString("type")
	if err != nil {
		return err
	}
	declaration, err := cmd.Flags().GetString("identifier")
	if err != nil {
		return err
	}
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}
	ddmUrl.Path = path.Join(ddmUrl.Path, "declaration", declarationType, declaration)
	var resp *http.Response
	err = getReqWithEnrollmentID(ddmUrl.String(), deviceID, &resp)
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

// ddmFn handles all logic for the various ddm commands
func ddmFn(cmd *cobra.Command, ddms []string) error {
	deviceID, err := cmd.Flags().GetString("ID")
	if err != nil {
		return err
	}
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}
	switch cmd.Use {
	case "token":
		ddmUrl.Path = path.Join(ddmUrl.Path, "tokens")
	case "declarations":
		ddmUrl.Path = path.Join(ddmUrl.Path, "declaration-items")
	case "errors":
		ddmUrl.Path = path.Join(ddmUrl.Path, "v1/ddm-errors", deviceID)
	default:
		return fmt.Errorf("%s is not a valid ddm type", cmd.Use)
	}
	var resp *http.Response
	err = getReqWithEnrollmentID(ddmUrl.String(), deviceID, &resp)
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

package ddm

import (
	"encoding/json"
	"fmt"
	"io"
	"path"

	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/macadmins/nanohubctl/internal/utils"
)

func ddmCmd() *cobra.Command {
	ddmCmd := &cobra.Command{
		Use:     "ddm",
		Short:   "ddm endpoint related operations",
		Long:    "ddm endpoint related operations",
		PreRunE: utils.ApplyPreExecFn,
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
	// ToDo - Fix these commands
	// tokenDdmCmd(),
	// declarationsDdmCmd(),
	// getDeclarationDdmCmd(),
	)

	return ddmCmd
}

// tokenDdmCmd shows the ddm token for a given device ID
func tokenDdmCmd() *cobra.Command {
	tokenDdmCmd := &cobra.Command{
		Use:     "token",
		Short:   "Show DDM token for a given device ID",
		Long:    "Show DDM token for a given device ID",
		PreRunE: utils.ApplyPreExecFn,
		RunE:    ddmFn,
	}

	return tokenDdmCmd
}

// declarationsDdmCmd shows all declaration items for a given device ID from the ddm endpoint
func declarationsDdmCmd() *cobra.Command {
	declarationsDdmCmd := &cobra.Command{
		Use:     "declaration-items",
		Short:   "Show declaration items from the ddm endpoint",
		Long:    "Show declaration items from the ddm endpoint",
		PreRunE: utils.ApplyPreExecFn,
		RunE:    ddmFn,
	}

	return declarationsDdmCmd
}

// getDeclarationDdmCmd gets a specific declaration type about a given device ID from the ddm endpoint
func getDeclarationDdmCmd() *cobra.Command {
	getDeclarationDdmCmd := &cobra.Command{
		Use:     "declaration-details",
		Short:   "Get a specific declaration type about a given device ID from the ddm endpoint",
		Long:    "Get a specific declaration type about agiven device ID from the ddm endpoint",
		PreRunE: utils.ApplyPreExecFn,
		RunE:    getDeclarationDdmFn,
	}

	getDeclarationDdmCmd.Flags().StringP("type", "t", "", "Type of the declaration to retrieve. (Configuration, Management, Assets, etc)")
	getDeclarationDdmCmd.Flags().StringP("identifier", "d", "", "Identifier of the declaration to retrieve.")
	getDeclarationDdmCmd.MarkFlagRequired("type")
	getDeclarationDdmCmd.MarkFlagRequired("identifier")
	getDeclarationDdmCmd.MarkFlagsRequiredTogether("type", "identifier")

	return getDeclarationDdmCmd
}

// getDeclarationDdmFn handles the declaration ddm endpoint
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
	ddmUrl, err := utils.GetDDMUrl()
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
	fmt.Println(utils.PrettyJsonPrint(jsonResponse))
	return nil
}

// ddmFn handles all logic for the various ddm commands
func ddmFn(cmd *cobra.Command, ddms []string) error {
	deviceID := viper.GetString("client_id")
	// deviceID, err := cmd.Flags().GetString("ID")
	// if err != nil {
	// 	return err
	// }
	ddmUrl, err := utils.GetDDMUrl()
	if err != nil {
		return err
	}
	fmt.Println(cmd.Name())
	switch cmd.Use {
	case "token":
		ddmUrl.Path = path.Join(ddmUrl.Path, "tokens")
		fmt.Println(ddmUrl.String())
	case "declarations":
		ddmUrl.Path = path.Join(ddmUrl.Path, "declaration-items")
	case "errors":
		ddmUrl.Path = path.Join(ddmUrl.Path, "ddm-errors", deviceID)
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
	fmt.Println(utils.PrettyJsonPrint(jsonResponse))
	return nil
}

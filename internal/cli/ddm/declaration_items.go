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

func declarationItemsCmd() *cobra.Command {
	declarationItemsCmd := &cobra.Command{
		Use:     "declaration-items",
		Short:   "Show all declaration items and ServerTokens for a client from the ddm endpoint",
		Long:    "Show all declaration items and ServerTokens for a client from the ddm endpoint",
		PreRunE: utils.ApplyPreExecFn,
		RunE:    declarationItemsDdmFn,
	}

	return declarationItemsCmd
}

func declarationItemsDdmFn(cmd *cobra.Command, args []string) error {
	deviceID := viper.GetString("client_id")
	ddmUrl, err := utils.GetDDMUrl()
	if err != nil {
		return err
	}
	ddmUrl.Path = path.Join(ddmUrl.Path, "declaration-items")

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

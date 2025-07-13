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

func tokenDdmCmd() *cobra.Command {
	tokenDdmCmd := &cobra.Command{
		Use:     "tokens",
		Short:   "Show DDM sync token for a given device ID",
		Long:    "Show DDM sync token for a given device ID",
		PreRunE: utils.ApplyPreExecFn,
		RunE:    tokensDdmFn,
	}

	return tokenDdmCmd
}

func tokensDdmFn(cmd *cobra.Command, args []string) error {
	deviceID := viper.GetString("client_id")
	ddmUrl, err := utils.GetDDMUrl()
	if err != nil {
		return err
	}
	ddmUrl.Path = path.Join(ddmUrl.Path, "tokens")

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

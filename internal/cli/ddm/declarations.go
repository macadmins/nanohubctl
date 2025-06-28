package ddm

import (
	"encoding/json"
	"fmt"
	"io"
	"path"

	"net/http"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/macadmins/nanohubctl/internal/utils"
)

func declarationsCmd() *cobra.Command {
	declarationsCmd := &cobra.Command{
		Use:     "declarations",
		Short:   fmt.Sprintf("This verb gets a list of all declarations"),
		Long:    fmt.Sprintf("This verb gets a list of all declarations"),
		PreRunE: utils.ApplyPreExecFn,
		RunE: func(cmd *cobra.Command, args []string) error {
			// if err := cmd.Help(); err != nil {
			// 	return err
			// }
			ddmUrl, err := utils.GetDDMUrl()
			if err != nil {
				return err
			}
			ddmUrl.Path = path.Join(ddmUrl.Path, "declarations")
			allDecls, nil := getAllDeclarations(ddmUrl)
			for _, decl := range allDecls {
				fmt.Println(decl)
			}
			return nil
		},
	}

	return declarationsCmd
}

func getAllDeclarations(ddmUrl *url.URL) ([]string, error) {
	var resp *http.Response
	err := getReq(ddmUrl.String(), &resp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	// Could be an array of strings or a proper dictionary
	var jsonResponse []string
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return nil, err
	}
	return jsonResponse, nil
}

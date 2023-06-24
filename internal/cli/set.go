package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"path"

	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func setCmd() *cobra.Command {
	setCmd := &cobra.Command{
		Use:     "set",
		Short:   fmt.Sprintf("This verb handles all set related operations"),
		Long:    fmt.Sprintf("This verb handles all set related operations"),
		PreRunE: applyPreExecFn,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}
	setCmd.AddCommand(
		addSetCmd(),
		getSetCmd(),
		deleteSetCmd(),
	)

	return setCmd
}

// getCmd handles getting sets on the server
func getSetCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:     "get",
		Short:   fmt.Sprintf("get a set"),
		Long:    fmt.Sprintf("get a set"),
		PreRunE: applyPreExecFn,
		RunE:    getSetFn,
	}

	getCmd.Flags().StringP("name", "n", "", "Name of the set to retrieve")
	getCmd.MarkFlagRequired("name")

	return getCmd
}

func getSetFn(cmd *cobra.Command, sets []string) error {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	fmt.Printf("Getting set for identifier %s\n", name)
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}
	ddmUrl.Path = path.Join(ddmUrl.Path, "v1/set-declarations", name)
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
	var jsonResponse []string
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return err
	}
	fmt.Println(PrettyJsonPrint(jsonResponse))
	return nil
}

// createCmd handles creating sets on the server
func addSetCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:     "add",
		Short:   fmt.Sprintf("create a set"),
		Long:    fmt.Sprintf("create a set"),
		PreRunE: applyPreExecFn,
		RunE:    addSetFn,
	}

	createCmd.Flags().StringP("name", "n", "", "Name of the set to add item to")
	createCmd.Flags().StringP("identifier", "i", "", "Identifier of the declaration to add to the set")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("identifier")

	return createCmd
}

func addSetFn(cmd *cobra.Command, sets []string) error {
	name, err := cmd.Flags().GetString("name")
	identifier, err := cmd.Flags().GetString("identifier")
	if err != nil {
		return err
	}
	fmt.Printf("Adding %s to set %s...\n", identifier, name)
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}

	resp, err := addOrDeleteSetItem("add", name, identifier, ddmUrl)

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		fmt.Printf("%s is already in %s", identifier, name)
	} else if resp.StatusCode == http.StatusNoContent {
		fmt.Printf("%s has been added to %s", identifier, name)
	} else {
		fmt.Println(resp.Status)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		return fmt.Errorf(string(body))
	}

	return nil
}

// getCmd handles getting sets on the server
func deleteSetCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:     "delete",
		Short:   fmt.Sprintf("delete a set"),
		Long:    fmt.Sprintf("delete a set"),
		PreRunE: applyPreExecFn,
		RunE:    deleteSetFn,
	}

	deleteCmd.Flags().StringP("name", "n", "", "Name of the set to delete the declaration from")
	deleteCmd.Flags().StringP("identifier", "i", "", "Identifier of the declaration to remove from the set")
	deleteCmd.MarkFlagRequired("name")
	deleteCmd.MarkFlagRequired("identifier")

	return deleteCmd
}

func deleteSetFn(cmd *cobra.Command, sets []string) error {
	name, err := cmd.Flags().GetString("name")
	identifier, err := cmd.Flags().GetString("identifier")
	if err != nil {
		return err
	}
	fmt.Printf("Adding %s to set %s...\n", identifier, name)
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}

	resp, err := addOrDeleteSetItem("delete", name, identifier, ddmUrl)

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		fmt.Printf("%s does not exist in %s", identifier, name)
	} else if resp.StatusCode == http.StatusNoContent {
		fmt.Printf("%s has been removed from %s", identifier, name)
	} else {
		fmt.Println(resp.Status)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		return fmt.Errorf(string(body))
	}

	return nil
}

func addOrDeleteSetItem(action, name, identifier string, ddmUrl *url.URL) (*http.Response, error) {
	// Set the path
	ddmUrl.Path = path.Join(ddmUrl.Path, "/v1/set-declarations/", name)
	// Add the query arguments
	q := ddmUrl.Query()
	q.Set("declaration", identifier)
	ddmUrl.RawQuery = q.Encode()
	var resp *http.Response
	var respErr error
	if action == "add" {
		respErr = putReq(ddmUrl.String(), &resp)
	} else if action == "delete" {
		respErr = deleteReq(ddmUrl.String(), &resp)
	}
	if respErr != nil {
		return nil, respErr
	}
	return resp, nil
}

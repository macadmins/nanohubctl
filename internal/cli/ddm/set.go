package ddm

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"path"

	"net/http"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/macadmins/nanohubctl/internal/utils"
)

// setCmd handles creation and management of declaration sets
func setCmd() *cobra.Command {
	setCmd := &cobra.Command{
		Use:     "set",
		Short:   fmt.Sprintf("This verb handles all set related operations"),
		Long:    fmt.Sprintf("This verb handles all set related operations"),
		PreRunE: utils.ApplyPreExecFn,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}
	setCmd.AddCommand(
		listSetsCmd(),
		addSetCmd(),
		getSetCmd(),
		deleteSetCmd(),
	)

	return setCmd
}

// listSetsCmd handles getting sets on the server
func listSetsCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "list",
		Short:   fmt.Sprintf("list all sets"),
		Long:    fmt.Sprintf("list all sets"),
		PreRunE: utils.ApplyPreExecFn,
		RunE:    listSetsFn,
	}

	return listCmd
}

func listSetsFn(cmd *cobra.Command, args []string) error {
	fmt.Printf("Listing all available sets\n")
	ddmUrl, err := utils.GetDDMUrl()
	if err != nil {
		return err
	}
	ddmUrl.Path = path.Join(ddmUrl.Path, "sets")
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
	var jsonResponse []string
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return err
	}
	fmt.Println(PrettyJsonPrint(jsonResponse))
	return nil
}

// getCmd handles getting sets on the server
func getSetCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:     "get [set name]",
		Short:   fmt.Sprintf("Get the declarations for a set"),
		Long:    fmt.Sprintf("Get the declarations for a set"),
		Args:    cobra.MinimumNArgs(1),
		PreRunE: utils.ApplyPreExecFn,
		RunE:    getSetFn,
	}

	return getCmd
}

func getSetFn(cmd *cobra.Command, args []string) error {
	name := args[0]
	fmt.Printf("Getting set for identifier %s\n\n", name)
	ddmUrl, err := utils.GetDDMUrl()
	if err != nil {
		return err
	}
	ddmUrl.Path = path.Join(ddmUrl.Path, "set-declarations", name)
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
	var jsonResponse []string
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return err
	}
	if jsonResponse == nil {
		fmt.Println("No declarations found")
		return nil
	}
	fmt.Println(PrettyJsonPrint(jsonResponse))
	return nil
}

// addSetCmd adds a declaration to a given set
func addSetCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:     "add",
		Short:   fmt.Sprintf("Add a declaration to a set"),
		Long:    fmt.Sprintf("Add a declaration to a set"),
		PreRunE: utils.ApplyPreExecFn,
		RunE:    addSetFn,
	}

	createCmd.Flags().StringP("name", "n", "", "Name of the set to add item to")
	createCmd.Flags().StringP("identifier", "i", "", "Identifier of the declaration to add to the set")
	createCmd.MarkFlagsRequiredTogether("name", "identifier")

	return createCmd
}

func addSetFn(cmd *cobra.Command, sets []string) error {
	name, err := cmd.Flags().GetString("name")
	identifier, err := cmd.Flags().GetString("identifier")
	if err != nil {
		return err
	}
	fmt.Printf("Adding %s to set %s...\n\n", identifier, name)
	ddmUrl, err := utils.GetDDMUrl()
	if err != nil {
		return err
	}

	resp, err := addOrDeleteSetItem("add", name, identifier, ddmUrl)

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		fmt.Printf("%s is already in %s", identifier, name)
	} else if resp.StatusCode == http.StatusNoContent {
		fmt.Printf("%s has been added to set: %s", identifier, name)
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

// deleteSetCmd deletes a declaration from a given set
func deleteSetCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:     "delete",
		Short:   fmt.Sprintf("Delete a declaration from a set"),
		Long:    fmt.Sprintf("Delete a declaration from a set"),
		PreRunE: utils.ApplyPreExecFn,
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
	fmt.Printf("Adding %s to set %s...\n\n", identifier, name)
	ddmUrl, err := utils.GetDDMUrl()
	if err != nil {
		return err
	}

	resp, err := addOrDeleteSetItem("delete", name, identifier, ddmUrl)

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		fmt.Printf("%s does not exist in %s", identifier, name)
	} else if resp.StatusCode == http.StatusNoContent {
		fmt.Printf("%s has been removed from set: %s", identifier, name)
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

// addOrDeleteSetItem handles http for add and remove, probably better to just duplicate the code.
func addOrDeleteSetItem(action, name, identifier string, ddmUrl *url.URL) (*http.Response, error) {
	// Set the path
	ddmUrl.Path = path.Join(ddmUrl.Path, "/set-declarations/", name)
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

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

// deviceCmd handles creation and management of declaration devices
func deviceCmd() *cobra.Command {
	deviceCmd := &cobra.Command{
		Use:     "device",
		Short:   fmt.Sprintf("This verb handles all device related operations"),
		Long:    fmt.Sprintf("This verb handles all device related operations"),
		PreRunE: applyPreExecFn,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}
	deviceCmd.PersistentFlags().StringP("ID", "i", "", "Identifier of the client to return status for")
	deviceCmd.MarkPersistentFlagRequired("ID")

	deviceCmd.AddCommand(
		applyDeviceCmd(),
		getDeviceCmd(),
		removeDeviceCmd(),
	)

	return deviceCmd
}

// getCmd handles getting devices on the server
func getDeviceCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:     "sets",
		Short:   fmt.Sprintf("Get a list of all sets a device has applied to it"),
		Long:    fmt.Sprintf("Get a list of all sets a device has applied to it"),
		PreRunE: applyPreExecFn,
		RunE:    getdeviceFn,
	}

	// getCmd.Flags().StringP("ID", "i", "", "Name of the device to retrieve")
	// getCmd.MarkFlagRequired("ID ")

	return getCmd
}

func getdeviceFn(cmd *cobra.Command, devices []string) error {
	deviceID, err := cmd.Flags().GetString("ID")
	if err != nil {
		return err
	}
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}
	ddmUrl.Path = path.Join(ddmUrl.Path, "v1/enrollment-sets", deviceID)
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
	var jsonResponse interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return err
	}
	fmt.Println(PrettyJsonPrint(jsonResponse))
	return nil
}

// createCmd handles creating devices on the server
func applyDeviceCmd() *cobra.Command {
	applyDeviceCmd := &cobra.Command{
		Use:     "apply",
		Short:   fmt.Sprintf("Apply a declaration set to a device"),
		Long:    fmt.Sprintf("Apply a declaration set to a device"),
		PreRunE: applyPreExecFn,
		RunE:    applyDeviceFn,
	}

	applyDeviceCmd.Flags().StringP("set", "s", "", "Name of the set to apply to the device")
	applyDeviceCmd.Flags().StringP("ID", "i", "", "Enrollment ID of the device")
	applyDeviceCmd.MarkFlagsRequiredTogether("set", "ID")

	return applyDeviceCmd
}

func applyDeviceFn(cmd *cobra.Command, devices []string) error {
	set, err := cmd.Flags().GetString("set")
	deviceID, err := cmd.Flags().GetString("ID")
	if err != nil {
		return err
	}
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}

	resp, err := addOrDeletedeviceItem("apply", deviceID, set, ddmUrl)

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		fmt.Printf("%s is already in %s", deviceID, set)
	} else if resp.StatusCode == http.StatusNoContent {
		fmt.Printf("%s has been added to %s", deviceID, set)
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

// getCmd handles getting devices on the server
func removeDeviceCmd() *cobra.Command {
	removeDeviceCmd := &cobra.Command{
		Use:     "remove",
		Short:   fmt.Sprintf("Remove a declaration set from a device"),
		Long:    fmt.Sprintf("Remove a declaration set from a device"),
		PreRunE: applyPreExecFn,
		RunE:    removeDeviceFn,
	}

	removeDeviceCmd.Flags().StringP("set", "s", "", "Name of the set to apply to the device")
	removeDeviceCmd.Flags().StringP("ID", "i", "", "Enrollment ID of the device")
	removeDeviceCmd.MarkFlagsRequiredTogether("set", "ID")

	return removeDeviceCmd
}

func removeDeviceFn(cmd *cobra.Command, devices []string) error {
	deviceID, err := cmd.Flags().GetString("ID")
	set, err := cmd.Flags().GetString("set")
	if err != nil {
		return err
	}
	fmt.Printf("Adding %s to device %s...\n", deviceID, set)
	ddmUrl, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}

	resp, err := addOrDeletedeviceItem("remove", deviceID, set, ddmUrl)

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		fmt.Printf("%s does not exist in %s", deviceID, set)
	} else if resp.StatusCode == http.StatusNoContent {
		fmt.Printf("%s has been removed from %s", deviceID, set)
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

func addOrDeletedeviceItem(action, deviceID, set string, ddmUrl *url.URL) (*http.Response, error) {
	// Device path for the enrollment set
	ddmUrl.Path = path.Join(ddmUrl.Path, "v1/enrollment-sets", deviceID)
	// Add the query arguments
	q := ddmUrl.Query()
	q.Set("set", set)
	ddmUrl.RawQuery = q.Encode()
	fmt.Println(ddmUrl.String())
	var resp *http.Response
	var respErr error
	if action == "apply" {
		respErr = putReq(ddmUrl.String(), &resp)
	} else if action == "remove" {
		respErr = deleteReq(ddmUrl.String(), &resp)
	}
	if respErr != nil {
		return nil, respErr
	}
	return resp, nil
}

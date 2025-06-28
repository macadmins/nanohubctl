package ddm

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func applyPreExecFn(cmd *cobra.Command, args []string) error {
	// Bind all the flags to a viper setting so we can use viper everywhere without thinking about it
	if err := viper.BindPFlag("url", cmd.Flags().Lookup("url")); err != nil {
		return errors.New("failed to bind url to viper")
	}
	if err := viper.BindPFlag("api_key", cmd.Flags().Lookup("api_key")); err != nil {
		return errors.New("failed to bind key to viper")
	}
	if err := viper.BindPFlag("client_id", cmd.Flags().Lookup("client_id")); err != nil {
		return errors.New("failed to bind id to viper")
	}
	// If declaration or declarations is called, skip the UUID check
	if !(cmd.Name() == "declarations" || cmd.Name() == "declaration" || cmd.Parent().Name() == "declaration") {
		clientUUID := viper.GetString("client_id")
		if !validUUID(clientUUID) {
			return errors.New("Invalid UUID provided")
		}
	}
	// Make sure mandatory values are present before continuing
	if viper.GetString("URL") == "" {
		return errors.New("Base Url must be provided!")
	}
	if viper.GetString("API_KEY") == "" {
		return errors.New("API Key must be provided!")
	}

	return nil
}

func putReq(url string, resp **http.Response) error {
	req, err := http.NewRequest("PUT", url, nil)
	req.ContentLength = 0
	auth := viper.GetString("api_user") + ":" + viper.GetString("api_key")
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encodedAuth)

	*resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func putJsonReq(url string, jsonBytes []byte, resp **http.Response) error {
	body := bytes.NewBuffer(jsonBytes)
	req, err := http.NewRequest("PUT", url, body)
	auth := viper.GetString("api_user") + ":" + viper.GetString("api_key")
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encodedAuth)

	*resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func getReq(url string, resp **http.Response) error {
	req, err := http.NewRequest("GET", url, nil)
	auth := viper.GetString("api_user") + ":" + viper.GetString("api_key")
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encodedAuth)

	*resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func getReqWithEnrollmentID(url, deviceID string, resp **http.Response) error {
	req, err := http.NewRequest("GET", url, nil)
	auth := viper.GetString("api_user") + ":" + viper.GetString("api_key")
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encodedAuth)
	req.Header.Add("X-Enrollment-ID", deviceID)

	*resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func deleteReq(url string, resp **http.Response) error {
	req, err := http.NewRequest("DELETE", url, nil)
	auth := viper.GetString("api_user") + ":" + viper.GetString("api_key")
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encodedAuth)

	*resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func PrettyJsonPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func validUUID(uuid string) bool {
	if len(uuid) == 36 || len(uuid) == 25 {
		return true
	}
	return false
}

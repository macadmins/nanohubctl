package cli

import (
	"bytes"
	"encoding/base64"
	"errors"
	"log"
	"net/http"

	"github.com/google/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func setLoggerOpts() {
	if vv {
		logger.SetLevel(2)
	}
	logger.SetFlags(log.LUTC)
}

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

	// Make sure mandatory values are present before continuing
	if viper.GetString("URL") == "" {
		return errors.New("Base Url must be provided!")
	}
	if viper.GetString("API_KEY") == "" {
		return errors.New("API Key must be provided!")
	}

	return nil
}

func putReq(url string, jsonBytes []byte, resp **http.Response) error {
	username := "kmfddm"

	body := bytes.NewBuffer(jsonBytes)
	req, err := http.NewRequest("PUT", url, body)
	auth := username + ":" + viper.GetString("api_key")
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encodedAuth)

	*resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}

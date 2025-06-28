package nanocmd

import (
	"encoding/base64"
	"net/http"

	"github.com/spf13/viper"
)

func postReq(url string, resp **http.Response) error {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
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

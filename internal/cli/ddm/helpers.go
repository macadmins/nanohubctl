package ddm

import (
	"bytes"
	"encoding/base64"
	"net/http"

	"github.com/spf13/viper"
)

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

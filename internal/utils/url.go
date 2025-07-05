package utils

import (
	"net/url"
	"path"

	"github.com/spf13/viper"
)

// GetDDMUrl returns the base nanohub URL with the proper DDM API path
func GetDDMUrl() (url.URL, error) {
	baseUrl := viper.GetString("url")
	ddmUrl, err := url.Parse(baseUrl)
	if err != nil {
		return url.URL{}, err
	}
	ddmUrl.Path = path.Join(ddmUrl.Path, "api/v1/ddm")
	return *ddmUrl, nil
}

// GetNanoCmdUrl returns the base nanohub URL with the proper NanoCMD API path
func GetNanoCMDUrl() (*url.URL, error) {
	baseUrl := viper.GetString("url")
	ddmUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	ddmUrl.Path = path.Join(ddmUrl.Path, "api/v1/nanocmd")
	return ddmUrl, nil
}

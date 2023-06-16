package cli

import (
	"errors"
	"log"

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
		return errors.New("Url must be provided!")
	}
	if viper.GetString("API_KEY") == "" {
		return errors.New("Key must be provided!")
	}

	return nil
}

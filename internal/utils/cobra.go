package utils

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func validUUID(uuid string) bool {
	if len(uuid) == 36 || len(uuid) == 25 {
		return true
	}
	return false
}

func ApplyPreExecFn(cmd *cobra.Command, args []string) error {
	// Bind all the flags to a viper setting so we can use viper everywhere without thinking about it
	if err := viper.BindPFlag("url", cmd.Flags().Lookup("url")); err != nil {
		return errors.New("failed to bind url to viper")
	}
	if err := viper.BindPFlag("api_key", cmd.Flags().Lookup("api_key")); err != nil {
		return errors.New("failed to bind key to viper")
	}
	if err := viper.BindPFlag("api_user", cmd.Flags().Lookup("api_user")); err != nil {
		return errors.New("failed to bind api_user to viper")
	}
	if err := viper.BindPFlag("client_id", cmd.Flags().Lookup("client_id")); err != nil {
		return errors.New("failed to bind id to viper")
	}

	// For DDM commands, check UUID validity (skip for certain commands)
	if cmd.Parent() != nil && cmd.Parent().Name() == "ddm" {
		if !(cmd.Name() == "declarations" || cmd.Name() == "declaration" || cmd.Parent().Name() == "declaration") {
			clientUUID := viper.GetString("client_id")
			if !validUUID(clientUUID) {
				return errors.New("Invalid UUID provided")
			}
		}
	}

	// Make sure mandatory values are present before continuing
	if viper.GetString("url") == "" {
		return errors.New("Base URL must be provided!")
	}
	if viper.GetString("api_key") == "" {
		return errors.New("API Key must be provided!")
	}

	return nil
}

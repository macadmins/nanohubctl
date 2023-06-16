package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func applyCmd() *cobra.Command {
	applyCmd := &cobra.Command{
		Use:     "apply",
		Short:   fmt.Sprintf("apply a declaration"),
		Long:    fmt.Sprintf("apply a declaration"),
		Args:    cobra.MinimumNArgs(1),
		PreRunE: applyPreExecFn,
		RunE:    applyFn,
	}
	// applyCmd.Flags().StringP("url", "u", viper.GetString("BASE_URL"), "URL of the ddm instance")
	// applyCmd.Flags().StringP("key", "k", viper.GetString("API_KEY"), "API for for accessing the ddm instanuce")
	// applyCmd.Flags().StringP("id", "i", viper.GetString("ID"), "ID of the machine to configure")

	return applyCmd
}

func applyFn(cmd *cobra.Command, declarations []string) error {
	fmt.Println("Applying a declaration")
	return nil
}

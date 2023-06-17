package cli

import (
	"context"
	"fmt"

	"github.com/google/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// These vars are available to every sub command
	debug bool
	vv    bool

	version string = "0.0.1"
)

func ExecuteWithContext(ctx context.Context) error {
	return rootCmd().ExecuteContext(ctx)
}

func rootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   fmt.Sprintf("ddmctl"),
		Short: "A command line tool for working with ddm",
		Long:  "A command line tool for working with declarative device management",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setLoggerOpts()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
		Args:    cobra.NoArgs,
		Version: version,
	}

	// At the rootCmd level, set these global flags that will be available to downstream cmds
	rootCmd.PersistentFlags().String("url", "", "URL of the ddm instance")
	rootCmd.PersistentFlags().String("api_key", "", "API key for the ddm instance")
	rootCmd.PersistentFlags().String("client_id", "", "Client ID to apply items to")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Run in debug mode")
	rootCmd.PersistentFlags().BoolVar(&vv, "vv", false, "Run in verbose logging mode")

	if vv {
		logger.SetLevel(2)
	}

	// Set up ENV namespace and ENV vars
	viper.SetEnvPrefix("DDM")
	viper.BindEnv("URL")
	viper.BindEnv("API_KEY")
	viper.BindEnv("CLIENT_ID")

	// Import subCmds into the rootCmd
	rootCmd.AddCommand(
		declarationCmd(),
		// getCmd(),
		// deleteCmd(),
	)

	return rootCmd
}

package cli

import (
	"context"
	"log"

	"github.com/google/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/macadmins/nanohubctl/internal/cli/ddm"
	"github.com/macadmins/nanohubctl/internal/cli/godeclr"
	"github.com/macadmins/nanohubctl/internal/cli/nanocmd"
)

var (
	// These vars are available to every sub command
	debug bool
	vv    bool

	version string = "1.0.5"
)

func setLoggerOpts() {
	if vv {
		logger.SetLevel(2)
	}
	logger.SetFlags(log.LUTC)
}

func ExecuteWithContext(ctx context.Context) error {
	return rootCmd().ExecuteContext(ctx)
}

func rootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "nanohubctl [ subcommand ]",
		Short: "A command line tool for working with nanohub",
		Long:  "A command line tool for working with nanohub APIs",
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
	rootCmd.PersistentFlags().String("api_user", "nanohub", "API key for the ddm instance")
	rootCmd.PersistentFlags().String("client_id", "", "Client ID to apply items to")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Run in debug mode")
	rootCmd.PersistentFlags().BoolVar(&vv, "vv", false, "Run in verbose logging mode")
	if vv {
		logger.SetLevel(2)
	}

	// Viper has an order of precedence for settings:
	// 1. PFlags
	// 2. ENV vars
	// 3. If both Flag and Env var are set, flag wins
	// 4. Defaults

	// Bind PFlags to viper settings
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api_key"))
	viper.BindPFlag("api_user", rootCmd.PersistentFlags().Lookup("api_user"))
	viper.BindPFlag("client_id", rootCmd.PersistentFlags().Lookup("client_id"))

	// Set up ENV namespace and ENV vars
	// All env vars will be prefixed with DDM
	viper.SetEnvPrefix("NANOHUB")
	viper.BindEnv("URL")
	viper.BindEnv("API_KEY")
	viper.BindEnv("API_USER")
	viper.BindEnv("CLIENT_ID")

	// Set defaults
	viper.SetDefault("api_user", "nanohub")

	// Import subCmds into the rootCmd
	rootCmd.AddCommand(
		ddm.RootCmd(),
		nanocmd.RootCmd(),
		godeclr.RootCmd(),
	)

	return rootCmd
}

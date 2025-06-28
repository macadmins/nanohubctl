package ddm

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	ddmRootCmd := &cobra.Command{
		Use:   "ddm",
		Short: "ddm operations on nanohub",
		Long:  "All declarative device management operations",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setLoggerOpts()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	// Set persistent flags for debug and verbose logging
	ddmRootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Run in debug mode")
	ddmRootCmd.PersistentFlags().BoolVar(&vv, "vv", false, "Run in verbose logging mode")

	ddmRootCmd.AddCommand(
		declarationsCmd(),
		declarationCmd(),
		setCmd(),
		deviceCmd(),
		ddmCmd(),
	)

	return ddmRootCmd
}

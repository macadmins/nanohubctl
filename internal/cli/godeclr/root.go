package godeclr

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	goDeclrRootCmd := &cobra.Command{
		Use:   "godeclr -id [IDENTIFIER] -token [TOKEN] -type [TYPE]",
		Short: "godeclr",
		Long:  "Generates DDM Declarations based on the device-management schemas",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	goDeclrRootCmd.AddCommand(
		TypesCmd(),
		TypeCmd(),
	)

	return goDeclrRootCmd
}

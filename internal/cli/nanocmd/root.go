package nanocmd

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	nanocmdRootCmd := &cobra.Command{
		Use:   "nanocmd",
		Short: "nanocmd operations on nanohub",
		Long:  "All declarative device management operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	nanocmdRootCmd.AddCommand(
		WorkflowCmd(),
	)

	return nanocmdRootCmd
}

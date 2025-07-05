package ddm

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	ddmRootCmd := &cobra.Command{
		Use:   "ddm",
		Short: "ddm operations on nanohub",
		Long:  "All declarative device management operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	ddmRootCmd.AddCommand(
		declarationsCmd(),
		declarationCmd(),
		setCmd(),
		deviceCmd(),
		syncCmd(),
		// ToDo - Fix ddmCmd stuff
		// Make it more clear how these differ from the other commands
		// tokenDdmCmd(),
		// declarationsDdmCmd(),
		// getDeclarationDdmCmd(),
	)

	return ddmRootCmd
}

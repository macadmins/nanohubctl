// package cli

// import (
// 	"fmt"

// 	"github.com/spf13/cobra"
// )

// func statusCmd() *cobra.Command {
// 	statusCmd := &cobra.Command{
// 		Use:     "status",
// 		Short:   fmt.Sprintf("This verb handles all status related operations"),
// 		Long:    fmt.Sprintf("This verb handles all status related operations"),
// 		PreRunE: applyPreExecFn,
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			if err := cmd.Help(); err != nil {
// 				return err
// 			}
// 			return nil
// 		},
// 	}

// 	// statusCmd.PersistentFlags().StringP("ID", "i", "", "Identifier of the client to return status for")
// 	// statusCmd.MarkPersistentFlagRequired("ID")
// 	statusCmd.AddCommand(
// 		declarationStatusCmd(),
// 		errorsCmd(),
// 		valuesCmd(),
// 	)

// 	return statusCmd
// }

// // declarationStatusCmd lists declarations for a specified device ID

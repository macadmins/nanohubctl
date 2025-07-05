package godeclr

import (
	"fmt"
	"os"
	"sort"

	"github.com/korylprince/go-adm/declarations"
	"github.com/spf13/cobra"
)

func TypesCmd() *cobra.Command {
	typesCmd := &cobra.Command{
		Use:   "types",
		Short: "types",
		Long:  "List all supported declaration types.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var typs []string
			for typ := range declarations.DeclarationMap {
				typs = append(typs, typ)
			}
			sort.Strings(typs)
			fmt.Println("Supported declaration types:")
			for _, typ := range typs {
				fmt.Println("\t" + typ)
			}
			os.Exit(0)
			return nil
		},
	}

	return typesCmd
}

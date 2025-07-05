package godeclr

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/korylprince/go-adm/declarations"
	"github.com/korylprince/go-adm/tagutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TypeCmd() *cobra.Command {
	typeCmd := &cobra.Command{
		Use:   "type [declaration type] [-full]",
		Short: "type",
		Long:  "Show an example of the specified declaration type",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverToken := ""
			if len(args) == 0 {
				return fmt.Errorf("declaration type must be specified")
			}
			typ := args[0]
			_, ok := declarations.DeclarationMap[typ]
			if !ok {
				return fmt.Errorf("unknown declaration type: %s", typ)
			}
			decl_identifier := viper.GetString("decl_identifier")
			if decl_identifier == "" {
				decl_identifier = uuid.New().String()
			}
			decl, err := declarations.NewFromType(typ, decl_identifier, serverToken)
			if err != nil {
				fmt.Println("could not generate declaration:", err)
				os.Exit(1)
			}
			var declobj any = decl
			if viper.GetBool("full") {
				payload := tagutil.FullFields(decl.Payload)
				if err = tagutil.SetDefaults(payload); err != nil {
					fmt.Println("could not fill out declaration:", err)
					os.Exit(1)
				}

				m := map[string]any{

					"Type":       decl.Type(),
					"Identifier": decl.Identifier,
					"Payload":    payload,
				}
				if decl.ServerToken != "" {
					m["ServerToken"] = decl.ServerToken
				}
				declobj = m
			}
			buf, err := json.MarshalIndent(declobj, "", "\t")
			if err != nil {
				fmt.Println("could not json marshal declaration:", err)
				os.Exit(1)
			}

			fmt.Println(string(buf))
			return nil
		},
	}

	typeCmd.Flags().StringP("type", "T", "", "declaration type. Use -types to list all supported types")
	typeCmd.Flags().BoolP("full", "f", false, "output all fields in the declaration")
	typeCmd.PersistentFlags().StringP("identifier", "I", "", "declaration identifier (auto-generated UUID if not specified)")

	viper.BindPFlag("type", typeCmd.Flags().Lookup("type"))
	viper.BindPFlag("full", typeCmd.Flags().Lookup("full"))
	viper.BindPFlag("decl_identifier", typeCmd.PersistentFlags().Lookup("identifier"))

	return typeCmd
}

package ddm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/macadmins/nanohubctl/internal/utils"
	"github.com/spf13/cobra"
)

func syncCmd() *cobra.Command {
	syncDirCmd := &cobra.Command{
		Use:     "sync /path/to/directory",
		Short:   "Sync directory with DDM",
		Long:    "Sync directory with DDM",
		Args:    cobra.ExactArgs(1),
		PreRunE: utils.ApplyPreExecFn,
		RunE:    syncDirFn,
	}
	return syncDirCmd
}

func syncDirFn(cmd *cobra.Command, args []string) error {
	dirPath := args[0]

	// Check if directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist", dirPath)
	}

	// Collect all JSON file paths
	var jsonPaths []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-JSON files
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".json") {
			return nil
		}

		jsonPaths = append(jsonPaths, path)
		if err != nil {
			return fmt.Errorf("error walking directory: %v", err)
		}

		return nil
	})
	if err != nil {
		return err
	}
	err = createDeclaration(jsonPaths...)
	if err != nil {
		return err
	}
	fmt.Printf("Finished syncing %s - processed %d declaration files\n", dirPath, len(jsonPaths))
	return nil
}

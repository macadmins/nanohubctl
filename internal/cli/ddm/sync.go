package ddm

import (
	"bufio"
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
	var declJSONPaths []string
	var setPaths []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-JSON files
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(path), ".json") {
			declJSONPaths = append(declJSONPaths, path)
			return nil
		}
		// Match files that start with the word "set" and end with ".txt"
		if strings.HasSuffix(path, ".txt") && strings.HasPrefix(filepath.Base(path), "set") {
			fmt.Printf("Processing %s\n", path)
			setPaths = append(setPaths, path)
		}

		if err != nil {
			return fmt.Errorf("error walking directory: %v", err)
		}

		return nil
	})
	if err != nil {
		return err
	}
	err = createDeclaration(declJSONPaths...)
	if err != nil {
		return err
	}
	err = syncSets(setPaths)
	if err != nil {
		return nil
	}
	fmt.Printf("Synced %d declarations to NanoHUB\n", len(declJSONPaths))
	return nil
}

func syncSets(setPaths []string) error {
	declSets := make(map[string][]string)
	for _, setPath := range setPaths {
		setName := setNameFromPath(setPath)
		declSets[setName] = []string{}
		file, err := os.Open(setPath)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(strings.TrimSpace(line), "#") || line == "" {
				continue
			}
			declSets[setName] = append(declSets[setName], strings.TrimSpace(line))
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file: %v\n", err)
		}
	}
	// Now process the declaratiosn for each set
	for setName, identifiers := range declSets {
		if len(identifiers) == 0 {
			fmt.Printf("No identifiers found for set %s, skipping...\n", setName)
			continue
		}
		err := addSet(setName, identifiers...)
		if err != nil {
			return err
		}
	}
	for setName, items := range declSets {
		fmt.Printf("Synced %d declarations in set '%s'\n", len(items), setName)
	}
	return nil
}

// Derive set name from file name and normalize it
func setNameFromPath(setName string) string {
	setName = filepath.Base(setName)
	setName = strings.TrimSpace(setName)
	setName = strings.TrimPrefix(setName, "set.")
	setName = strings.TrimSuffix(setName, ".txt")
	setName = strings.ToLower(setName)
	return setName
}

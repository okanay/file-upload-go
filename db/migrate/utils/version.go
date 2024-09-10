package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func GetMigrationVersion(dbURL string) (int, error) {
	cmd := exec.Command("migrate", "-database", dbURL, "-path", "db/migration", "version")

	//cmd := exec.Command("migrate", "-database", dbURL, "-path", "sql/migration", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("error running migrate version command: %v, output: %s", err, string(output))
	}

	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		// If there's no output, we assume it's the initial state (version 0)
		return 0, nil
	}

	// Try to parse the last word of the output as the version number
	parts := strings.Fields(outputStr)
	if len(parts) == 0 {
		return 0, fmt.Errorf("unexpected output format from migrate version command: %s", outputStr)
	}

	versionStr := parts[len(parts)-1]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return 0, fmt.Errorf("error parsing version number '%s': %v", versionStr, err)
	}

	return version, nil
}

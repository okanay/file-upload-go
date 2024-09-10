package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/okanay/file-upload-go/db/migrate/utils"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	// Env Configuration
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get the database URL from the environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL is not set in .env file")
	}

	// Get current migration version
	currentVersion, err := utils.GetMigrationVersion(dbURL)
	if err != nil {
		log.Fatalf("Error getting current migration version: %v", err)
	}

	fmt.Printf("Current migration version: %d\n", currentVersion)

	// Ask user for target version
	var targetVersion int
	fmt.Print("Enter the target migration version (0 for full rollback): ")
	_, err = fmt.Scanf("%d", &targetVersion)
	if err != nil {
		log.Fatalf("Invalid input: %v", err)
	}

	if targetVersion >= currentVersion {
		fmt.Println("Target version should be less than current version")
		return
	}

	// Confirm with user
	if targetVersion == 0 {
		fmt.Println("WARNING: You are about to rollback all migrate. This will revert your database to its initial state.")
	}
	fmt.Printf("Are you sure you want to migrate down to version %d? [y/N] ", targetVersion)
	var confirm string
	fmt.Scanf("%s", &confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Println("Migration cancelled")
		return
	}

	// Execute migration
	var cmd *exec.Cmd
	if targetVersion == 0 {
		cmd = exec.Command("migrate", "-database", dbURL, "-path", "db/migration", "down")
	} else {
		cmd = exec.Command("migrate", "-database", dbURL, "-path", "db/migration", "goto", strconv.Itoa(targetVersion))
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error running migrate command: %v", err)
	}

	if targetVersion == 0 {
		fmt.Println("All migrate have been rolled back successfully. The database is now in its initial state.")
	} else {
		fmt.Printf("Migration to version %d completed successfully\n", targetVersion)
	}
}

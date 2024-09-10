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

	// Get current migration version and dirty state
	currentVersion, err := utils.GetMigrationVersion(dbURL)

	if err != nil {
		log.Fatalf("Error getting current migration status: %v", err)
	}

	fmt.Printf("Current migration version: %d\n", currentVersion)

	// Ask user for migration type
	var migrationType string
	fmt.Print("Enter migration type (normal/force): ")
	fmt.Scanf("%s", &migrationType)

	if migrationType != "normal" && migrationType != "force" {
		log.Fatalf("Invalid migration type. Please enter 'normal' or 'force'.")
	}

	// Ask user for target version
	var targetVersion int
	fmt.Print("Enter the target migration version: ")
	_, err = fmt.Scanf("%d", &targetVersion)
	if err != nil {
		log.Fatalf("Invalid input: %v", err)
	}

	if migrationType == "normal" {
		if targetVersion >= currentVersion {
			fmt.Println("For normal migration, target version should be less than current version")
			return
		}
	} else { // force migration
		fmt.Println("\nWARNING: You are about to perform a force migration.")
		fmt.Println("This will NOT change your database schema, but only update the version number.")
		fmt.Println("Use this option only if you are sure that your database schema")
		fmt.Println("matches the desired version and you want to resolve a 'dirty' state.")
	}

	// Confirm with user
	fmt.Printf("\nAre you sure you want to %s migrate to version %d? [y/N] ", migrationType, targetVersion)
	var confirm string
	fmt.Scanf("%s", &confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Println("Migration cancelled")
		return
	}

	// Execute migration
	var cmd *exec.Cmd
	if migrationType == "force" {
		cmd = exec.Command("migrate", "-database", dbURL, "-path", "db/migration", "force", strconv.Itoa(targetVersion))
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

	fmt.Printf("%s migration to version %d completed successfully\n", migrationType, targetVersion)
}

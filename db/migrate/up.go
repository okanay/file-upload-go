package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/okanay/file-upload-go/db/migrate/utils"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
		if strings.Contains(err.Error(), "no migration") {
			fmt.Println("No migrate have been applied yet. The database is in its initial state.")
			currentVersion = 0
		} else {
			log.Fatalf("Error getting current migration status: %v", err)
		}
	}

	fmt.Printf("Current migration version: %d\n", currentVersion)

	// Ask user for number of steps
	var steps int
	fmt.Print("Enter the number of migration steps to apply (or press Enter for all): ")
	input := ""
	fmt.Scanln(&input)
	if input != "" {
		steps, err = strconv.Atoi(input)
		if err != nil {
			log.Fatalf("Invalid input: %v", err)
		}
	}

	var cmd *exec.Cmd
	if steps > 0 {
		cmd = exec.Command("migrate", "-database", dbURL, "-path", "db/migration", "up", strconv.Itoa(steps))
	} else {
		cmd = exec.Command("migrate", "-database", dbURL, "-path", "db/migration", "up")
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error running migrate command: %v", err)
	}

	// Get new migration version
	newVersion, err := utils.GetMigrationVersion(dbURL)
	if err != nil {
		log.Fatalf("Error getting new migration status: %v", err)
	}

	fmt.Printf("Migration successful. New version: %d\n", newVersion)

}

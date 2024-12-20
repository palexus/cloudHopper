package chop

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

var configFile = "/Users/alexanderpreis/Projects/infologistix/cloudHopper/chop.yaml"

func (config *Configuration) ReadAccounts() {
	// Execute the gcloud command
	cmd := exec.Command("gcloud", "config", "configurations", "list")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing gcloud command:", err)
		return
	}

	// Parse the output
	scanner := bufio.NewScanner(&out)
	accounts := make(map[string]bool) // To store unique accounts
	for scanner.Scan() {
		line := scanner.Text()

		// Skip the header row and empty lines
		if strings.Contains(line, "ACCOUNT") || line == "" {
			continue
		}

		// Split columns by spaces (assuming columns are space-separated)
		columns := strings.Fields(line)
		if len(columns) >= 3 {
			account := columns[2] // The third column contains the account
			accounts[account] = true
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading command output:", err)
		return
	}

	// Add accounts to chop
	for account := range accounts {
		config.AddAccount(account)
		fmt.Println("Adding account:", account)
	}

	config.SaveConfigurationToYAML(configFile)
}

func (config *Configuration) ReadProjects(account string) {
	// Execute the gcloud command
	cmd := exec.Command("gcloud", "projects", "list")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing gcloud command:", err)
		return
	}

	// Parse the output
	scanner := bufio.NewScanner(&out)
	projects := make(map[string]bool) // To store unique projects
	for scanner.Scan() {
		line := scanner.Text()

		// Skip the header row and empty lines
		if strings.Contains(line, "PROJECT_ID") || line == "" {
			continue
		}

		// Split columns by spaces (assuming columns are space-separated)
		columns := strings.Fields(line)
		if len(columns) >= 3 {
			project := columns[0] // The second column contains the projectID
			projects[project] = true
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading command output:", err)
		return
	}

	// Add projects to chop
	for project := range projects {
		config.AddProjectToActiveAccount(account, project)
		fmt.Println("Adding project:", project)
	}

	config.SaveConfigurationToYAML(configFile)
}

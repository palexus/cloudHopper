package chop

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Machine represents a machine in a project
type Machine struct {
	Name      string
	LastUsage time.Time
}

// Project represents a project in an account
type Project struct {
	Name     string
	Machines map[string]Machine
}

// Account represents an account with multiple projects
type Account struct {
	Name     string
	Projects map[string]Project
}

// Configuration contains all accounts and the currently active account/project
type Configuration struct {
	Accounts       map[string]Account
	ActiveAccount  string            // Tracks the currently active account
	ActiveProjects map[string]string // Tracks active projects per account
}

// Initializes a new configuration
func NewConfiguration() Configuration {
	return Configuration{
		Accounts: make(map[string]Account),
	}
}

// Adds a new account
func (configs *Configuration) AddAccount(account string) {
	if _, exists := configs.Accounts[account]; exists {
		return // Account already exists
	}
	configs.Accounts[account] = Account{
		Name:     account,
		Projects: make(map[string]Project),
	}
}

// Sets the active account
func (configs *Configuration) SetActiveAccount(account string) error {
	if _, exists := configs.Accounts[account]; !exists {
		return errors.New("account does not exist")
	}
	configs.ActiveAccount = account
	configs.ActiveProjects = make(map[string]string) // Reset the active project
	return nil
}

// AddProjectToActiveAccount adds a project to the active account
func (configs *Configuration) AddProjectToActiveAccount(account string, project string) error {
	// Ensure the account exists
	acc, exists := configs.Accounts[account]
	if !exists {
		return errors.New("account does not exist")
	}

	// Check if the project already exists
	if _, exists := acc.Projects[project]; exists {
		return nil // Project already exists
	}

	// Add the new project
	if acc.Projects == nil {
		acc.Projects = make(map[string]Project)
	}
	acc.Projects[project] = Project{
		Name:     project,
		Machines: make(map[string]Machine),
	}

	// Update the account in the configuration
	configs.Accounts[account] = acc
	return nil
}

// Sets the active project for a specific account
func (configs *Configuration) SetActiveProjectForAccount(account string, project string) error {
	// Ensure the account exists
	acc, exists := configs.Accounts[account]
	if !exists {
		return errors.New("account does not exist")
	}

	// Ensure the project exists in the account
	if _, exists := acc.Projects[project]; !exists {
		return errors.New("project does not exist in the account")
	}

	// Set the active project for the account
	if configs.ActiveProjects == nil {
		configs.ActiveProjects = make(map[string]string)
	}
	configs.ActiveProjects[account] = project
	return nil
}

// Adds a machine to the active project for a specific account
func (configs *Configuration) AddMachineToActiveProject(account string, machine string) error {
	// Ensure the account exists
	acc, exists := configs.Accounts[account]
	if !exists {
		return errors.New("account does not exist")
	}

	// Ensure the account has an active project
	activeProject, activeExists := configs.ActiveProjects[account]
	if !activeExists || activeProject == "" {
		return errors.New("no active project for the specified account")
	}

	// Get the active project
	project, exists := acc.Projects[activeProject]
	if !exists {
		return errors.New("active project not found in the account")
	}

	// Add the machine to the active project
	if _, exists := project.Machines[machine]; exists {
		return nil // Machine already exists
	}

	project.Machines[machine] = Machine{
		Name:      machine,
		LastUsage: time.Now(),
	}
	acc.Projects[activeProject] = project
	configs.Accounts[account] = acc
	return nil
}

// DeleteAccount removes an account and all its projects and machines
func (configs *Configuration) DeleteAccount(account string) error {
	if _, exists := configs.Accounts[account]; !exists {
		return errors.New("account does not exist")
	}

	delete(configs.Accounts, account)

	// If the active account is deleted, unset it
	if configs.ActiveAccount == account {
		configs.ActiveAccount = ""
		delete(configs.ActiveProjects, account)
	}
	return nil
}

// DeleteProject removes a project from an account
func (configs *Configuration) DeleteProject(account string, project string) error {
	// Ensure the account exists
	acc, exists := configs.Accounts[account]
	if !exists {
		return errors.New("account does not exist")
	}

	// Ensure the project exists
	if _, exists := acc.Projects[project]; !exists {
		return errors.New("project does not exist")
	}

	delete(acc.Projects, project)

	// If the active project for the account is deleted, unset it
	if configs.ActiveProjects[account] == project {
		delete(configs.ActiveProjects, account)
	}
	return nil
}

// DeleteMachine removes a machine from a project in an account
func (configs *Configuration) DeleteMachine(account string, project string, machine string) error {
	// Ensure the account exists
	acc, exists := configs.Accounts[account]
	if !exists {
		return errors.New("account does not exist")
	}

	// Ensure the project exists
	proj, exists := acc.Projects[project]
	if !exists {
		return errors.New("project does not exist in the account")
	}

	// Ensure the machine exists
	if _, exists := proj.Machines[machine]; !exists {
		return errors.New("machine does not exist in the project")
	}

	delete(proj.Machines, machine)
	//acc.Projects[project] = proj
	//configs.Accounts[account] = acc
	return nil
}

// UnsetActiveAccount unsets the currently active account
func (configs *Configuration) UnsetActiveAccount() {
	if configs.ActiveAccount == "" {
		fmt.Println("No active account set")
		return
	}
	configs.ActiveAccount = ""
}

// UnsetActiveProjectForAccount unsets the active project for a specific account
func (configs *Configuration) UnsetActiveProjectForAccount(account string) error {
	// If no account is provided, unset the active project for the active account
	if account == "" {
		if configs.ActiveAccount == "" {
			return errors.New("no active account set")

		}
		if _, exists := configs.ActiveProjects[configs.ActiveAccount]; !exists {
			return errors.New("no active project set for the active account")
		}
		// Unset the active project for the account
		delete(configs.ActiveProjects, configs.ActiveAccount)
		fmt.Println("Active project unset for account:", configs.ActiveAccount)
		return nil
	}

	// If account is specified via flag, unset the active project for that account
	if _, exists := configs.Accounts[account]; !exists {
		return errors.New("account does not exist")
	}
	if _, exists := configs.ActiveProjects[account]; !exists {
		return errors.New("no active project set for the specified account")
	}
	// Unset the active project for the specified account
	delete(configs.ActiveProjects, account)
	fmt.Println("Active project unset for account:", account)
	return nil
}

// SaveConfigurationToYAML saves the Configuration to a YAML file
func (configs *Configuration) SaveConfigurationToYAML(filename string) error {
	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Encode the configuration to YAML
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(configs); err != nil {
		return fmt.Errorf("failed to encode configuration to YAML: %w", err)
	}

	return nil
}

// ReadConfigurationFromYAML loads a Configuration from a YAML file
func (configs *Configuration) ReadConfigurationFromYAML(filename string) error {
	// Open the YAML file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode the YAML into the Configuration struct
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(configs); err != nil {
		return fmt.Errorf("failed to decode YAML into configuration: %w", err)
	}

	return nil
}

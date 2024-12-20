/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"palexus/chop/cmd/chop"
	"sort"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var config chop.Configuration

var configFile = "/Users/alexanderpreis/Projects/infologistix/cloudHopper/chop.yaml"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chop",
	Short: "SSH into your favorite machines.",
	Long: `
      __             ____ __                      
 ____/ /__  __ _____/ / // /__  ___  ___  ___ ____
/ __/ / _ \/ // / _  / _  / _ \/ _ \/ _ \/ -_) __/
\__/_/\___/\_,_/\_,_/_//_/\___/ .__/ .__/\__/_/   
                             /_/  /_/             
This little tool helps you to navigate and log into your various machines
	in the different projects of your different accounts.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use one of the subcommands: add, delete, list, save, load")
	},
}

// Set subcommands for 'set'
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the active account or project",
	Run: func(cmd *cobra.Command, args []string) {
		// Provide a default message if no subcommand is provided
		fmt.Println("Please specify 'account' or 'project' to set")
	},
}

// Set the active account
var setAccountCmd = &cobra.Command{
	Use:   "account [account_name]",
	Short: "Set the active account",
	Args:  cobra.ExactArgs(1), // Ensure that exactly one argument (account name) is passed
	Run: func(cmd *cobra.Command, args []string) {
		account := args[0]

		// Set the active account
		err := config.SetActiveAccount(account)
		if err != nil {
			fmt.Println("Error setting account:", err)
		} else {
			fmt.Println("Active account set to:", account)
		}

		// Save configuration after adding
		save_err := config.SaveConfigurationToYAML(configFile)
		if save_err != nil {
			fmt.Println("Error saving configuration:", save_err)
		}
	},
}

// Set the active project for the current account
var setProjectCmd = &cobra.Command{
	Use:   "project [project_name]",
	Short: "Set the active project for the active account",
	Args:  cobra.ExactArgs(1), // Ensure that exactly one argument (project name) is passed
	Run: func(cmd *cobra.Command, args []string) {
		project := args[0]
		account, _ := cmd.Flags().GetString("account")

		// If no account is provided via flag or argument, use the active account
		if account == "" {
			if config.ActiveAccount == "" {
				fmt.Fprintln(os.Stderr, "No active account set or provided. Use --account flag or provide account argument.")
				cmd.Help()
				return
			}
			account = config.ActiveAccount
		}

		// Set the active project for the specified account
		err := config.SetActiveProjectForAccount(account, project)
		if err != nil {
			fmt.Println("Error setting project:", err)
		} else {
			fmt.Println("Active project for account", account, "set to:", project)
		}

		// Save configuration after adding
		save_err := config.SaveConfigurationToYAML(configFile)
		if save_err != nil {
			fmt.Println("Error saving configuration:", save_err)
		}
	},
}

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Unset account or project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify 'account' or 'project' to unset")
	},
}

// Unset active account
var unsetAccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Unset the active account",
	Run: func(cmd *cobra.Command, args []string) {
		config.UnsetActiveAccount()
		fmt.Println("Active account unset")
	},
}

// Unset active project for the current account
var unsetProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Unset the active project for the active account",
	Run: func(cmd *cobra.Command, args []string) {
		account, _ := cmd.Flags().GetString("account")

		set_err := config.UnsetActiveProjectForAccount(account)
		if set_err != nil {
			fmt.Println("Error unsetting project:", set_err)
		}

		// Save configuration after adding
		err := config.SaveConfigurationToYAML(configFile)
		if err != nil {
			fmt.Println("Error saving configuration:", err)
		}
	},
}

// Add subcommands for 'add'
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add accounts, projects, or machines",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify 'account', 'project', or 'machine' to add")
	},
}

// Add multiple accounts
var addAccountCmd = &cobra.Command{
	Use:   "account [account_names...]",
	Short: "Add one or more accounts",
	Args:  cobra.MinimumNArgs(1), // Ensure at least one account name is provided
	Run: func(cmd *cobra.Command, args []string) {
		for _, account := range args {
			// Add each account
			config.AddAccount(account)
			fmt.Println("Account added:", account)
		}
		// Save the configuration after adding accounts
		err := config.SaveConfigurationToYAML(configFile)
		if err != nil {
			fmt.Println("Error saving configuration:", err)
		}
	},
}

// Add multiple projects to the active account (with optional --account flag)
var addProjectCmd = &cobra.Command{
	Use:   "project [project_names...]",
	Short: "Add one or more projects to an account",
	Args:  cobra.MinimumNArgs(1), // Ensure at least one project name is provided
	Run: func(cmd *cobra.Command, args []string) {
		// Get the --account flag value
		account, _ := cmd.Flags().GetString("account")

		// If no account is provided via flag, use the active account
		if account == "" {
			if config.ActiveAccount == "" {
				fmt.Fprintln(os.Stderr, "No account set. Please provide an account using the --account flag or 'chop set account <account>'")
				cmd.Help()
				return
			}
			account = config.ActiveAccount
		}

		for _, project := range args {
			// Add each project to the specified account
			err := config.AddProjectToActiveAccount(account, project)
			if err != nil {
				fmt.Println("Error adding project:", err)
			} else {
				fmt.Println("Project added to", account, ":", project)
			}
		}

		// Save the configuration after adding projects
		err := config.SaveConfigurationToYAML(configFile)
		if err != nil {
			fmt.Println("Error saving configuration:", err)
		}
	},
}

// Add multiple machines to the active project in the active account (with --account and --project flags)
var addMachineCmd = &cobra.Command{
	Use:   "machine [machine_names...]",
	Short: "Add one or more machines to the active project in the active account",
	Args:  cobra.MinimumNArgs(1), // Ensure at least one machine name is provided
	Run: func(cmd *cobra.Command, args []string) {
		// Get the --account and --project flags
		account, _ := cmd.Flags().GetString("account")
		project, _ := cmd.Flags().GetString("project")

		// Ensure the account is set (either via flag or active account)
		if account == "" {
			if config.ActiveAccount == "" {
				fmt.Fprintln(os.Stderr, "No active account set. Please provide an account using --account or 'chop set account <account>'")
				cmd.Help()
				return
			}
			account = config.ActiveAccount
		}

		// Ensure the project is set (either via flag or active project)
		if project == "" {
			activeProject, activeExists := config.ActiveProjects[account]
			if !activeExists || activeProject == "" {
				fmt.Fprintln(os.Stderr, "No active project for the active account. Please provide a project using --project or 'chop set project <project>'")
				cmd.Help()
				return
			}
			project = activeProject
		}

		for _, machine := range args {
			// Add each machine to the specified project of the specified account
			err := config.AddMachineToActiveProject(account, machine)
			if err != nil {
				fmt.Println("Error adding machine:", err)
			} else {
				fmt.Println("Machine added to", account, "->", project, ":", machine)
			}
		}

		// Save the configuration after adding machines
		err := config.SaveConfigurationToYAML(configFile)
		if err != nil {
			fmt.Println("Error saving configuration:", err)
		}
	},
}

// ternary is a helper function that returns trueValue if condition is true, otherwise falseValue.
func ternary(condition bool, trueValue, falseValue string) string {
	if condition {
		return trueValue
	}
	return falseValue
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts, projects, and machines",
	Run: func(cmd *cobra.Command, args []string) {
		// Define colors for active account and project
		activeAccountColor := color.New(color.FgGreen).SprintFunc()
		activeProjectColor := color.New(color.FgCyan).SprintFunc()

		// Create a new simpletable
		table := simpletable.New()
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Text: "ACCOUNT"},
				{Text: "PROJECT"},
				{Text: "MACHINE"},
			},
		}

		// Collect account names and sort them
		accountNames := make([]string, 0, len(config.Accounts))
		for accountName := range config.Accounts {
			accountNames = append(accountNames, accountName)
		}
		sort.Strings(accountNames)

		// Iterate through the sorted accounts and populate the table rows
		for _, accountName := range accountNames {
			account := config.Accounts[accountName]
			accountDisplay := accountName
			if accountName == config.ActiveAccount {
				accountDisplay = activeAccountColor(accountName) + " (active)"
			}

			// Collect project names and sort them
			projectNames := make([]string, 0, len(account.Projects))
			for projectName := range account.Projects {
				projectNames = append(projectNames, projectName)
			}
			sort.Strings(projectNames)

			// Flag to print the account name only once
			printAccount := true

			for _, projectName := range projectNames {
				project := account.Projects[projectName]
				projectDisplay := projectName
				if config.ActiveProjects[accountName] == projectName {
					projectDisplay = activeProjectColor(projectName) + " (active)"
				}

				// Collect machine names and sort them
				machineNames := make([]string, 0, len(project.Machines))
				for machineName := range project.Machines {
					machineNames = append(machineNames, machineName)
				}
				sort.Strings(machineNames)

				// Flag to print the project name only once
				printProject := true

				for _, machineName := range machineNames {
					table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
						{Text: ternary(printAccount, accountDisplay, "")},
						{Text: ternary(printProject, projectDisplay, "")},
						{Text: machineName},
					})

					// After the first machine is printed, suppress further project name printing
					printProject = false
					// After the first machine in the project is printed, suppress further account name printing
					printAccount = false
				}

				// If no machines, still print the project row
				if len(machineNames) == 0 {
					table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
						{Text: ternary(printAccount, accountDisplay, "")},
						{Text: ternary(printProject, projectDisplay, "")},
						{Text: "-"},
					})
					printProject = false
					printAccount = false
				}
			}

			// If no projects, still print the account row
			if len(account.Projects) == 0 {
				table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
					{Text: ternary(printAccount, accountDisplay, "")},
					{Text: "-"},
					{Text: "-"},
				})
			}
		}

		styleList := []*simpletable.Style{
			//simpletable.StyleUnicode,
			//simpletable.StyleCompact,
			//simpletable.StyleCompactClassic,
			//simpletable.StyleCompactLite,
			//simpletable.StyleRounded,
			//simpletable.StyleMarkdown,
			simpletable.StyleDefault,
		}

		for _, style := range styleList { // Use "_" to ignore the index and directly get the value
			// Configure table style
			table.SetStyle(style)

			// Print the table
			fmt.Println(table.String())
		}
	},
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove an account, project, or machine",
	Long:  "Remove an account, a project under a specific account, or a machine from a project.",
}

var rmAccountCmd = &cobra.Command{
	Use:   "account [account_name]",
	Short: "Remove an account",
	Args:  cobra.MinimumNArgs(1), // Ensure at least one account name is provided
	Run: func(cmd *cobra.Command, args []string) {
		for _, account := range args {
			// Remove each account
			err := config.DeleteAccount(account)
			if err != nil {
				fmt.Println("Error removing account:", err)
			} else {
				fmt.Println("Account removed:", account)
			}
		}
		// Save configuration after changes
		save_err := config.SaveConfigurationToYAML(configFile)
		if save_err != nil {
			fmt.Println("Error saving configuration:", save_err)
		}
	},
}

// Remove one or more projects from an account (with optional --account flag)
var rmProjectCmd = &cobra.Command{
	Use:   "project [project_names...]",
	Short: "Remove one or more projects from an account",
	Args:  cobra.MinimumNArgs(1), // Ensure at least one project name is provided
	Run: func(cmd *cobra.Command, args []string) {
		// Get the --account flag value
		account, _ := cmd.Flags().GetString("account")

		// If no account is provided via flag, use the active account
		if account == "" {
			if config.ActiveAccount == "" {
				fmt.Fprintln(os.Stderr, "No account set. Please provide an account using the --account flag or 'chop set account <account>'")
				cmd.Help()
				return
			}
			account = config.ActiveAccount
		}

		for _, project := range args {
			// Remove each project from the specified account
			err := config.DeleteProject(account, project)
			if err != nil {
				fmt.Println("Error removing project:", err)
			} else {
				fmt.Println("Project removed from", account, ":", project)
			}
		}

		// Save the configuration after removing projects
		err := config.SaveConfigurationToYAML(configFile)
		if err != nil {
			fmt.Println("Error saving configuration:", err)
		}
	},
}

// Remove one or more machines from a project in an account (with --account and --project flags)
var rmMachineCmd = &cobra.Command{
	Use:   "machine [machine_names...]",
	Short: "Remove one or more machines from a project",
	Args:  cobra.MinimumNArgs(1), // Ensure at least one machine name is provided
	Run: func(cmd *cobra.Command, args []string) {
		// Get the --account and --project flags
		account, _ := cmd.Flags().GetString("account")
		project, _ := cmd.Flags().GetString("project")

		// Ensure the account is set (either via flag or active account)
		if account == "" {
			if config.ActiveAccount == "" {
				fmt.Fprintln(os.Stderr, "No active account set. Please provide an account using --account or 'chop set account <account>'")
				cmd.Help()
				return
			}
			account = config.ActiveAccount
		}

		// Ensure the project is set (either via flag or active project)
		if project == "" {
			activeProject, activeExists := config.ActiveProjects[account]
			if !activeExists || activeProject == "" {
				fmt.Fprintln(os.Stderr, "No active project for the active account. Please provide a project using --project or 'chop set project <project>'")
				cmd.Help()
				return
			}
			project = activeProject
		}

		for _, machine := range args {
			// Remove each machine from the specified project of the specified account
			err := config.DeleteMachine(account, project, machine)
			if err != nil {
				fmt.Println("Error removing machine:", err)
			} else {
				fmt.Println("Machine removed from", account, "->", project, ":", machine)
			}
		}

		// Save the configuration after removing machines
		err := config.SaveConfigurationToYAML(configFile)
		if err != nil {
			fmt.Println("Error saving configuration:", err)
		}
	},
}

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune the entire configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config.ActiveAccount = ""
		config.ActiveProjects = make(map[string]string)
		config.Accounts = make(map[string]chop.Account)

		reader := bufio.NewReader(os.Stdin)

		// Prompt the user
		fmt.Print("Are you sure you want to do it? [yes/no]: ")

		// Read user input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		// Normalize the input
		input = strings.TrimSpace(strings.ToLower(input))

		// Process the input
		switch input {
		case "yes", "y":
			fmt.Println("You confirmed: Proceeding...")
		case "no", "n":
			fmt.Println("You declined: Aborting...")
			return
		default:
			fmt.Println("Invalid input. Please type 'yes' or 'no'.")
			return
		}

		// Save the configuration after removing machines
		save_err := config.SaveConfigurationToYAML(configFile)
		if save_err != nil {
			fmt.Println("Error saving configuration:", save_err)
		}
		fmt.Println("Done.")
	},
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetches accounts, projects or machines",
	Long:  "Fetches accounts, projects or machines from GCP. (Azure, AWS are not supported, yet)",
}

var fetchAccountCmd = &cobra.Command{
	Use:   "accounts",
	Short: "fetches accounts",
	Long:  "Fetches accounts from your gcloud configurations. (Azure, AWS are not supported, yet)",
	Run: func(cmd *cobra.Command, args []string) {
		config.ReadAccounts()
	},
}

var fetchProjectCmd = &cobra.Command{
	Use:   "projects",
	Short: "fetches accounts",
	Long:  "Fetches accounts from your gcloud configurations. (Azure, AWS are not supported, yet)",
	Run: func(cmd *cobra.Command, args []string) {
		account, _ := cmd.Flags().GetString("account")
		config.ReadProjects(account)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Load the configuration file at startup, if it exists
	err := config.ReadConfigurationFromYAML(configFile)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
	}

	// Ensure the accounts map is initialized if no data was loaded
	if config.Accounts == nil {
		config.Accounts = make(map[string]chop.Account)
		config.SaveConfigurationToYAML(configFile)
		fmt.Println("Created empty configuration file at:", configFile)
	}

	// ************ ADD ***************
	// Add the subcommands to the 'add' parent command
	addCmd.AddCommand(addAccountCmd)
	addCmd.AddCommand(addProjectCmd)
	addCmd.AddCommand(addMachineCmd)
	addProjectCmd.Flags().String("account", "", "The Account where you want to set the Project")
	addMachineCmd.Flags().String("account", "", "The Account where you want to set the Machine")
	addMachineCmd.Flags().String("project", "", "The Project where you want to set the Machine")
	rootCmd.AddCommand(addCmd)

	// ********** LIST ***********
	rootCmd.AddCommand(listCmd)

	// ********** SET **************
	setProjectCmd.Flags().String("account", "", "Account to set the project for (optional)")
	setCmd.AddCommand(setAccountCmd)
	setCmd.AddCommand(setProjectCmd)
	rootCmd.AddCommand(setCmd)

	// ******** REMOVE *************
	rmProjectCmd.Flags().String("account", "", "Specify the account to remove projects from (if not provided, active account will be used)")
	rmMachineCmd.Flags().String("account", "", "Specify the account to remove machines from (if not provided, active account will be used)")
	rmMachineCmd.Flags().String("project", "", "Specify the project to remove machines from (if not provided, active project will be used)")
	rmCmd.AddCommand(rmAccountCmd)
	rmCmd.AddCommand(rmProjectCmd)
	rmCmd.AddCommand(rmMachineCmd)
	rootCmd.AddCommand(rmCmd)

	// ********* UNSET ***********
	unsetCmd.AddCommand(unsetAccountCmd)
	unsetCmd.AddCommand(unsetProjectCmd)
	unsetProjectCmd.Flags().String("account", "", "Account name where to unset the project")
	rootCmd.AddCommand(unsetCmd)

	// ********** PRUNE ************
	rootCmd.AddCommand(pruneCmd)

	// ********** FETCH ************
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.AddCommand(fetchAccountCmd)
	fetchCmd.AddCommand(fetchProjectCmd)
	fetchProjectCmd.Flags().String("account", "", "In which account do you wish to fetch the projects?")
}

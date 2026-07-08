// Package main provides the entry point for the wailsm CLI tool.
// It handles command-line argument parsing and routes requests to the
// appropriate command handlers in the commands package.
package main

import (
	"fmt"
	"os"

	"github.com/sweet-juice/sweetjuice/cmd/commands"
)

func main() {
	if len(os.Args) < 2 {
		commands.ShowUsage()
	}

	// Route the command line arguments to the corresponding handler functions
	switch os.Args[1] {
	case "--new":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: Please specify a project directory name.")
			os.Exit(1)
		}
		commands.CreateNewProject(os.Args[2])
	case "--refresh":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: Please specify a target platform: 'android' or 'ios'")
			os.Exit(1)
		}
		commands.ExecuteRefresh(os.Args[2])
	case "--build":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Error: Missing arguments. Usage: wailsm --build <platform> <debug|release>")
			os.Exit(1)
		}
		commands.ExecuteBuild(os.Args[2], os.Args[3])
	case "--run":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: Please specify a target platform: 'android' or 'ios'")
			os.Exit(1)
		}
		commands.ExecuteRun(os.Args[2])
	case "--add":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: Please provide a valid plugin repository path.")
			os.Exit(1)
		}
		commands.ManagePlugin("add", os.Args[2])
	case "--remove":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: Please provide a valid plugin repository path.")
			os.Exit(1)
		}
		commands.ManagePlugin("remove", os.Args[2])
	case "-h", "--help":
		commands.ShowUsage()
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown option '%s'\n", os.Args[1])
		commands.ShowUsage()
	}
}

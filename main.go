package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jasonleenaylor/git-branch-cleaner/branchcleaner"
)

var standardBranches = []string{"master", "main", "develop", "release/*"}

func main() {
	// Check if help is the first argument or no arguments are provided
	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage()
		return
	}

	// Default filter list (empty)
	var filterBranches []string
	var dryRun bool
	args := os.Args[1:]

	// Check for --all and enforce it to be used only with --dry-run
	if args[0] == "--all" {
		if len(args) > 2 || (len(args) == 2 && args[1] != "--dry-run") {
			fmt.Println("--all can only be used alone or with --dry-run.")
			printUsage()
			return
		}
		// Add the current branch to the filter list
		currentBranch, err := getCurrentBranch()
		if err != nil {
			fmt.Println("Error retrieving current branch:", err)
			return
		}
		filterBranches = append(filterBranches, currentBranch)
		if len(args) == 2 && args[1] == "--dry-run" {
			dryRun = true
		}
	} else {
		// Parse remaining arguments
		for i, arg := range args {
			switch {
			case arg == "--standard":
				if i > 0 {
					fmt.Println("")
				}

				// Add standard branches to the filter list
				filterBranches = append(filterBranches, standardBranches...)
			case strings.HasPrefix(arg, "--exclude:"):
				// Exclude branches passed with --exclude:<branch-pattern>
				filterBranches = append(filterBranches, strings.TrimPrefix(arg, "--exclude:"))
			case arg == "--dry-run":
				dryRun = true
			default:
				// If an unknown argument is passed, print usage and exit
				fmt.Printf("Unknown argument: %s\n", arg)
				printUsage()
				return
			}
		}
	}

	// List all local branches
	branches, err := getLocalBranches()
	if err != nil {
		fmt.Println("Error retrieving local branches:", err)
		return
	}

	// Filter branches based on the provided args and filter list
	filteredBranches := branchcleaner.FilterBranches(branches, filterBranches)

	// Print initial dry-run message
	if dryRun {
		fmt.Println("Dry Run: The following branches would be deleted:")
	}

	// Delete the filtered branches or simulate it in dry-run mode
	err = deleteBranches(filteredBranches, dryRun)
	if err != nil {
		fmt.Println("Error deleting branches:", err)
	} else if !dryRun {
		fmt.Println("Branches deleted:", strings.Join(filteredBranches, ", "))
	}
}

// Get all local branches
func getLocalBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "--list")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Split the output into individual branches
	branches := strings.Split(string(output), "\n")
	var result []string
	for _, branch := range branches {
		branch = strings.TrimSpace(branch)
		if branch != "" {
			result = append(result, branch)
		}
	}
	return result, nil
}

// Get the current branch
func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// Delete the branches using 'git branch -D' or simulate in dry-run mode
func deleteBranches(branches []string, dryRun bool) error {
	for _, branch := range branches {
		if dryRun {
			// Print branch to be deleted in dry-run mode
			fmt.Println(branch)
		} else {
			// Execute the 'git branch -D' command to delete the branch
			cmd := exec.Command("git", "branch", "-D", branch)
			err := cmd.Run()
			if err != nil {
				return fmt.Errorf("failed to delete branch %s: %w", branch, err)
			}
		}
	}
	return nil
}

// Print usage statement with standard branches
func printUsage() {
	fmt.Println("Git Branch Cleaner Usage:")
	fmt.Println("Usage: git-branch-cleaner [--standard] [--exclude:<branch-pattern>] [--all] [--dry-run] [--help]")
	fmt.Println("\nOptions:")
	fmt.Println("  --standard                 Include the standard branches (master, main, develop, release/*) in the filter list.")
	fmt.Println("  --exclude:<branch-pattern> Exclude branches matching the pattern from deletion.")
	fmt.Println("  --all                      Delete all branches except the current branch (must be used alone or with --dry-run).")
	fmt.Println("  --dry-run                  Preview the branches that would be deleted without actually deleting them.")
	fmt.Println("  --help                     Show this help message.")
	fmt.Println("\nStandard branches that will never be deleted (unless excluded):")
	fmt.Printf("  %s\n", strings.Join(standardBranches, ", "))
	fmt.Println("\nExample Usage:")
	fmt.Println("  git-branch-cleaner --standard --exclude:feature/* --exclude:bugfix/*")
	fmt.Println("  git-branch-cleaner --all")
	fmt.Println("  git-branch-cleaner --dry-run --standard")
	fmt.Println("  git-branch-cleaner --help")
}

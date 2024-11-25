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

	// Parse remaining arguments
	for i, arg := range args {
		switch {
		case arg == "--all":
			if len(args) > 2 || (len(args) == 2 && args[1] != "--dry-run") {
				fmt.Println("--all can only be used alone or with --dry-run.")
				printUsage()
				return
			}
			// Add the current branch to the filter list
			_, currentBranch, err := getLocalBranches()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			filterBranches = []string{currentBranch}
			fmt.Printf("Current branch: %s\n", currentBranch)
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

	// List all local branches
	branches, currentBranch, err := getLocalBranches()
	if err != nil {
		fmt.Println("Error retrieving local branches:", err)
		return
	}

	// Filter branches based on the provided args and filter list
	filteredBranches := branchcleaner.FilterBranches(branches, filterBranches)

	if branchesContains(filteredBranches, currentBranch) {
		fmt.Print("Deleting the current branch is not supported. Change to a different branch.")
		return
	}

	// Print initial dry-run message
	if dryRun {
		fmt.Println("Dry Run: The following branches would be deleted:")
	}

	// Delete the filtered branches or simulate it in dry-run mode
	err = deleteBranches(filteredBranches, dryRun)
	if err != nil {
		fmt.Println("Error deleting branches:", err)
	}
}

func branchesContains(branches []string, test string) bool {
	for _, b := range branches {
		if b == test {
			return true
		}
	}
	return false
}

func getLocalBranches() ([]string, string, error) {
	cmd := exec.Command("git", "branch")
	output, err := cmd.Output()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get branches: %v", err)
	}

	var branches []string
	var currentBranch string

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		if branch == "" {
			continue
		}
		if strings.HasPrefix(branch, "*") {
			currentBranch = strings.TrimPrefix(branch, "* ")
		}
		branches = append(branches, strings.TrimPrefix(branch, "* "))
	}

	return branches, currentBranch, nil
}

// Delete the branches using 'git branch -D' or simulate in dry-run mode
func deleteBranches(branches []string, dryRun bool) error {
	for _, branch := range branches {
		// Print branch to be deleted
		fmt.Println(branch)
		if !dryRun {
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

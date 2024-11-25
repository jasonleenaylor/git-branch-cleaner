// branchcleaner.go
package branchcleaner

import (
	"strings"
	"fmt"
)

// FilterBranches applies the exclude filters to the list of branches.
func FilterBranches(branches []string, exclude []string) []string {
	var filteredBranches []string

	for _, branch := range branches {
		if shouldExclude(exclude, branch) {
			continue
		}
		filteredBranches = append(filteredBranches, branch)
	}

	return filteredBranches
}

// shouldExclude checks if a branch should be excluded based on the exclude patterns.
func shouldExclude(exclude []string, branch string) bool {
	for _, pattern := range exclude {
		if matchPattern(pattern, branch) {
			return true
		}
	}
	return false
}

// matchPattern matches a branch name against a wildcard pattern (e.g., release/*).
func matchPattern(pattern, branch string) bool {
	if strings.HasSuffix(pattern, "/*") {
		return strings.HasPrefix(branch, strings.TrimSuffix(pattern, "/*"))
	}
	return branch == pattern
}

// CleanBranches deletes the branches based on the logic.
func CleanBranches(branches []string, exclude []string) {
	filtered := FilterBranches(branches, exclude)
	for _, branch := range filtered {
		// Delete branch (this is just a placeholder action)
		fmt.Printf("Deleting branch: %s\n", branch)
	}
}


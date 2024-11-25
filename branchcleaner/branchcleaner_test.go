package branchcleaner

import (
	"testing"
)

func TestFilterBranches(t *testing.T) {
	tests := []struct {
		branches  []string
		exclude   []string
		expected  []string
		expectErr bool // Flag to indicate if we expect the test to fail
	}{
		{
			// Test branch wildcard exclusion
			branches: []string{"master", "develop", "release/1.0"},
			exclude:  []string{"release/*"},
			expected: []string{"master", "develop"},
		},
		{
			// Test specific prefix style branch name exclusions
			branches: []string{"master", "develop", "feature/new-feature", "release/1.0"},
			exclude:  []string{"feature/new-feature"},
			expected: []string{"master", "develop", "release/1.0"},
		},
		{
			// Test normal branch name exclusion
			branches:  []string{"master", "main", "develop"},
			exclude:   []string{"master"},
			expected:  []string{"main", "develop"}, // intentionally wrong expected
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			// Call the function to filter the branches
			actual := FilterBranches(test.branches, test.exclude)

			// Compare the actual result with the expected
			if len(actual) != len(test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, actual)
				if test.expectErr {
					t.Log("Test expected to fail, handling failure.")
				}
			} else {
				for i, branch := range actual {
					if branch != test.expected[i] {
						t.Errorf("Expected %v at index %d, got %v", test.expected[i], i, branch)
					}
				}
			}
		})
	}
}

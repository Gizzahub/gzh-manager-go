package github

import (
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	org := "your-organization" // Replace with the actual organization name

	repos, err := List(org)
	if err != nil {
		t.Fatalf("Failed to list repositories: %v", err)
	}

	fmt.Printf("Repositories in organization %s:\n", org)
	for _, repo := range repos {
		fmt.Println(repo)
	}
}

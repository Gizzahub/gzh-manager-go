package github_org

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type RepoInfo struct {
	DefaultBranch string `json:"default_branch"`
}

func GetDefaultBranch(org string, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", org, repo)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get repository info: %s", resp.Status)
	}

	var repoInfo RepoInfo
	if err := json.NewDecoder(resp.Body).Decode(&repoInfo); err != nil {
		return "", err
	}

	return repoInfo.DefaultBranch, nil
}

func List(org string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos", org)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get repositories: %s", resp.Status)
	}

	var repos []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var repoNames []string
	for _, repo := range repos {
		repoNames = append(repoNames, repo.Name)
	}

	return repoNames, nil
}

func Clone(org string, repo string, branch string) error {
	if branch == "" {
		defaultBranch, err := GetDefaultBranch(org, repo)
		if err != nil {
			return fmt.Errorf("failed to get default branch: %w", err)
		}
		branch = defaultBranch
	}

	cloneURL := fmt.Sprintf("https://github.com/%s/%s.git", org, repo)
	cmd := exec.Command("git", "clone", "-b", branch, cloneURL)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}

// RefreshAll synchronizes the repositories in the targetPath with the repositories in the given organization.
func RefreshAll(targetPath string, org string) error {
	// Get all directories inside targetPath
	targetRepos, err := getDirectories(targetPath)
	if err != nil {
		return fmt.Errorf("failed to get directories in target path: %w", err)
	}

	// Get all repositories from the organization
	orgRepos, err := List(org)
	if err != nil {
		return fmt.Errorf("failed to list repositories from organization: %w", err)
	}

	// Determine repos to delete (targetRepos - orgRepos)
	reposToDelete := difference(targetRepos, orgRepos)

	// Delete repos that are not in the organization
	for _, repo := range reposToDelete {
		repoPath := filepath.Join(targetPath, repo)
		if err := os.RemoveAll(repoPath); err != nil {
			return fmt.Errorf("failed to delete repository %s: %w", repo, err)
		}
	}

	// Clone or reset hard HEAD for each repository in the organization
	for _, repo := range orgRepos {
		repoPath := filepath.Join(targetPath, repo)
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			// Clone the repository if it does not exist
			if err := Clone(org, repo, ""); err != nil {
				return fmt.Errorf("failed to clone repository %s: %w", repo, err)
			}
		} else {
			// Reset hard HEAD if the repository already exists
			cmd := exec.Command("git", "-C", repoPath, "reset", "--hard", "HEAD")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to reset repository %s: %w", repo, err)
			}
		}
	}

	return nil
}

// getDirectories returns a list of directory names in the given path.
func getDirectories(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
}

// difference returns the elements in 'a' that are not in 'b'.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

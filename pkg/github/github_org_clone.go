package github

import (
	"bytes"
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

func Clone(targetPath string, org string, repo string, branch string) error {
	if branch == "" {
		defaultBranch, err := GetDefaultBranch(org, repo)
		if err != nil {
			return fmt.Errorf("failed to get default branch: %w", err)
		}
		branch = defaultBranch
	}

	cloneURL := fmt.Sprintf("https://github.com/%s/%s.git", org, repo)
	var out bytes.Buffer
	var stderr bytes.Buffer
	//cmd := exec.Command("git", "clone", "-b", branch, cloneURL, targetPath)
	cmd := exec.Command("git", "clone", cloneURL, targetPath)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(stderr.String())
		fmt.Println(out.String())
		return fmt.Errorf("Clone Failed  (url: %s, branch: %s, targetPath: %s, err: %w)\n", cloneURL, branch, targetPath, err)
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
			return fmt.Errorf("failed to delete repository %s: %w", repoPath, err)
		}
	}

	// Clone or reset hard HEAD for each repository in the organization
	for _, repo := range orgRepos {
		repoPath := filepath.Join(targetPath, repo)
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			// Clone the repository if it does not exist
			if err := Clone(repoPath, org, repo, ""); err != nil {
				//fmt.Printf("failed to clone repository %s: %w\n", repoPath, err)
				return fmt.Errorf("failed to clone repository %s: %w", repoPath, err)
			}
		} else {
			// Reset hard HEAD if the repository already exists
			cmd := exec.Command("git", "-C", repoPath, "reset", "--hard", "HEAD")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to reset repository %s: %w", repoPath, err)
			}
			fmt.Printf("Repo Clone or Reset Success: %s\n", repoPath)
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

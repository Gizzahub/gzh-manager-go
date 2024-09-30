package gitea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type RepoInfo struct {
	DefaultBranch string `json:"default_branch"`
}

func GetDefaultBranch(org string, repo string) (string, error) {
	url := fmt.Sprintf("https://gitea.com/api/v1/repos/%s/%s", org, repo)
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
	url := fmt.Sprintf("https://gitea.com/api/v1/orgs/%s/repos", org)
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

	cloneURL := fmt.Sprintf("https://gitea.com/%s/%s.git", org, repo)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("git", "clone", "-b", branch, cloneURL, targetPath)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(stderr.String())
		fmt.Println(out.String())
		return fmt.Errorf("Clone Failed (url: %s, branch: %s, targetPath: %s, err: %w)\n", cloneURL, branch, targetPath, err)
	}

	return nil
}

func RefreshAll(targetPath string, org string) error {
	repos, err := List(org)
	if err != nil {
		return fmt.Errorf("failed to list repositories: %w", err)
	}

	for _, repo := range repos {
		err := Clone(targetPath, org, repo, "")
		if err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	}

	return nil
}

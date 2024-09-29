package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type GitLabRepoInfo struct {
	DefaultBranch string `json:"default_branch"`
}

func GetDefaultBranch(group string, repo string) (string, error) {
	url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s%%2F%s", group, repo)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get repository info: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var gitLabRepo GitLabRepoInfo
	err = json.Unmarshal(body, &gitLabRepo)
	if err != nil {
		return "", err
	}

	return gitLabRepo.DefaultBranch, nil
}

func listGroupRepos(group string, allRepos *[]string) error {
	url := fmt.Sprintf("https://gitlab.com/api/v4/groups/%s/projects", group)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get repositories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get repositories: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var repos []struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(body, &repos)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		*allRepos = append(*allRepos, repo.Name)
	}

	// Get subgroups
	subgroupsURL := fmt.Sprintf("https://gitlab.com/api/v4/groups/%s/subgroups", group)
	subgroupsResp, err := http.Get(subgroupsURL)
	if err != nil {
		return fmt.Errorf("failed to get subgroups: %w", err)
	}
	defer subgroupsResp.Body.Close()

	if subgroupsResp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get subgroups: %s", subgroupsResp.Status)
	}

	subgroupsBody, err := io.ReadAll(subgroupsResp.Body)
	if err != nil {
		return err
	}

	var subgroups []struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(subgroupsBody, &subgroups)
	if err != nil {
		return err
	}

	for _, subgroup := range subgroups {
		err := listGroupRepos(subgroup.ID, allRepos)
		if err != nil {
			return err
		}
	}

	return nil
}

func List(group string) ([]string, error) {
	var allRepos []string
	err := listGroupRepos(group, &allRepos)
	if err != nil {
		return nil, err
	}
	return allRepos, nil
}

func Clone(targetPath string, group string, repo string, branch string) error {
	if branch == "" {
		defaultBranch, err := GetDefaultBranch(group, repo)
		if err != nil {
			return fmt.Errorf("failed to get default branch: %w", err)
		}
		branch = defaultBranch
	}

	cloneURL := fmt.Sprintf("https://gitlab.com/%s/%s.git", group, repo)
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

func RefreshAll(targetPath string, group string) error {
	// Get all directories inside targetPath
	targetRepos, err := getDirectories(targetPath)
	if err != nil {
		return fmt.Errorf("failed to get directories in target path: %w", err)
	}

	// Get all repositories from the group
	groupRepos, err := List(group)
	if err != nil {
		return fmt.Errorf("failed to list repositories from group: %w", err)
	}

	// Determine repos to delete (targetRepos - groupRepos)
	reposToDelete := difference(targetRepos, groupRepos)

	// Delete repos that are not in the group
	for _, repo := range reposToDelete {
		repoPath := filepath.Join(targetPath, repo)
		if err := os.RemoveAll(repoPath); err != nil {
			return fmt.Errorf("failed to delete repository %s: %w", repoPath, err)
		}
	}

	// Clone or reset hard HEAD for each repository in the group
	for _, repo := range groupRepos {
		repoPath := filepath.Join(targetPath, repo)
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			// Clone the repository if it does not exist
			if err := Clone(repoPath, group, repo, ""); err != nil {
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

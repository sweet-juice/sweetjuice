package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// EnsureGHCLILoggedIn checks if the 'gh' CLI is installed.
func EnsureGHCLILoggedIn() {
	if !CommandExists("gh") {
		Fatal("GitHub CLI (gh) missing", fmt.Errorf("please install the GitHub CLI: https://cli.github.com"))
	}
}

// ForkAndCloneRepo forks a repository and clones it locally if it doesn't exist.
func ForkAndCloneRepo(repoOwner, repoName, targetLocalPath string) {
	if DirExists(targetLocalPath) {
		return
	}

	fmt.Printf("Forking and cloning %s/%s to %s...\n", repoOwner, repoName, targetLocalPath)

	// Fork the repo (fails gracefully if already forked)
	_ = exec.Command("gh", "repo", "fork", repoOwner+"/"+repoName, "--clone=false").Run()

	// Get current user
	out, err := exec.Command("gh", "api", "user", "-q", ".login").Output()
	if err != nil {
		Fatal("Failed to get GitHub user info", err)
	}
	user := strings.TrimSpace(string(out))

	// Clone the fork
	RunCmd("gh", "repo", "clone", user+"/"+repoName, targetLocalPath)
}

// WaitForActionFinish waits for the latest GitHub Action run to complete and returns its ID.
func WaitForActionFinish(repoPath string) string {
	fmt.Println("Waiting for GitHub Action to start and complete...")

	origWd, _ := os.Getwd()
	_ = os.Chdir(repoPath)
	defer func() { _ = os.Chdir(origWd) }()

	for i := 0; i < 60; i++ { // Wait up to 10 minutes
		time.Sleep(10 * time.Second)

		out, err := exec.Command("gh", "run", "list", "--limit", "1", "--json", "status,conclusion,databaseId").Output()
		if err != nil {
			continue
		}

		output := string(out)
		if strings.Contains(output, "\"status\":\"completed\"") {
			if strings.Contains(output, "\"conclusion\":\"success\"") {
				// Get the run ID
				parts := strings.Split(output, "\"databaseId\":")
				if len(parts) > 1 {
					idPart := strings.Split(parts[1], "}")[0]
					idPart = strings.Split(idPart, ",")[0]
					idPart = strings.Split(idPart, "]")[0]
					return strings.Trim(strings.TrimSpace(idPart), "\"")
				}
			} else if strings.Contains(output, "\"conclusion\":\"failure\"") {
				Fatal("GitHub Action failed", fmt.Errorf("check 'gh run view' for details"))
			}
		}
		fmt.Print(".")
	}

	Fatal("Timeout waiting for GitHub Action", fmt.Errorf("the build is taking too long"))
	return ""
}

// DownloadArtifact downloads an artifact from a GitHub run.
func DownloadArtifact(repoPath, runID, artifactName, destPath string) {
	fmt.Printf("\nDownloading artifact %s...\n", artifactName)

	origWd, _ := os.Getwd()
	_ = os.Chdir(repoPath)
	defer func() { _ = os.Chdir(origWd) }()

	RunCmd("gh", "run", "download", runID, "--name", artifactName, "--dir", destPath)
}

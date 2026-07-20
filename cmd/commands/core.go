// Package commands implements the core orchestration logic for the wailsm CLI.
// It manages high-level workflows such as project bootstrapping, plugin management,
// and platform-specific pipeline execution.
package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sweet-juice/sweetjuice/cmd/android"
	"github.com/sweet-juice/sweetjuice/cmd/ios"
	"github.com/sweet-juice/sweetjuice/cmd/utils"
)

const (
	// Version is the current version of the wailsm CLI.
	Version = "v1.4.0"
)

// ShowUsage prints the help information for the Sweet Juice CLI.
func ShowUsage() {
	fmt.Println("Sweet Juice Toolchain CLI (juice)")
	fmt.Println("Usage:")
	fmt.Println("  juice --new <project_name>        Create a fresh project from the template")
	fmt.Println("  juice --refresh <platform>        Run platform sync: 'android' or 'ios'")
	fmt.Println("  juice --build <platform> <mode>   Compile binaries: 'debug' (APK/IPA), 'release' (APK/IPA), or 'bundle' (AAB)")
	fmt.Println("  juice --run <platform>            Compile, install, and execute application via ADB or xtool")
	fmt.Println("  juice --run-cross <platform>      Cloud build (GitHub Actions), install, and execute")
	fmt.Println("  juice --setup cross               Setup GitHub Action based cross-compilation for iOS")
	fmt.Println("  juice --add <plugin-url>          Install a native Go/Mobile plugin")
	fmt.Println("  juice --remove <plugin-url>       Uninstall a native Go/Mobile plugin")
	os.Exit(1)
}

// CreateNewProject bootstraps a new Sweet Juice project in the specified directory.
// It locates the template from the Go module cache and scaffolds the native projects.
func CreateNewProject(targetDir string) {
	fmt.Fprintf(os.Stdout, "=== Creating Project: %s [%s] ===\n", targetDir, Version)
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Directory '%s' already exists. Aborting.\n", targetDir)
		os.Exit(1)
	}

	if _, err := exec.LookPath("go"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Required system tool dependency 'go' is missing.\n")
		os.Exit(1)
	}

	fmt.Println("Locating Sweet Juice core...")
	// 1. Try to find the module in the current context (handles development in-repo)
	out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/sweet-juice/sweetjuice").Output()
	coreDir := strings.TrimSpace(string(out))

	// 2. If not found or empty, try to get it from module cache/remote
	if err != nil || coreDir == "" {
		fmt.Println("Core not found in local context, checking module cache...")
		// Use @latest to ensure it's looked up in the global index if not in current go.mod
		out, err = exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/sweet-juice/sweetjuice@latest").Output()
		coreDir = strings.TrimSpace(string(out))
	}

	// 3. Fallback: if we are still empty, try to find the binary's parent if it looks like the repo
	if coreDir == "" {
		exePath, _ := os.Executable()
		// If running from dist/linux/juice, go up 3 levels
		repoCandidate := filepath.Join(filepath.Dir(exePath), "..", "..")
		if utils.FileExists(filepath.Join(repoCandidate, "go.mod")) {
			coreDir, _ = filepath.Abs(repoCandidate)
		}
	}

	if coreDir == "" {
		utils.Fatal("Failed to locate Sweet Juice core directory", fmt.Errorf("please ensure you have github.com/sweet-juice/sweetjuice installed or are running from the source repo"))
	}

	templatePath := filepath.Join(coreDir, "AppTemplate")

	if !utils.DirExists(templatePath) {
		utils.Fatal("Could not find AppTemplate in core module directory", fmt.Errorf("path missing: %s", templatePath))
	}

	fmt.Println("Scaffolding project from local template...")
	if err := utils.CopyDirectory(templatePath, targetDir); err != nil {
		utils.Fatal("Failed to copy template to target directory", err)
	}

	// Remove node_modules from target if it was copied (to prevent path issues)
	_ = os.RemoveAll(filepath.Join(targetDir, "frontend", "node_modules"))

	// Update go.mod in the target directory to point to the local coreDir for development
	// or remove the replace if we want to use the module cache version.
	// Presuming the user wants to work with the version they just used to scaffold.
	targetGoMod := filepath.Join(targetDir, "go.mod")
	if utils.FileExists(targetGoMod) && coreDir != "" {
		fmt.Println("Configuring project dependencies...")
		// Rename the module to the target directory name
		utils.RunCmd("go", "mod", "edit", "-module="+targetDir, targetGoMod)
		// Remove existing replace if any
		utils.RunCmd("go", "mod", "edit", "-dropreplace=github.com/sweet-juice/sweetjuice", targetGoMod)
		// Add new replace pointing to the absolute coreDir
		utils.RunCmd("go", "mod", "edit", "-replace=github.com/sweet-juice/sweetjuice="+coreDir, targetGoMod)
	}

	android.SetupAndroidLocalProperties(targetDir)

	origWd, _ := os.Getwd()
	_ = os.Chdir(targetDir)
	defer func() { _ = os.Chdir(origWd) }()

	// Trigger frontend build as specified in config.ini
	utils.BuildFrontend()

	// Initialize Go Mobile build tools only if missing
	utils.EnsureGoMobileTools()

	if utils.FileExists("go.mod") {
		fmt.Println("Valid Go module context detected. Binding tool tracking dependencies...")
		utils.RunCmd("go", "mod", "tidy")
		// gobind is required as a tool for the project
		if !utils.CommandExists("gobind") {
			utils.RunCmd("go", "get", "-tool", "golang.org/x/mobile/cmd/gobind")
		}
	}

	fmt.Fprintf(os.Stdout, "=== Setup complete! Your project is ready in ./%s ===\n", targetDir)
}

// ExecuteRefresh triggers a platform-specific synchronization pass.
// For Android, this typically runs gomobile bind and syncs native plugins.
// For iOS, it runs gomobile bind natively.
func ExecuteRefresh(platform string) {
	if platform != "android" && platform != "ios" {
		fmt.Fprintln(os.Stderr, "Error: Please specify a valid target platform: 'android' or 'ios'")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "=== Executing Platform Refresh: %s ===\n", platform)
	if platform == "android" {
		android.RefreshPipeline()
	} else {
		ios.RefreshPipeline()
	}
}

// ExecuteBuild compiles the project for the specified platform and mode using native tools.
// It automatically triggers a refresh pass before compilation.
func ExecuteBuild(platform, mode string) {
	mode = strings.ToLower(mode)
	if mode != "debug" && mode != "release" && mode != "bundle" {
		fmt.Fprintln(os.Stderr, "Error: Invalid build mode specified. Use 'debug' (APK/IPA), 'release' (APK/IPA), or 'bundle' (AAB).")
		os.Exit(1)
	}

	ExecuteRefresh(platform)

	if platform == "android" {
		android.BuildPipeline(mode)
	} else {
		ios.BuildPipeline(mode)
	}
}

// ExecuteRun compiles, installs, and launches the application on a connected device.
func ExecuteRun(platform string) {
	if platform == "android" {
		ExecuteBuild("android", "debug")
		android.RunPipeline()
	} else if platform == "ios" {
		ExecuteBuild("ios", "debug")
		ios.RunPipeline()
	} else {
		fmt.Fprintln(os.Stderr, "Error: Invalid platform. Use 'android' or 'ios'.")
		os.Exit(1)
	}
}

// ExecuteRunCross performs a cloud-based build using GitHub Actions for platforms
// that cannot be built locally (e.g., iOS on Linux).
func ExecuteRunCross(platform string) {
	if platform != "ios" {
		fmt.Println("Cross-build is currently only optimized for 'ios'. Running standard local build for others...")
		ExecuteRun(platform)
		return
	}

	config := utils.LoadConfig()
	githubUser := config.GetOrDefault("cross", "github_user", "")
	crossRepoPath := config.GetOrDefault("cross", "cross_repo_path", "")

	if githubUser == "" || crossRepoPath == "" {
		fmt.Println("Cross-build not configured. Please run 'juice --setup cross' first.")
		os.Exit(1)
	}

	fmt.Printf("=== Initiating Cloud Cross-Build for %s ===\n", platform)
	utils.EnsureGHCLILoggedIn()

	// Ensure the repo path exists
	if !utils.DirExists(crossRepoPath) {
		fmt.Printf("Cross-build repository not found at %s. Attempting to restore...\n", crossRepoPath)
		utils.RunCmd("gh", "repo", "clone", githubUser+"/juice-cross", crossRepoPath)
	}

	fmt.Println("Syncing codebase to cloud builder...")
	// Clear old sync files (but keep .git and .github)
	files, _ := os.ReadDir(crossRepoPath)
	for _, f := range files {
		if f.Name() == ".git" || f.Name() == ".github" {
			continue
		}
		_ = os.RemoveAll(filepath.Join(crossRepoPath, f.Name()))
	}

	// We only sync the parts needed for binding: Go code and Frontend dist
	if err := utils.CopyDirectory(".", crossRepoPath); err != nil {
		utils.Fatal("Failed to sync code to cross-repo", err)
	}

	// Ensure the workflow is up-to-date from the internal template if possible
	// For now we assume the fork has a valid bind.yml as we just set it up

	// Cleanup destination from local-specific files
	_ = os.RemoveAll(filepath.Join(crossRepoPath, "native"))
	_ = os.RemoveAll(filepath.Join(crossRepoPath, "build"))
	_ = os.RemoveAll(filepath.Join(crossRepoPath, "temps"))

	// Build frontend locally first to ensure dist is populated
	utils.BuildFrontend()

	// Push changes to trigger GitHub Action
	origWd, _ := os.Getwd()
	_ = os.Chdir(crossRepoPath)
	utils.RunCmd("git", "add", ".")
	// It's okay if commit fails because of no changes
	_ = exec.Command("git", "commit", "-m", "chore: automated cross-build sync").Run()
	utils.RunCmd("git", "push", "origin", "main")
	_ = os.Chdir(origWd)

	// Wait for action
	utils.WaitForActionFinish(crossRepoPath)

	fmt.Println("\nDownloading built bindings from GitHub Release...")
	iosNativePath := filepath.Join("native", "ios")
	_ = os.MkdirAll(iosNativePath, 0755)

	zipPath := filepath.Join(iosNativePath, "Sweetjuice.xcframework.zip")
	releaseURL := fmt.Sprintf("https://github.com/%s/juice-cross/releases/latest/download/Sweetjuice.xcframework.zip", githubUser)

	if err := utils.DownloadFile(releaseURL, zipPath); err != nil {
		utils.Fatal("Failed to download built framework from release", err)
	}

	fmt.Println("Extracting framework...")
	if err := utils.UnzipTarget(zipPath, iosNativePath); err != nil {
		utils.Fatal("Failed to extract framework", err)
	}
	_ = os.Remove(zipPath)

	fmt.Println("Cloud build integration successful. Launching on local hardware...")
	ios.RunPipeline()
}

// ExecuteSetup handles toolchain or project-level configuration tasks.
func ExecuteSetup(target string) {
	if target != "cross" {
		fmt.Fprintf(os.Stderr, "Error: Unknown setup target '%s'. Available: 'cross'\n", target)
		os.Exit(1)
	}

	fmt.Println("=== Setting up Cross-Compilation Environment ===")
	utils.EnsureGHCLILoggedIn()

	// Get current GitHub user
	out, err := exec.Command("gh", "api", "user", "-q", ".login").Output()
	if err != nil {
		utils.Fatal("Failed to get GitHub user info", err)
	}
	githubUser := strings.TrimSpace(string(out))
	fmt.Printf("Detected GitHub user: %s\n", githubUser)

	home, _ := os.UserHomeDir()
	crossRepoPath := filepath.Join(home, ".sweetjuice", "juice-cross")

	fmt.Println("\nStep 1: Forking sweet-juice/juice-cross...")
	fmt.Println("Please ensure you have manually forked https://github.com/sweet-juice/juice-cross to your account.")

	fmt.Println("\nStep 2: Cloning your fork locally...")
	if !utils.DirExists(filepath.Dir(crossRepoPath)) {
		_ = os.MkdirAll(filepath.Dir(crossRepoPath), 0755)
	}

	if utils.DirExists(crossRepoPath) {
		fmt.Printf("Repository already exists at %s. Updating...\n", crossRepoPath)
		origWd, _ := os.Getwd()
		_ = os.Chdir(crossRepoPath)
		utils.RunCmd("git", "pull", "origin", "main")
		_ = os.Chdir(origWd)
	} else {
		utils.RunCmd("gh", "repo", "clone", githubUser+"/juice-cross", crossRepoPath)
	}

	fmt.Println("\nStep 3: Updating project configuration...")
	if utils.FileExists("config.ini") {
		// Update config.ini with these details
		// We use a simple sed or replacement logic here
		// In a real CLI, we might use a dedicated INI parser
		data, _ := os.ReadFile("config.ini")
		content := string(data)
		content = strings.Replace(content, "github_user =", "github_user = "+githubUser, 1)
		content = strings.Replace(content, "cross_repo_path =", "cross_repo_path = "+crossRepoPath, 1)
		_ = os.WriteFile("config.ini", []byte(content), 0644)
		fmt.Println("Successfully updated config.ini")
	} else {
		fmt.Println("Warning: config.ini not found in current directory. Configuration skipped.")
	}

	fmt.Println("\n=== Cross-compilation setup complete! ===")
}

// ManagePlugin handles the installation and removal of native Go/Mobile plugins.
// It manages go.mod dependencies and synchronizes native source files into the intermediate staging area.
func ManagePlugin(action, pluginRepo string) {
	if !utils.DirExists(".plugins") || !utils.FileExists("go.mod") {
		fmt.Fprintln(os.Stderr, "Error: You must execute plugin commands from the root of a wailsm project directory.")
		os.Exit(1)
	}

	if action == "add" {
		fmt.Fprintf(os.Stdout, "=== Installing Plugin: %s ===\n", pluginRepo)
		utils.RunCmd("go", "get", pluginRepo)

		out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", pluginRepo).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not resolve source location for module: %s\n", pluginRepo)
			os.Exit(1)
		}
		goModSrc := strings.TrimSpace(string(out))

		androidSrc := filepath.Join(goModSrc, "android")
		if utils.DirExists(androidSrc) {
			fmt.Println("Syncing Android native directory...")
			if err := utils.CopyDirectory(androidSrc, filepath.Join(".plugins", "android")); err != nil {
				utils.Fatal("Failed to sync Android native directory", err)
			}
		}

		iosSrc := filepath.Join(goModSrc, "ios")
		if utils.DirExists(iosSrc) {
			fmt.Println("Syncing iOS native directory...")
			if err := utils.CopyDirectory(iosSrc, filepath.Join(".plugins", "ios")); err != nil {
				utils.Fatal("Failed to sync iOS native directory", err)
			}
		}
		fmt.Fprintf(os.Stdout, "=== Plugin %s added successfully! ===\n", pluginRepo)
	} else if action == "remove" {
		fmt.Fprintf(os.Stdout, "=== Removing Plugin: %s ===\n", pluginRepo)
		pluginDirname := filepath.Base(pluginRepo)
		_ = os.RemoveAll(filepath.Join(".plugins", "android", pluginDirname))
		_ = os.RemoveAll(filepath.Join(".plugins", "ios", pluginDirname))

		utils.RunCmd("go", "mod", "edit", "-droprequire="+pluginRepo)
		utils.RunCmd("go", "mod", "tidy")
	}
}

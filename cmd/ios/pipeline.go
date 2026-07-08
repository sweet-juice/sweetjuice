// Package ios implements the iOS-specific build and deployment pipeline using xtool.
// It allows for iOS development and compilation on systems where xtool and related
// toolchains are natively installed.
package ios

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sweet-juice/sweetjuice/cmd/utils"
)

var (
	// CleanOutput determines if previous build artifacts should be purged.
	CleanOutput = "true"
	// IOSTarget defines the iOS architecture to build for.
	IOSTarget = "ios"
	// XFrameworkName is the filename of the generated iOS framework.
	XFrameworkName = "sweetjuice.xcframework"
)

func applyConfig() {
	config := utils.LoadConfig()
	target := config.GetOrDefault("build", "ios_target", "ios")
	if target != "ios" && !strings.Contains(target, "/") {
		IOSTarget = "ios/" + target
	} else {
		IOSTarget = target
	}
	XFrameworkName = config.GetOrDefault("ios", "xframework_name", "sweetjuice.xcframework")
}

// ValidateIOSEnvironment ensures the iOS build tools are available.
func ValidateIOSEnvironment() {
	if runtime.GOOS != "darwin" {
		fmt.Fprintln(os.Stderr, "[ios] Warning: iOS development is traditionally only supported on macOS. Proceeding as xtool is reported to be cross-platform.")
	}

	if !utils.CommandExists("xtool") {
		utils.Fatal("xtool missing", fmt.Errorf("please install xtool: https://github.com/mizage/xtool"))
	}

	if !utils.CommandExists("gomobile") {
		utils.Fatal("gomobile tool missing", fmt.Errorf("please run 'go install golang.org/x/mobile/cmd/gomobile@latest' and 'gomobile init'"))
	}
}

// ScaffoldProject uses xtool to create the initial iOS project structure.
func ScaffoldProject(name string) {
	fmt.Printf("Scaffolding iOS project '%s' using xtool...\n", name)

	iosDir := filepath.Join("native", "ios")
	if !utils.DirExists(iosDir) {
		_ = os.MkdirAll(iosDir, 0755)
	}

	// Step 1: Ensure xtool is available
	if !utils.CommandExists("xtool") {
		utils.Fatal("xtool is required for iOS project scaffolding", fmt.Errorf("xtool command not found"))
	}

	origWd, _ := os.Getwd()
	_ = os.Chdir(iosDir)
	defer func() { _ = os.Chdir(origWd) }()

	// Step 2: Run xtool init
	fmt.Println("Running xtool scaffolding sequence...")
	utils.RunCmd("xtool", "init", "--name", name)

	fmt.Println("iOS project scaffolded successfully.")
}

// RefreshPipeline triggers the native iOS build toolchain setup and Go binding generation.
func RefreshPipeline() {
	ValidateIOSEnvironment()
	applyConfig()
	fmt.Println("Starting iOS toolchain refresh and Go binding generation...")

	iosDir := filepath.Join("native", "ios")
	targetSrcDir := filepath.Join(iosDir, "Sources", "Plugins")
	stagingPluginsDir := filepath.Join(".plugins", "ios")
	outputPath := filepath.Join(iosDir, XFrameworkName)

	if !utils.DirExists(iosDir) {
		fmt.Fprintln(os.Stderr, "Error: Native iOS path layout missing. Ensure you have the iOS template files in native/ios.")
		os.Exit(1)
	}

	if CleanOutput == "true" {
		if utils.DirExists(outputPath) {
			fmt.Println("Cleaning previous Go bindings...")
			_ = os.RemoveAll(outputPath)
		}
	}

	// Step 1: Run gomobile bind
	fmt.Println("Generating Go bindings (XFramework) for iOS...")
	utils.RunCmd("gomobile", "bind", "-target="+IOSTarget, "-o", outputPath, ".")

	// Step 2: Sync plugins
	if utils.DirExists(stagingPluginsDir) && !utils.DirEmpty(stagingPluginsDir) {
		fmt.Println("Syncing iOS native plugins...")
		_ = os.MkdirAll(targetSrcDir, 0755)
		if err := utils.CopyDirectory(stagingPluginsDir, targetSrcDir); err != nil {
			utils.Fatal("Failed syncing plugin package trees inside iOS workspace", err)
		}
	}

	fmt.Printf("iOS Platform Refresh complete. Bindings generated in %s\n", outputPath)
}

// BuildPipeline triggers the compilation of the Swift/iOS application using xtool.
func BuildPipeline(mode string) {
	ValidateIOSEnvironment()
	applyConfig()
	fmt.Printf("Building iOS application in mode: %s via xtool...\n", mode)

	iosDir := filepath.Join("native", "ios")
	if !utils.DirExists(iosDir) {
		fmt.Fprintln(os.Stderr, "Error: Native iOS path layout missing.")
		os.Exit(1)
	}

	origWd, _ := os.Getwd()
	_ = os.Chdir(iosDir)
	defer func() { _ = os.Chdir(origWd) }()

	// Map generic modes to Xcode configurations
	config := "Debug"
	if mode == "release" || mode == "bundle" {
		config = "Release"
	}

	// We use 'xtool build' to compile the project
	fmt.Printf("Executing xtool build sequence (%s)...\n", config)
	utils.RunCmd("xtool", "build", "--configuration", config)

	fmt.Println("iOS application built successfully.")
}

// RunPipeline handles deployment to a real device via xtool.
func RunPipeline() {
	ValidateIOSEnvironment()
	fmt.Println("Preparing device deployment...")

	iosDir := filepath.Join("native", "ios")
	if !utils.DirExists(iosDir) {
		fmt.Fprintln(os.Stderr, "Error: Native iOS path layout missing.")
		os.Exit(1)
	}

	origWd, _ := os.Getwd()
	_ = os.Chdir(iosDir)
	defer func() { _ = os.Chdir(origWd) }()

	// Use 'xtool run' to deploy to the device.
	fmt.Println("Deploying to connected iOS device...")
	utils.RunCmd("xtool", "run")
}

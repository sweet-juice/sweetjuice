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
	XFrameworkName = "Sweetjuice.xcframework"
)

func applyConfig() {
	config := utils.LoadConfig()
	target := config.GetOrDefault("build", "ios_target", "ios")
	if target != "ios" && !strings.Contains(target, "/") {
		IOSTarget = "ios/" + target
	} else {
		IOSTarget = target
	}
	// We force the name to match Package.swift convention
	XFrameworkName = "Sweetjuice.xcframework"
}

// ValidateIOSEnvironment ensures the iOS build tools are available.
func ValidateIOSEnvironment() {
	if !utils.CommandExists("xtool") {
		utils.Fatal("xtool missing", fmt.Errorf("please install xtool: https://github.com/mizage/xtool"))
	}

	utils.EnsureGoMobileTools()
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

// RefreshPipeline triggers the native iOS build toolchain setup and Go compilation.
func RefreshPipeline() {
	ValidateIOSEnvironment()
	applyConfig()

	// Build frontend first
	utils.BuildFrontend()

	fmt.Println("Starting iOS toolchain refresh and Go compilation...")

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

	// Generate Go bindings (XCFramework) via gomobile bind
	fmt.Println("Generating Go bindings (XCFramework) for iOS via gomobile...")
	if runtime.GOOS != "darwin" {
		fmt.Println("[ios] Note: Local iOS binding on Linux/Windows requires a specialized gomobile toolchain.")
		fmt.Println("      It is recommended to use 'juice --run-cross ios' for GitHub Action based builds.")
	}

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

	// Map generic modes to xtool configurations
	config := "debug"
	if mode == "release" || mode == "bundle" {
		config = "release"
	}

	// We use 'xtool dev build' to compile the project
	fmt.Printf("Executing xtool build sequence (%s)...\n", config)
	utils.RunCmd("xtool", "dev", "build", "--configuration", config)

	fmt.Println("iOS application built successfully.")
}

// RunPipeline handles deployment to a real device via xtool.
func RunPipeline() {
	ValidateIOSEnvironment()
	fmt.Println("Preparing device deployment via xtool...")

	iosDir := filepath.Join("native", "ios")
	if !utils.DirExists(iosDir) {
		fmt.Fprintln(os.Stderr, "Error: Native iOS path layout missing.")
		os.Exit(1)
	}

	origWd, _ := os.Getwd()
	_ = os.Chdir(iosDir)
	defer func() { _ = os.Chdir(origWd) }()

	// Use 'xtool dev run' to deploy to the device.
	fmt.Println("Deploying to connected iOS device...")
	utils.RunCmd("xtool", "dev", "run")
}



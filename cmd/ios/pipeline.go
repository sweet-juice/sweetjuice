// Package ios implements the iOS-specific build and deployment pipeline using xtool.
// It allows for iOS development and compilation on systems where xtool and related
// toolchains are natively installed.
package ios

import (
	"fmt"
	"os"
	"os/exec"
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

// GetIOSSDKPath returns the path to the Darwin SDK provided by xtool.
func GetIOSSDKPath() string {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".swiftpm", "swift-sdks", "darwin.artifactbundle")
	if utils.DirExists(path) {
		return path
	}
	return ""
}

// ValidateIOSEnvironment ensures the iOS build tools are available.
func ValidateIOSEnvironment() {
	if runtime.GOOS != "darwin" {
		sdkPath := GetIOSSDKPath()
		if sdkPath == "" {
			utils.Fatal("iOS SDK missing", fmt.Errorf("could not locate xtool SDK at ~/.swiftpm/swift-sdks/darwin.artifactbundle. Please run 'xtool setup'"))
		}
		fmt.Printf("[ios] Using xtool SDK at %s\n", sdkPath)
	}

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

	// Step 1: Generate Go bindings
	if runtime.GOOS == "darwin" {
		fmt.Println("Generating Go bindings (XFramework) for iOS via gomobile...")
		utils.RunCmd("gomobile", "bind", "-target="+IOSTarget, "-o", outputPath, ".")
	} else {
		fmt.Println("Generating Go bindings (XFramework) for iOS via xtool cross-toolchain...")
		crossBindIOS(outputPath)
	}

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

func crossBindIOS(outputPath string) {
	sdkPath := GetIOSSDKPath()
	// Find iPhoneOS SDK root inside the artifact bundle
	iphoneSDKBase := filepath.Join(sdkPath, "Developer", "Platforms", "iPhoneOS.platform", "Developer", "SDKs")
	entries, _ := os.ReadDir(iphoneSDKBase)
	var iphoneSDK string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "iPhoneOS") && strings.HasSuffix(entry.Name(), ".sdk") {
			iphoneSDK = filepath.Join(iphoneSDKBase, entry.Name())
			break
		}
	}

	if iphoneSDK == "" {
		utils.Fatal("iOS SDK not found inside artifact bundle", fmt.Errorf("missing iPhoneOS.sdk"))
	}

	fmt.Printf("[ios] Using SDK: %s\n", iphoneSDK)

	// Bindings generation
	tempDir, _ := os.MkdirTemp("", "sweetjuice-ios-bind")
	defer os.RemoveAll(tempDir)

	fmt.Println("  -> Generating Objective-C and Go bindings...")
	// Generate bindings into the temp dir.
	// We use the current directory package.
	utils.RunCmd("gobind", "-lang=go,objc", "-outdir", tempDir, ".")

	// Get current module name
	modName := "project"
	if data, err := os.ReadFile("go.mod"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "module ") {
				modName = strings.TrimSpace(strings.TrimPrefix(line, "module "))
				break
			}
		}
	}

	// Find the actual path of sweetjuice module
	sweetjuicePath := ""
	if out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/sweet-juice/sweetjuice").Output(); err == nil {
		sweetjuicePath = strings.TrimSpace(string(out))
	}

	// Create a temporary main package to build the c-archive
	buildDir := filepath.Join(tempDir, "build")
	_ = os.MkdirAll(buildDir, 0755)

	// Copy the generated Go and Header files directly into our main package
	gobindSrcDir := filepath.Join(tempDir, "src", "gobind")
	files, _ := os.ReadDir(gobindSrcDir)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".go") {
			data, _ := os.ReadFile(filepath.Join(gobindSrcDir, file.Name()))
			// Change package from 'gobind' to 'main' so it can be compiled together
			content := strings.Replace(string(data), "package gobind", "package main", 1)

			// If it's seq.go, add our required side-effect imports
			if file.Name() == "seq.go" {
				extraImports := fmt.Sprintf("\nimport _ \"golang.org/x/mobile/bind\"\nimport _ \"%s\"\n", modName)
				content = strings.Replace(content, "package main", "package main"+extraImports, 1)
			}

			_ = os.WriteFile(filepath.Join(buildDir, file.Name()), []byte(content), 0644)
		} else if strings.HasSuffix(file.Name(), ".h") || strings.HasSuffix(file.Name(), ".m") {
			// Copy headers and implementation files as well for cgo to find and compile them
			data, _ := os.ReadFile(filepath.Join(gobindSrcDir, file.Name()))
			_ = os.WriteFile(filepath.Join(buildDir, file.Name()), []byte(data), 0644)
		}
	}

	// Create a temporary go.mod for the build
	absPath, _ := filepath.Abs(".")
	var goModLines []string
	goModLines = append(goModLines, "module build", "go 1.24", "", "require (")
	goModLines = append(goModLines, fmt.Sprintf("\t%s v0.0.0", modName))
	if sweetjuicePath != "" {
		goModLines = append(goModLines, "\tgithub.com/sweet-juice/sweetjuice v0.0.0")
	}
	goModLines = append(goModLines, ")", "", fmt.Sprintf("replace %s => %s", modName, absPath))
	if sweetjuicePath != "" {
		goModLines = append(goModLines, fmt.Sprintf("replace github.com/sweet-juice/sweetjuice => %s", sweetjuicePath))
	}
	_ = os.WriteFile(filepath.Join(buildDir, "go.mod"), []byte(strings.Join(goModLines, "\n")), 0644)

	// Compile Go to static library
	fmt.Println("  -> Compiling Go for ios/arm64...")
	swiftBin := "/usr/lib/swift/bin"
	llvmRanlib := filepath.Join(swiftBin, "llvm-ranlib")
	llvmAr := filepath.Join(swiftBin, "llvm-ar")

	origWd, _ := os.Getwd()
	_ = os.Chdir(buildDir)
	defer func() { _ = os.Chdir(origWd) }()

	fmt.Println("  -> Resolving dependencies...")
	utils.RunCmd("go", "mod", "tidy")

	// Use a unique name for the core library to avoid header conflicts with gobind
	cmd := exec.Command("go", "build", "-buildmode=c-archive", "-tags", "ios", "-o", "SweetjuiceCore.a", ".")
	cmd.Env = append(os.Environ(),
		"PATH="+swiftBin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"GOOS=ios",
		"GOARCH=arm64",
		"CGO_ENABLED=1",
		"AR="+llvmAr,
		"CC=clang -target arm64-apple-ios17.0 -isysroot "+iphoneSDK,
		"CXX=clang++ -target arm64-apple-ios17.0 -isysroot "+iphoneSDK,
		"CGO_CFLAGS=-I. -D__GOBIND_DARWIN__",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		utils.Fatal("Go compilation for iOS failed", err)
	}

	// Create XFramework structure
	frameworkDir := filepath.Join(origWd, outputPath, "ios-arm64", "Sweetjuice.framework")
	_ = os.MkdirAll(filepath.Join(frameworkDir, "Headers"), 0755)
	_ = os.MkdirAll(filepath.Join(frameworkDir, "Modules"), 0755)

	// Copy binary. IMPORTANT: Rename to match framework name
	targetBinary := filepath.Join(frameworkDir, "Sweetjuice")
	if err := utils.MoveFile(filepath.Join(buildDir, "SweetjuiceCore.a"), targetBinary); err != nil {
		utils.Fatal("Failed to move compiled library to framework", err)
	}

	// Add symbol index using the correct LLVM tool for Apple targets
	fmt.Println("  -> Indexing framework binary...")
	utils.RunCmd(llvmRanlib, targetBinary)

	// Copy ALL headers from gobind
	files, _ = os.ReadDir(gobindSrcDir)
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".h") {
			data, _ := os.ReadFile(filepath.Join(gobindSrcDir, name))
			content := string(data)
			// Force __GOBIND_DARWIN__ in seq.h to fix nstring errors
			if name == "seq.h" {
				content = "#ifndef __GOBIND_DARWIN__\n#define __GOBIND_DARWIN__\n#endif\n" + content
			}
			_ = os.WriteFile(filepath.Join(frameworkDir, "Headers", name), []byte(content), 0644)
		}
	}

	// Copy and patch the generated C core header
	coreHeaderPath := filepath.Join(buildDir, "SweetjuiceCore.h")
	if !utils.FileExists(coreHeaderPath) {
		// Fallback to lowercase if needed
		coreHeaderPath = filepath.Join(buildDir, "sweetjuicecore.h")
	}

	if data, err := os.ReadFile(coreHeaderPath); err == nil {
		coreHeaderContent := string(data)
		// Change include to avoid confusion with the gobind sweetjuice.h
		coreHeaderContent = strings.ReplaceAll(coreHeaderContent, "#include \"sweetjuice.h\"", "#include \"Sweetjuice.objc.h\"")
		_ = os.WriteFile(filepath.Join(frameworkDir, "Headers", "SweetjuiceCore.h"), []byte(coreHeaderContent), 0644)
	}

	// Create Umbrella Header
	umbrellaHeader := ""
	headers, _ := os.ReadDir(filepath.Join(frameworkDir, "Headers"))
	for _, h := range headers {
		name := h.Name()
		// Only include top-level headers, skip the umbrella itself
		if strings.HasSuffix(name, ".h") && name != "Sweetjuice.h" {
			umbrellaHeader += fmt.Sprintf("#import \"%s\"\n", name)
		}
	}
	_ = os.WriteFile(filepath.Join(frameworkDir, "Headers", "Sweetjuice.h"), []byte(umbrellaHeader), 0644)

	// Create Module Map
	moduleMap := `framework module Sweetjuice {
  umbrella header "Sweetjuice.h"
  export *
  module * { export * }
}`
	_ = os.WriteFile(filepath.Join(frameworkDir, "Modules", "module.modulemap"), []byte(moduleMap), 0644)

	// Create Info.plist for framework
	plistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>Sweetjuice</string>
	<key>CFBundleIdentifier</key>
	<string>com.sweetjuice.core</string>
	<key>CFBundleInfoDictionaryVersion</key>
	<string>6.0</string>
	<key>CFBundleName</key>
	<string>Sweetjuice</string>
	<key>CFBundlePackageType</key>
	<string>FMWK</string>
	<key>CFBundleShortVersionString</key>
	<string>1.0</string>
	<key>CFBundleVersion</key>
	<string>1.0</string>
	<key>CFBundleSupportedPlatforms</key>
	<array>
		<string>iPhoneOS</string>
	</array>
	<key>MinimumOSVersion</key>
	<string>17.0</string>
</dict>
</plist>`
	_ = os.WriteFile(filepath.Join(frameworkDir, "Info.plist"), []byte(plistContent), 0644)

	// Create Info.plist for XCFramework
	xcPlistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AvailableLibraries</key>
	<array>
		<dict>
			<key>LibraryIdentifier</key>
			<string>ios-arm64</string>
			<key>LibraryPath</key>
			<string>Sweetjuice.framework</string>
			<key>SupportedArchitectures</key>
			<array>
				<string>arm64</string>
			</array>
			<key>SupportedPlatform</key>
			<string>ios</string>
		</dict>
	</array>
	<key>CFBundlePackageType</key>
	<string>XFWK</string>
	<key>XCFrameworkFormatVersion</key>
	<string>1.0</string>
</dict>
</plist>`
	_ = os.WriteFile(filepath.Join(origWd, outputPath, "Info.plist"), []byte(xcPlistContent), 0644)

	fmt.Println("  -> XCFramework assembled successfully with module map.")
}

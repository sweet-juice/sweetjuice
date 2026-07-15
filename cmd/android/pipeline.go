// Package android implements the Android-specific build and deployment pipeline.
// It manages Gradle compilation, AAR generation via gomobile, and ADB device orchestration.
package android

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sweet-juice/sweetjuice/cmd/utils"
)

var (
	// GomobileTarget defines the Android ABI architecture to build for.
	GomobileTarget = "android/arm64"
	// AndroidAPI specifies the minimum Android SDK level for gomobile bind.
	AndroidAPI = "21"
	// AarName is the filename of the generated Android archive.
	AarName = "sweetjuice.aar"
	// CleanOutput determines if previous build artifacts should be purged.
	CleanOutput = "true"
)

// AndroidManifest represents a subset of the Android XML manifest for metadata extraction.
type AndroidManifest struct {
	XMLName     xml.Name    `xml:"manifest"`
	PackageName string      `xml:"package,attr"`
	Application Application `xml:"application"`
}

// Application represents the application tag in the Android manifest.
type Application struct {
	Activities []Activity `xml:"activity"`
}

// Activity represents an activity tag in the Android manifest.
type Activity struct {
	Name          string         `xml:"http://schemas.android.com/apk/res/android name,attr"`
	IntentFilters []IntentFilter `xml:"intent-filter"`
}

// IntentFilter represents an intent-filter tag in the Android manifest.
type IntentFilter struct {
	Actions    []Action   `xml:"action"`
	Categories []Category `xml:"category"`
}

// Action represents an action tag in the Android manifest.
type Action struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

// Category represents a category tag in the Android manifest.
type Category struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

// SetupAndroidLocalProperties discovers the Android SDK path and writes local.properties.
func SetupAndroidLocalProperties(targetDir string) {
	sdkPath := GetAndroidSDKPath()
	if sdkPath != "" {
		androidDir := filepath.Join(targetDir, "native", "android")
		propsPath := filepath.Join(androidDir, "local.properties")
		if !utils.DirExists(androidDir) {
			_ = os.MkdirAll(androidDir, 0755)
		}
		content := fmt.Sprintf("sdk.dir=%s\n", filepath.ToSlash(sdkPath))
		_ = os.WriteFile(propsPath, []byte(content), 0644)
		fmt.Printf("[android] SDK found at %s. Updated local.properties\n", sdkPath)
	} else {
		fmt.Fprintln(os.Stderr, "[android] Warning: Could not locate Android SDK automatically. Please set ANDROID_HOME.")
	}
}

// GetAndroidSDKPath returns the path to the Android SDK, checking environment variables and default install locations.
func GetAndroidSDKPath() string {
	// 1. Check environment variables
	if path := os.Getenv("ANDROID_HOME"); path != "" && utils.DirExists(path) {
		return path
	}
	if path := os.Getenv("ANDROID_SDK_ROOT"); path != "" && utils.DirExists(path) {
		return path
	}

	// 2. Check default platform-specific installation paths (Android Studio defaults)
	home, _ := os.UserHomeDir()
	var paths []string

	switch runtime.GOOS {
	case "darwin":
		paths = []string{filepath.Join(home, "Library", "Android", "sdk")}
	case "linux":
		paths = []string{
			filepath.Join(home, "Android", "Sdk"),
			filepath.Join(home, "android-sdk"),
			"/usr/lib/android-sdk",
			"/opt/android-sdk",
		}
	case "windows":
		paths = []string{
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Android", "Sdk"),
		}
	}

	for _, p := range paths {
		if utils.DirExists(p) {
			return p
		}
	}

	return ""
}

// ValidateAndroidEnvironment ensures the Android SDK and required tools are available.
func ValidateAndroidEnvironment() {
	sdkPath := GetAndroidSDKPath()
	if sdkPath == "" {
		utils.Fatal("Android SDK not found", fmt.Errorf("could not locate SDK in environment or default locations. Please install Android Studio or set ANDROID_HOME"))
	}

	// Ensure local.properties exists in the current project context
	if utils.DirExists(filepath.Join("native", "android")) {
		SetupAndroidLocalProperties(".")
	}

	utils.EnsureGoMobileTools()
}

func applyConfig() {
	config := utils.LoadConfig()
	target := config.GetOrDefault("build", "android_target", "arm64")
	if !strings.Contains(target, "/") {
		GomobileTarget = "android/" + target
	} else {
		GomobileTarget = target
	}
	AndroidAPI = config.GetOrDefault("android", "min_api", "21")
	AarName = config.GetOrDefault("android", "aar_name", "sweetjuice.aar")
}

// RefreshPipeline runs the gomobile bind command to sync Go code with the Android project.
// It also synchronizes native plugin source files.
func RefreshPipeline() {
	ValidateAndroidEnvironment()
	applyConfig()

	// Build frontend first
	utils.BuildFrontend()

	outputPath := filepath.Join("native", "android", "app", "libs")
	targetJavaSrcDir := filepath.Join("native", "android", "app", "src", "main", "java")
	stagingPluginsDir := filepath.Join(".plugins", "android")

	_ = os.MkdirAll(outputPath, 0755)

	if CleanOutput == "true" {
		files, _ := filepath.Glob(filepath.Join(outputPath, "*.aar"))
		for _, f := range files {
			_ = os.Remove(f)
		}
	}

	fmt.Fprintf(os.Stdout, "Building %s for target %s...\n", AarName, GomobileTarget)
	utils.RunCmd("gomobile", "bind", "-target="+GomobileTarget, "-androidapi="+AndroidAPI, "-o", filepath.Join(outputPath, AarName), ".")

	if utils.DirExists(stagingPluginsDir) && !utils.DirEmpty(stagingPluginsDir) {
		_ = os.MkdirAll(targetJavaSrcDir, 0755)
		if err := utils.CopyDirectory(stagingPluginsDir, targetJavaSrcDir); err != nil {
			utils.Fatal("Failed syncing plugin package trees inside Android workspace", err)
		}
	}
}

// BuildPipeline executes the Gradle build process to generate APK or AAB files.
func BuildPipeline(mode string) {
	ValidateAndroidEnvironment()
	applyConfig()
	androidDir := filepath.Join("native", "android")
	if !utils.DirExists(androidDir) {
		fmt.Fprintln(os.Stderr, "Error: Native Android path layout missing.")
		os.Exit(1)
	}

	gradleCmd := "./gradlew"
	if utils.IsWindowsHost() {
		gradleCmd = "gradlew.bat"
	}

	origWd, _ := os.Getwd()
	_ = os.Chdir(androidDir)
	defer func() { _ = os.Chdir(origWd) }()

	var targetTask string
	switch mode {
	case "debug":
		targetTask = "assembleDebug"
	case "release":
		targetTask = "assembleRelease"
	case "bundle":
		targetTask = "bundleRelease"
	}

	utils.RunCmd(gradleCmd, targetTask)
}

// RunPipeline installs the compiled APK to a connected device via ADB and launches it.
func RunPipeline() {
	ValidateAndroidEnvironment()

	adbTool := "adb"
	if _, err := exec.LookPath("adb"); err != nil {
		// Try to find it in the SDK platform-tools
		sdkPath := GetAndroidSDKPath()
		adbPath := filepath.Join(sdkPath, "platform-tools", "adb")
		if runtime.GOOS == "windows" {
			adbPath += ".exe"
		}
		if _, err := os.Stat(adbPath); err == nil {
			adbTool = adbPath
		} else {
			fmt.Fprintln(os.Stderr, "Error: 'adb' tool missing and not found in SDK platform-tools.")
			os.Exit(1)
		}
	}

	apkPath := filepath.Join("native", "android", "app", "build", "outputs", "apk", "debug", "app-debug.apk")
	if !utils.FileExists(apkPath) {
		fmt.Fprintf(os.Stderr, "Error: APK not found at %s. Did you build it?\n", apkPath)
		os.Exit(1)
	}

	fmt.Println("Installing APK to device...")
	utils.RunCmd(adbTool, "install", "-r", apkPath)

	manifestPath := filepath.Join("native", "android", "app", "src", "main", "AndroidManifest.xml")
	packageName, launcherActivity, err := parseManifestDetails(manifestPath)
	if err != nil {
		return
	}

	// Fallback for packageName from config if not in manifest
	if packageName == "" {
		config := utils.LoadConfig()
		packageName = config.GetOrDefault("app", "package", "com.sweetjuice.app")
	}

	// Ensure launcher activity is fully qualified if it starts with a dot
	if strings.HasPrefix(launcherActivity, ".") {
		launcherActivity = packageName + launcherActivity
	}

	fmt.Printf("Launching application %s/%s...\n", packageName, launcherActivity)
	utils.RunCmd(adbTool, "shell", "am", "start", "-n", packageName+"/"+launcherActivity)
}

func parseManifestDetails(manifestPath string) (string, string, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", "", err
	}
	var manifest AndroidManifest
	if err := xml.Unmarshal(data, &manifest); err != nil {
		return "", "", err
	}

	packageName := manifest.PackageName
	launcherActivity := ""
	for _, act := range manifest.Application.Activities {
		for _, filter := range act.IntentFilters {
			for _, action := range filter.Actions {
				if action.Name == "android.intent.action.MAIN" {
					launcherActivity = act.Name
					break
				}
			}
		}
	}
	return packageName, launcherActivity, nil
}

package sweetjuice

import (
	"embed"
	"fmt"
	"time"

	"github.com/sweet-juice/sweetjuice/plugins/biometric"
	"github.com/sweet-juice/sweetjuice/plugins/daemon"
	"github.com/sweet-juice/sweetjuice/plugins/filepicker"
	"github.com/sweet-juice/sweetjuice/plugins/logger"
	"github.com/sweet-juice/sweetjuice/plugins/notification"
	"github.com/sweet-juice/sweetjuice/plugins/permission"
	"github.com/sweet-juice/sweetjuice/plugins/workmanager"
	"github.com/sweet-juice/sweetjuice/core"
)

//go:embed frontend/*
var assets embed.FS

// Global plugin instances
var (
	permPlugin      = permission.NewPlugin()
	wmPlugin        = workmanager.NewPlugin()
	notifyPlugin    = notification.NewPlugin()
	biometricPlugin = biometric.NewPlugin()
	filePicker      = filepicker.NewPlugin()
	daemonPlugin    = daemon.NewPlugin()
	log             = logger.NewPlugin("[Sweet Juice]")
)

// NativeCallHandler is an interface that Java can implement to handle calls from Go.
// Defining it here ensures gomobile generates it in the 'sweetjuice' Java package.
type NativeCallHandler interface {
	OnNativeCall(method string, args string) string
}

// SetNativeCallHandler registers the Java-side handler for Go-to-Native calls.
func SetNativeCallHandler(handler NativeCallHandler) {
	core.SetNativeCallHandler(handler)
}

// StartApplication initializes the mobile backend runtime from Android.
func StartApplication() string {
	helloService := NewAppService()

	app := core.NewApplication(core.Options{
		Name:   "Sweet Juice",
		Assets: assets,
		Bind: []interface{}{
			helloService, // Always bind the main application service
			permPlugin,
			wmPlugin,
			notifyPlugin,
			biometricPlugin,
			filePicker,
			daemonPlugin,
		},
		OnStart: func(app *core.Application) error {
			// Initialize plugins
			if err := log.Init(app); err != nil {
				return err
			}
			if err := permPlugin.Init(app); err != nil {
				return err
			}
			if err := wmPlugin.Init(app); err != nil {
				return err
			}
			if err := notifyPlugin.Init(app); err != nil {
				return err
			}
			if err := biometricPlugin.Init(app); err != nil {
				return err
			}
			if err := filePicker.Init(app); err != nil {
				return err
			}
			if err := daemonPlugin.Init(app); err != nil {
				return err
			}

			log.Info("Go Application started!")

			// Example: Register a background task
			wmPlugin.RegisterTask("sync_analytics", func() error {
				log.Info("Background task sync_analytics invoked")
				notifyPlugin.Post(notification.Notification{
					ID:         100,
					Title:      "Periodic Post",
					Body:       "Hello from Sweet Juice!",
					Importance: notification.ImportanceHigh,
				})
				return nil
			})

			// Schedule the background task to run every 15 minutes after a 30-second delay
			time.AfterFunc(30*time.Second, func() {
				if status, err := permPlugin.Check("android.permission.POST_NOTIFICATIONS"); status != "granted" {
					log.Error("Permission check failed: %v", err)
					log.Info("Status: %s", status)
				} else {
					log.Debug("Enqueuing periodic background work...")
					wmPlugin.EnqueuePeriodic("sync_analytics", 15, &workmanager.Constraints{
						NetworkType: workmanager.NetworkNotRequired,
					})
				}
			})

			return nil
		},
	})

	if err := app.Run(); err != nil {
		log.Error("Application failed to run: %v", err)
		return fmt.Sprintf(`{"error":"%s"}`, err.Error())
	}

	return `{"status":"started"}`
}

// Below functions are called from Java to handle messages/events from the frontend or to perform native actions.

// HandleMessageFromFrontend processes messages sent from the JavaScript frontend.
func HandleMessageFromFrontend(methodKey string, jsonArgsPayload string) string {
	return core.HandleMessageFromFrontend(methodKey, jsonArgsPayload)
}

// HandleNativeAction processes calls from Go to Java and returns the result back to Go.
func HandleNativeAction(methodKey string, jsonArgsPayload string) string {
	return core.HandleNativeAction(methodKey, jsonArgsPayload)
}

func RequestAssetBytes(urlPath string) []byte {
	return core.NewMobileBridge().RequestAssetBytes(urlPath)
}

// RequestAssetMime retrieves the MIME type for a given asset path.
func RequestAssetMime(urlPath string) string {
	return core.NewMobileBridge().RequestAssetMime(urlPath)
}

// PollNativeEvent allows Go to check for any events sent from Java and retrieve them as a JSON string.
func PollNativeEvent() string {
	return core.NewMobileBridge().PollNativeEvent()
}

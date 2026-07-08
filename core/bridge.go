package core

import (
	"encoding/json"
	"fmt"
)

var globalAppInstance *Application

// NativeCallHandler is an interface implemented by the mobile platform (Java/Obj-C)
// to handle calls originating from the Go layer.
type NativeCallHandler interface {
	OnNativeCall(method string, args string) string
}

var globalNativeHandler NativeCallHandler

// SetNativeCallHandler registers the platform-specific handler.
func SetNativeCallHandler(handler NativeCallHandler) {
	globalNativeHandler = handler
}

// CallNativePlatform calls the registered native handler from Go.
func CallNativePlatform(method string, args string) string {
	if globalNativeHandler == nil {
		return `{"error": "No native handler registered"}`
	}
	return globalNativeHandler.OnNativeCall(method, args)
}

func SetGlobalApp(app *Application) {
	globalAppInstance = app
}

// MobileBridge is a dedicated wrapper struct that gobind can parse into a Java Class.
type MobileBridge struct{}

func NewMobileBridge() *MobileBridge {
	return &MobileBridge{}
}

// CallGoBackend provides a class-mapped method for the JNI layer to invoke.
func (b *MobileBridge) CallGoBackend(methodKey string, jsonArgsPayload string) string {
	if globalAppInstance == nil {
		return `{"error": "Application core runtime context not active"}`
	}

	var rawArgs []json.RawMessage
	if err := json.Unmarshal([]byte(jsonArgsPayload), &rawArgs); err != nil {
		return fmt.Sprintf(`{"error": "Failed to extract arguments payload: %s"}`, err.Error())
	}

	result, err := globalAppInstance.InvokeCall(methodKey, rawArgs)
	if err != nil {
		return fmt.Sprintf(`{"error": %q}`, err.Error())
	}

	responsePayload, _ := json.Marshal(map[string]interface{}{
		"result": result,
	})
	return string(responsePayload)
}

func (b *MobileBridge) PollNativeEvent() string {
	if globalAppInstance == nil {
		return ""
	}
	return globalAppInstance.Events.PollNativeEvent()
}

// HandleMessageFromFrontend is a package-level helper exposed to gomobile-generated Java wrappers.
func HandleMessageFromFrontend(methodKey string, jsonArgsPayload string) string {
	bridge := NewMobileBridge()
	return bridge.CallGoBackend(methodKey, jsonArgsPayload)
}

// HandleNativeAction is a package-level helper exposed to Java plugin packages.
func HandleNativeAction(methodKey string, jsonArgsPayload string) string {
	if globalAppInstance == nil {
		return `{"error": "Application core runtime context not active"}`
	}

	var rawArgs []json.RawMessage
	if err := json.Unmarshal([]byte(jsonArgsPayload), &rawArgs); err != nil {
		return fmt.Sprintf(`{"error": "Failed to extract arguments payload: %s"}`, err.Error())
	}

	result, err := globalAppInstance.InvokeNativeCall(methodKey, rawArgs)
	if err != nil {
		return fmt.Sprintf(`{"error": %q}`, err.Error())
	}

	responsePayload, _ := json.Marshal(map[string]interface{}{
		"result": result,
	})
	return string(responsePayload)
}

// RequestAssetBytes acts as a direct memory array provider for the native layout
func (b *MobileBridge) RequestAssetBytes(urlPath string) []byte {
	if globalAppInstance == nil {
		return []byte("Asset Server Offline")
	}
	return globalAppInstance.ReadAsset(urlPath).Data
}

// RequestAssetMime returns correct content headers to the WebView controller
func (b *MobileBridge) RequestAssetMime(urlPath string) string {
	if globalAppInstance == nil {
		return "text/plain"
	}
	return globalAppInstance.ReadAsset(urlPath).MimeType
}

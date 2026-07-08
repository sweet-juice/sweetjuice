// Package sweetjuice provides the core runtime for Sweet Juice applications.
// It manages the application lifecycle, service binding, and the bridge between
// Go and the native mobile platform.
package core

import (
	"embed"
	"encoding/json"
	"fmt"
)

// Options defines the application configuration following the Sweet Juice v3 pattern.
type Options struct {
	// Name is the name of the application.
	Name string
	// Assets is the embedded filesystem containing the frontend assets.
	Assets embed.FS
	// Bind is a slice of struct instances to export to the frontend.
	Bind []interface{}
	// OnStart is a callback triggered when the application has initialized.
	OnStart func(app *Application) error
}

// WindowOptions defines initial window configuration.
type WindowOptions struct {
	Title string
}

// Window represents the application window (WebView container).
type Window struct {
	Options WindowOptions
}

type nativeMethod func([]json.RawMessage) (interface{}, error)

// Application represents the global lifecycle state container.
// It manages the communication between the Go backend and the native mobile environment.
type Application struct {
	Name          string
	Window        *Window
	methods       map[string]interface{} // Store as generic interfaces to hide reflect types from gobind
	nativeMethods map[string]nativeMethod
	Events        *EventBus
	options       Options
}

// NewApplication creates a new instance of the Sweet Juice application.
func NewApplication(options Options) *Application {
	return &Application{
		Name:          options.Name,
		Window:        &Window{},
		methods:       make(map[string]interface{}),
		nativeMethods: make(map[string]nativeMethod),
		Events:        NewEventBus(),
		options:       options,
	}
}

// NewWindow configures and returns the application window.
func (a *Application) NewWindow(opts WindowOptions) *Window {
	a.Window.Options = opts
	return a.Window
}

// Run starts the application engine and initializes bindings.
func (a *Application) Run() error {
	fmt.Printf("[%s] Initializing Sweet Juice core engine...\n", a.Name)
	if err := a.parseBindings(); err != nil {
		return fmt.Errorf("failed to parse structural bindings: %w", err)
	}

	if a.nativeMethods == nil {
		a.nativeMethods = make(map[string]nativeMethod)
	}

	SetGlobalApp(a)

	if a.options.OnStart != nil {
		return a.options.OnStart(a)
	}

	return nil
}

// RegisterNativeMethod registers a Go function as a "Native Method" that can be called
// directly from the Java/Objective-C layer using the HandleNativeAction bridge.
func (a *Application) RegisterNativeMethod(methodKey string, fn func([]json.RawMessage) (interface{}, error)) {
	if a.nativeMethods == nil {
		a.nativeMethods = make(map[string]nativeMethod)
	}
	a.nativeMethods[methodKey] = fn
}

// InvokeNativeCall executes a previously registered native method.
func (a *Application) InvokeNativeCall(methodKey string, rawArgs []json.RawMessage) (interface{}, error) {
	if fn, exists := a.nativeMethods[methodKey]; exists {
		return fn(rawArgs)
	}
	return nil, fmt.Errorf("native method identity '%s' not registered with application", methodKey)
}

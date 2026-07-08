package filepicker

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

type FilePickerPlugin struct {
	app *core.Application
}

type PickerOptions struct {
	MimeType string `json:"mime_type"`
	Multiple bool   `json:"multiple"`
}

type PickerResult struct {
	URI      string   `json:"uri,omitempty"`
	URIs     []string `json:"uris,omitempty"`
	Multiple bool     `json:"multiple"`
	Error    string   `json:"error,omitempty"`
}

// NewPlugin creates a new instance of the FilePicker plugin.
func NewPlugin() *FilePickerPlugin {
	return &FilePickerPlugin{}
}

// Name returns the plugin name.
func (p *FilePickerPlugin) Name() string {
	return "filepicker"
}

// Init initializes the plugin and registers the native callback for picker results.
func (p *FilePickerPlugin) Init(app *core.Application) error {
	p.app = app

	app.RegisterNativeMethod("filepicker:result", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no result provided")
		}

		var result PickerResult
		if err := json.Unmarshal(args[0], &result); err != nil {
			return nil, err
		}

		app.Events.Emit("filepicker:result", result)
		return map[string]string{"status": "received"}, nil
	})

	return nil
}

// PickFile triggers the native file picker.
// Results are delivered asynchronously via the "filepicker:result" event.
func (p *FilePickerPlugin) PickFile(options PickerOptions) (string, error) {
	if options.MimeType == "" {
		options.MimeType = "*/*"
	}
	payload, _ := json.Marshal(options)
	return core.CallNativePlatform("filepicker:pickFile", string(payload)), nil
}

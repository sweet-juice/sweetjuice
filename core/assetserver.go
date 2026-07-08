package core

import (
	"fmt"
	"io"
	"path"
	"strings"
)

// AssetResponse wraps file reading results for the native client container.
type AssetResponse struct {
	Data       []byte
	MimeType   string
	StatusCode int
}

// ReadAsset extracts content directly out of the developer's embedded workspace.
func (a *Application) ReadAsset(urlPath string) *AssetResponse {
	// Clean up paths arriving from the WebView container
	cleanPath := path.Clean(strings.TrimPrefix(urlPath, "/"))
	if cleanPath == "." || cleanPath == "" {
		cleanPath = "frontend/index.html"
	} else if !strings.HasPrefix(cleanPath, "frontend/") {
		cleanPath = "frontend/" + cleanPath
	}

	file, err := a.options.Assets.Open(cleanPath)
	if err != nil {
		return &AssetResponse{
			Data:       []byte(fmt.Sprintf("Asset not found: %s", cleanPath)),
			MimeType:   "text/plain",
			StatusCode: 404,
		}
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return &AssetResponse{
			Data:       []byte("Internal asset extraction failure"),
			MimeType:   "text/plain",
			StatusCode: 500,
		}
	}

	return &AssetResponse{
		Data:       data,
		MimeType:   getMimeType(cleanPath),
		StatusCode: 200,
	}
}

func getMimeType(filePath string) string {
	switch path.Ext(filePath) {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

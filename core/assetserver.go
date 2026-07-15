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
		cleanPath = "index.html"
	}

	// Potential locations to search for the asset
	searchPaths := []string{
		cleanPath,
		path.Join("frontend", cleanPath),
		path.Join("frontend", "dist", cleanPath),
		path.Join("dist", cleanPath),
	}

	var data []byte
	var finalPath string
	found := false

	for _, p := range searchPaths {
		file, err := a.options.Assets.Open(p)
		if err == nil {
			data, _ = io.ReadAll(file)
			file.Close()
			finalPath = p
			found = true
			break
		}
	}

	if !found {
		return &AssetResponse{
			Data:       []byte(fmt.Sprintf("Asset not found: %s", cleanPath)),
			MimeType:   "text/plain",
			StatusCode: 404,
		}
	}

	return &AssetResponse{
		Data:       data,
		MimeType:   getMimeType(finalPath),
		StatusCode: 200,
	}
}

func getMimeType(filePath string) string {
	switch path.Ext(filePath) {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js", ".mjs":
		return "text/javascript"
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

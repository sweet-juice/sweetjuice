package utils

import (
	"bufio"
	"os"
	"strings"
)

// Config represents the parsed config.ini data.
type Config struct {
	App     map[string]string
	Build   map[string]string
	Android map[string]string
	IOS     map[string]string
}

// LoadConfig reads and parses config.ini from the current directory.
func LoadConfig() *Config {
	config := &Config{
		App:     make(map[string]string),
		Build:   make(map[string]string),
		Android: make(map[string]string),
		IOS:     make(map[string]string),
	}

	file, err := os.Open("config.ini")
	if err != nil {
		return config // Return default/empty if not found
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentSection string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.ToLower(line[1 : len(line)-1])
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch currentSection {
		case "app":
			config.App[key] = value
		case "build":
			config.Build[key] = value
		case "android":
			config.Android[key] = value
		case "ios":
			config.IOS[key] = value
		}
	}

	return config
}

// GetOrDefault returns the value for a key in a section, or a default value if not found.
func (c *Config) GetOrDefault(section, key, defaultValue string) string {
	var m map[string]string
	switch strings.ToLower(section) {
	case "app":
		m = c.App
	case "build":
		m = c.Build
	case "android":
		m = c.Android
	case "ios":
		m = c.IOS
	}

	if val, ok := m[key]; ok && val != "" {
		return val
	}
	return defaultValue
}

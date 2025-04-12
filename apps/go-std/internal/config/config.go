package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// Config provides a unified interface for both JSONC and environment variables
type ConfigMap struct {
	// Cache commonly used values
	cachedValues struct {
		port   int
		debug  bool
		apiKey string
		// etc...
	}
	// Raw data
	data map[string]interface{}
}

var (
	instance *ConfigMap
	once     sync.Once
)

func findEnvFile() (string, error) {
	// Possible locations for .env file
	locations := []string{
		".",           // Current directory
		"..",          // Parent directory
		"../..",       // Parent's parent directory
		"../../..",    // Root directory
		"apps/go-std", // Go app directory
	}

	for _, loc := range locations {
		path := filepath.Join(loc, ".env")
		if _, err := os.Stat(path); err == nil {
			abs, err := filepath.Abs(path)
			if err != nil {
				continue
			}
			return abs, nil
		}
	}

	return "", fmt.Errorf(".env file not found in any of the expected locations")
}

// NewConfig initializes configuration from JSONC file
func Config() (*ConfigMap, error) {
	var initErr error

	once.Do(func() {
		fmt.Println("Initializing config")
		instance = &ConfigMap{
			data: make(map[string]interface{}),
		}

		envPath, err := findEnvFile()
		if err != nil {
			initErr = fmt.Errorf("error finding .env file: %w", err)
			return
		}

		fmt.Printf("Loading .env from: %s\n", envPath)
		if err := godotenv.Load(envPath); err != nil {
			initErr = fmt.Errorf("error loading .env file: %w", err)
			return
		}

		// Cache frequently used values from env
		if port := os.Getenv("port"); port != "" {
			if p, err := strconv.Atoi(port); err == nil {
				instance.cachedValues.port = p
			}
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return instance, nil
}

// LoadJSON loads and merges JSON configuration on top of environment variables
func (c *ConfigMap) LoadJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read JSON config file: %w", err)
	}

	cleanJSON := removeJSONCComments(string(data))

	// Create a temporary map for the new JSON data
	jsonData := make(map[string]interface{})
	if err := json.Unmarshal([]byte(cleanJSON), &jsonData); err != nil {
		return fmt.Errorf("failed to parse JSON config: %w", err)
	}

	// Merge JSON data into existing config
	// JSON values will override existing values if they exist
	for k, v := range jsonData {
		c.data[k] = v
	}

	// Update cached values if they were overridden
	if port, err := c.GetInt("port"); err == nil {
		c.cachedValues.port = port
	}
	if port_, err := c.GetInt("app_port"); err == nil {
		c.cachedValues.port = port_
	}

	return nil
}

// GetString returns string value from either env or JSONC
func (c *ConfigMap) GetString(key string) string {
	// First check environment
	if val := os.Getenv(key); val != "" {
		return val
	}

	// Then check JSONC
	if val, ok := c.data[key]; ok {
		switch v := val.(type) {
		case string:
			return v
		default:
			return fmt.Sprintf("%v", v)
		}
	}

	return ""
}

// GetInt returns int value from either env or JSONC
func (c *ConfigMap) GetInt(key string) (int, error) {
	// First check environment
	if val := os.Getenv(key); val != "" {
		return strconv.Atoi(val)
	}

	// Then check JSONC
	if val, ok := c.data[key]; ok {
		switch v := val.(type) {
		case float64:
			return int(v), nil
		case int:
			return v, nil
		case string:
			return strconv.Atoi(v)
		}
	}

	return 0, fmt.Errorf("key not found or invalid type: %s", key)
}

// GetBool returns boolean value from either env or JSONC
func (c *ConfigMap) GetBool(key string) (bool, error) {
	// First check environment
	if val := os.Getenv(key); val != "" {
		return strconv.ParseBool(val)
	}

	// Then check JSONC
	if val, ok := c.data[key]; ok {
		switch v := val.(type) {
		case bool:
			return v, nil
		case string:
			return strconv.ParseBool(v)
		case float64:
			return v != 0, nil
		case int:
			return v != 0, nil
		}
	}

	return false, fmt.Errorf("key not found or invalid type: %s", key)
}

// removeJSONCComments removes JSONC style comments
func removeJSONCComments(input string) string {
	lines := strings.Split(input, "\n")
	var result []string

	for _, line := range lines {
		if idx := strings.Index(line, "//"); idx >= 0 {
			line = line[:idx]
		}
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// Use cached values
func (c *ConfigMap) Port() int {
	return c.cachedValues.port // Direct access, no lookup
}

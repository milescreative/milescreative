package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config provides a unified interface for both JSONC and environment variables
type Config struct {
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
func NewConfig(jsonPath string) (*Config, error) {
	c := &Config{}

	envPath, err := findEnvFile()
	if err != nil {
		return nil, fmt.Errorf("error finding .env file: %w", err)
	}

	fmt.Printf("Loading .env from: %s\n", envPath)
	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Load JSON config for constants
	if err := c.loadJSON(jsonPath); err != nil {
		return nil, err
	}

	// Cache frequently used values
	port, _ := c.GetInt("port")
	c.cachedValues.port = port

	return c, nil
}

func (c *Config) loadJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read JSON config file: %w", err)
	}

	cleanJSON := removeJSONCComments(string(data))

	if err := json.Unmarshal([]byte(cleanJSON), &c.data); err != nil {
		return fmt.Errorf("failed to parse JSON config: %w", err)
	}

	return nil
}

// GetString returns string value from either env or JSONC
func (c *Config) GetString(key string) string {
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
func (c *Config) GetInt(key string) (int, error) {
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
func (c *Config) GetBool(key string) (bool, error) {
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
func (c *Config) Port() int {
	return c.cachedValues.port // Direct access, no lookup
}

func main() {
	// Example usage
	configPath := filepath.Join("config", "test.jsonc")
	cfg, err := NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(os.Getenv("DB_HOST"))

	// Example of how to use it in your app:

	// Get app constants from JSONC
	appName := cfg.GetString("app_name")
	appVersion := cfg.GetString("app_version")

	// Get environment-specific values
	dbHost := cfg.GetString("DB_HOST") // from env
	dbPort := cfg.GetString("DB_PORT") // from env
	apiKey := cfg.GetString("API_KEY") // from env

	fmt.Println("Application Configuration:")
	fmt.Printf("Name: %s\n", appName)
	fmt.Printf("Version: %s\n", appVersion)
	fmt.Printf("Database Host: %s\n", dbHost)
	fmt.Printf("Database Port: %s\n", dbPort)
	fmt.Printf("API Key: %s\n", apiKey)
}

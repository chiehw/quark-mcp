package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	configDirName       = ".quark-mcp"
	legacyConfigDirName = ".quark-nd-disk"
	configDirPerm       = os.FileMode(0700)
	configFilePerm      = os.FileMode(0600)
)

type Config struct {
	Cookie string `json:"cookie"`
}

// lookupHome is replaced in tests.
var lookupHome = os.UserHomeDir

// DefaultConfigPath is the canonical JSON config path: ~/.quark-mcp/config.json
var DefaultConfigPath string

func init() {
	DefaultConfigPath = defaultConfigPath()
}

func homeDir() string {
	home, err := lookupHome()
	if err != nil || home == "" {
		home = os.Getenv("HOME")
	}
	return home
}

func defaultConfigPath() string {
	return filepath.Join(homeDir(), configDirName, "config.json")
}

func legacyConfigPath() string {
	return filepath.Join(homeDir(), legacyConfigDirName, "config.json")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isManagedConfigDir(dir string) bool {
	base := filepath.Base(dir)
	return base == configDirName || base == legacyConfigDirName
}

// ResolvePath returns path if set. Otherwise it prefers ~/.quark-mcp/config.json,
// then the legacy ~/.quark-nd-disk/config.json if that file still exists.
func ResolvePath(path string) string {
	if path != "" {
		return path
	}
	current := defaultConfigPath()
	if fileExists(current) {
		return current
	}
	legacy := legacyConfigPath()
	if fileExists(legacy) {
		return legacy
	}
	return current
}

func cookieRequiredError(path string) error {
	return fmt.Errorf("cookie is required: set QUARK_COOKIE or configure cookie in %s", path)
}

func tightenFilePermissions(path string) error {
	if err := os.Chmod(path, configFilePerm); err != nil && runtime.GOOS != "windows" {
		return fmt.Errorf("failed to restrict config file permissions: %w", err)
	}
	return nil
}

// Load reads the config. QUARK_COOKIE takes precedence over the JSON file.
func Load(path string) (*Config, error) {
	if cookie := strings.TrimSpace(os.Getenv("QUARK_COOKIE")); cookie != "" {
		return &Config{Cookie: cookie}, nil
	}

	path = ResolvePath(path)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, cookieRequiredError(path)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	cfg.Cookie = strings.TrimSpace(cfg.Cookie)
	if cfg.Cookie == "" {
		return nil, cookieRequiredError(path)
	}

	_ = tightenFilePermissions(path)
	return &cfg, nil
}

// LoadOrEmpty reads the config file, returns empty config if not exists
func LoadOrEmpty(path string) (*Config, error) {
	path = ResolvePath(path)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Init creates a new config file at the canonical path unless path is set
func Init(path string) error {
	if path == "" {
		path = defaultConfigPath()
	}

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("config file already exists at %s", path)
	}

	cfg := &Config{
		Cookie: "",
	}

	return Save(path, cfg)
}

// Save writes the config to file with user-only permissions
func Save(path string, cfg *Config) error {
	if path == "" {
		path = defaultConfigPath()
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, configDirPerm); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	if isManagedConfigDir(dir) {
		if err := os.Chmod(dir, configDirPerm); err != nil && runtime.GOOS != "windows" {
			return fmt.Errorf("failed to restrict config directory permissions: %w", err)
		}
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, configFilePerm); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	if err := tightenFilePermissions(path); err != nil {
		return err
	}

	return nil
}

func maskCookie(cookie string) string {
	if len(cookie) > 20 {
		return cookie[:10] + "..." + cookie[len(cookie)-10:]
	}
	if len(cookie) > 5 {
		return cookie[:5] + "..."
	}
	return "(set)"
}

// Set sets a config value by key
func Set(path, key, value string) error {
	path = ResolvePath(path)
	cfg, err := LoadOrEmpty(path)
	if err != nil {
		return err
	}

	switch key {
	case "cookie":
		cfg.Cookie = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	return Save(path, cfg)
}

// Get gets a config value by key
func Get(path, key string) (string, error) {
	path = ResolvePath(path)
	cfg, err := LoadOrEmpty(path)
	if err != nil {
		return "", err
	}

	switch key {
	case "cookie":
		if cfg.Cookie == "" {
			return "", fmt.Errorf("cookie is not set")
		}
		return maskCookie(cfg.Cookie), nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}

// Show shows all config values
func Show(path string) (map[string]string, error) {
	path = ResolvePath(path)
	cfg, err := LoadOrEmpty(path)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	if cfg.Cookie != "" {
		result["cookie"] = maskCookie(cfg.Cookie)
	} else {
		result["cookie"] = "(not set)"
	}

	return result, nil
}

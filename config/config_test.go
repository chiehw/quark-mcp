package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestLoadPrefersQUARK_COOKIEWithoutConfigFile(t *testing.T) {
	t.Setenv("QUARK_COOKIE", "env-cookie-value")

	cfg, err := Load(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Cookie != "env-cookie-value" {
		t.Fatalf("Cookie = %q, want env value", cfg.Cookie)
	}
}

func TestLoadPrefersQUARK_COOKIEOverJSON(t *testing.T) {
	t.Setenv("QUARK_COOKIE", "env-cookie-value")
	path := writeConfig(t, `{"cookie":"file-cookie-value"}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Cookie != "env-cookie-value" {
		t.Fatalf("Cookie = %q, want env to win", cfg.Cookie)
	}
}

func TestLoadFallsBackToJSON(t *testing.T) {
	t.Setenv("QUARK_COOKIE", "")
	path := writeConfig(t, `{"cookie":"file-cookie-value"}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Cookie != "file-cookie-value" {
		t.Fatalf("Cookie = %q, want file value", cfg.Cookie)
	}
}

func TestLoadRequiresCookieFromEitherSource(t *testing.T) {
	t.Setenv("QUARK_COOKIE", "")
	path := filepath.Join(t.TempDir(), "config.json")

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want missing-cookie error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "QUARK_COOKIE") || !strings.Contains(msg, path) {
		t.Fatalf("error %q should mention QUARK_COOKIE and the JSON path", msg)
	}
}

func TestLoadEmptyJSONCookieRequiresEnvOrFile(t *testing.T) {
	t.Setenv("QUARK_COOKIE", "")
	path := writeConfig(t, `{"cookie":"   "}`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want missing-cookie error")
	}
	if !strings.Contains(err.Error(), "QUARK_COOKIE") {
		t.Fatalf("error %q should mention QUARK_COOKIE", err)
	}
}

func TestSaveUsesUserOnlyPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix permission bits are not enforced on windows")
	}

	path := filepath.Join(t.TempDir(), ".quark-mcp", "config.json")
	if err := Save(path, &Config{Cookie: "secret"}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat(file) error = %v", err)
	}
	if perm := fileInfo.Mode().Perm(); perm != 0600 {
		t.Fatalf("config file perm = %o, want 0600", perm)
	}

	dirInfo, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatalf("Stat(dir) error = %v", err)
	}
	if perm := dirInfo.Mode().Perm(); perm != 0700 {
		t.Fatalf("config dir perm = %o, want 0700", perm)
	}
}

func TestGetMasksCookie(t *testing.T) {
	path := writeConfig(t, `{"cookie":"abcdefghijklmnopqrstuvwxyz"}`)

	got, err := Get(path, "cookie")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if strings.Contains(got, "abcdefghijklmnopqrstuvwxyz") {
		t.Fatalf("Get() leaked cookie: %q", got)
	}
	if !strings.Contains(got, "...") {
		t.Fatalf("Get() = %q, want masked value", got)
	}
}

func TestResolvePathPrefersNewConfigOverLegacy(t *testing.T) {
	t.Setenv("QUARK_COOKIE", "")
	home := t.TempDir()
	lookupHome = func() (string, error) { return home, nil }
	t.Cleanup(func() { lookupHome = os.UserHomeDir })

	newPath := filepath.Join(home, ".quark-mcp", "config.json")
	legacyPath := filepath.Join(home, ".quark-nd-disk", "config.json")
	if err := os.MkdirAll(filepath.Dir(newPath), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newPath, []byte(`{"cookie":"new-cookie"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(legacyPath, []byte(`{"cookie":"legacy-cookie"}`), 0600); err != nil {
		t.Fatal(err)
	}

	if got := ResolvePath(""); got != newPath {
		t.Fatalf("ResolvePath() = %q, want %q", got, newPath)
	}
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Cookie != "new-cookie" {
		t.Fatalf("Cookie = %q, want new-cookie", cfg.Cookie)
	}
}

func TestLoadFallsBackToLegacyConfigPath(t *testing.T) {
	t.Setenv("QUARK_COOKIE", "")
	home := t.TempDir()
	lookupHome = func() (string, error) { return home, nil }
	t.Cleanup(func() { lookupHome = os.UserHomeDir })

	legacyPath := filepath.Join(home, ".quark-nd-disk", "config.json")
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(legacyPath, []byte(`{"cookie":"legacy-cookie"}`), 0600); err != nil {
		t.Fatal(err)
	}

	if got := ResolvePath(""); got != legacyPath {
		t.Fatalf("ResolvePath() = %q, want legacy %q", got, legacyPath)
	}
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Cookie != "legacy-cookie" {
		t.Fatalf("Cookie = %q, want legacy-cookie", cfg.Cookie)
	}
}

func writeConfig(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

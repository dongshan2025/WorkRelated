package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInitAndGetters(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")
	content := `app:
  name: testapp
  port: 9000
debug: false
`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	// 初始化
	v = nil // reset global viper instance
	if err := Init(cfgPath, ""); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if got := GetString("app.name"); got != "testapp" {
		t.Fatalf("app.name = %q, want %q", got, "testapp")
	}
	if got := GetInt("app.port"); got != 9000 {
		t.Fatalf("app.port = %d, want %d", got, 9000)
	}
	if got := GetBool("debug"); got != false {
		t.Fatalf("debug = %v, want %v", got, false)
	}

	// Unmarshal
	var out struct {
		App struct {
			Name string `mapstructure:"name"`
			Port int    `mapstructure:"port"`
		} `mapstructure:"app"`
		Debug bool `mapstructure:"debug"`
	}
	if err := Unmarshal(&out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.App.Name != "testapp" || out.App.Port != 9000 || out.Debug != false {
		t.Fatalf("unmarshal wrong: %+v", out)
	}
}

func TestEnvFallback(t *testing.T) {
	// ensure no config file
	v = nil
	os.Unsetenv("TEST_APP_NAME")
	os.Setenv("TEST_APP_NAME", "envapp")
	defer os.Unsetenv("TEST_APP_NAME")

	if err := Init("", "TEST"); err != nil {
		t.Fatalf("Init with env failed: %v", err)
	}

	if got := GetString("app.name"); got != "envapp" {
		t.Fatalf("app.name from env = %q, want %q", got, "envapp")
	}
}

func TestWatchConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")
	content := `app:
  name: watchapp
  port: 7000
debug: false
`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	v = nil
	if err := Init(cfgPath, ""); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	ch := make(chan struct{}, 1)
	WatchConfig(func() {
		select {
		case ch <- struct{}{}:
		default:
		}
	})

	// modify file
	newContent := `app:
  name: watchapp2
  port: 7001
debug: true
`
	if err := os.WriteFile(cfgPath, []byte(newContent), 0644); err != nil {
		t.Fatalf("write new config: %v", err)
	}

	select {
	case <-ch:
		// ok
	case <-time.After(3 * time.Second):
		t.Fatalf("config change not detected within timeout")
	}
}

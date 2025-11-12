package config

import (
	"fmt"
	"log"
	"time"

	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var v *viper.Viper

// Init 初始化配置。configFile 可以是完整路径（带扩展名）或目录/文件名。
// 如果 configFile 为空，viper 会尝试在当前目录查找 config.(yaml|json|toml) 等。
func Init(configFile string, envPrefix string) error {
	if v == nil {
		v = viper.New()
	}

	// 1) 支持指定完整文件路径
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.AddConfigPath(".")
		v.SetConfigName("config")
	}

	// 2) 环境变量覆盖
	if envPrefix != "" {
		v.SetEnvPrefix(envPrefix)
	}
	// 将 '.' 替换为 '_' 以支持从环境变量读取嵌套 key（如 APP_NAME -> app.name）
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		// 如果找不到文件也不是致命错误，可以允许仅使用 env
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Infof("no config file found, continuing with environment variables only")
			return nil
		}
		logger.Errorf("read config error: %v", err)
		return fmt.Errorf("read config: %w", err)
	}

	logger.Infof("using config file: %s", v.ConfigFileUsed())
	return nil
}

// WatchConfig 当配置文件改变时调用回调
func WatchConfig(onChange func()) {
	if v == nil {
		return
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		logger.Infof("config changed: %s", e.Name)
		if onChange != nil {
			onChange()
		}
	})
}

// Logger is a minimal logging interface used by this package. You can supply
// your own implementation via SetLogger.
type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Debugf(format string, args ...any)
}

type stdLogger struct{}

func (stdLogger) Infof(format string, args ...any)  { log.Printf("[INFO] "+format, args...) }
func (stdLogger) Errorf(format string, args ...any) { log.Printf("[ERROR] "+format, args...) }
func (stdLogger) Debugf(format string, args ...any) { log.Printf("[DEBUG] "+format, args...) }

var logger Logger = stdLogger{}

// SetLogger sets a custom logger for the config package. Passing nil will be ignored.
func SetLogger(l Logger) {
	if l == nil {
		return
	}
	logger = l
}

// Get helpers
func GetString(key string) string { return v.GetString(key) }
func GetInt(key string) int       { return v.GetInt(key) }
func GetBool(key string) bool     { return v.GetBool(key) }
func GetDuration(key string) time.Duration {
	return v.GetDuration(key)
}

// Unmarshal 将配置映射到结构体
func Unmarshal(rawVal any) error {
	if v == nil {
		return fmt.Errorf("viper not initialized")
	}
	return v.Unmarshal(rawVal)
}

// UnmarshalKey 将某个 key 映射到结构体
func UnmarshalKey(key string, rawVal any) error {
	if v == nil {
		return fmt.Errorf("viper not initialized")
	}
	return v.UnmarshalKey(key, rawVal)
}

// GetViper 返回底层 viper 实例以便高级用法
func GetViper() *viper.Viper { return v }

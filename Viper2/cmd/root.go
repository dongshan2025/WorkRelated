package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"viper2/config"

	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	envPrefix   string
	verbose     bool
	logEncoding string
	// version can be set at build time with -ldflags "-X viper2/cmd.version=..."
	version = "v0.1.0"
)

const (
	defaultConfigFile = "./config.yaml"
	defaultEnvPrefix  = "MYAPP"
)

// LoadResult describes how configuration was loaded.
// LoadedFromFile == true means the effective configuration came from a file.
// false means configuration was initialized from environment variables only.
type LoadResult struct {
	LoadedFromFile bool
}

// rootCmd is the base command for the app
var rootCmd = &cobra.Command{
	Use:   "viper2",
	Short: "Viper2 demo app (config + zap)",
	Long:  "Demo application that shows viper-based config and zap integration.",
	Example: `  # run with default config (./config.yaml) or fall back to env vars
  viper2 serve

  # run using a specific config file
  viper2 serve --config ./custom.yaml

  # run with environment variables only (no file)
  viper2 serve --config "" --env-prefix MYAPP2

  # print version
  viper2 version
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// show help by default
		return cmd.Help()
	},
	// PersistentPreRunE runs before any subcommand's Run/RunE. Use it to
	// validate global flags (like --log-encoding) early and fail fast.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if logEncoding != "" {
			le := logEncoding
			if le != "json" && le != "console" {
				return fmt.Errorf("invalid --log-encoding %q: allowed values are 'json' or 'console'", le)
			}
		}
		return nil
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// persistent flags are available to all subcommands
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", defaultConfigFile, "path to config file (if empty will use only environment variables)")
	rootCmd.PersistentFlags().StringVar(&envPrefix, "env-prefix", defaultEnvPrefix, "environment variable prefix to use when loading config from env")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output (sets default log.level=debug when not configured)")
	rootCmd.PersistentFlags().StringVar(&logEncoding, "log-encoding", "", "log encoding to use (json|console) - overrides config file and environment")

	// add subcommands
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(validateCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the demo server (reads config and initializes logger)",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration (from file if present, otherwise from env)
		res, err := loadConfig(cfgFile, envPrefix)
		if err != nil {
			return err
		}

		// If verbose flag set and log.level not explicitly configured (file/env),
		// set a default level to debug. This uses viper's SetDefault so that
		// a value from file or env still overrides it.
		if verbose {
			if v := config.GetViper(); v != nil && !v.IsSet("log.level") {
				v.SetDefault("log.level", "debug")
			}
		}

		// Setup zap from config (if any) and use it to log the startup banner.
		l, err := config.SetupZapFromConfig()
		if err != nil {
			// log the error and continue with standard logger
			log.Printf("warning: failed to setup zap from config: %v", err)
		}
		// gather some runtime metadata for structured startup logging
		v := config.GetViper()
		cfgFileUsed := ""
		if v != nil {
			cfgFileUsed = v.ConfigFileUsed()
		}

		if l != nil {
			// ensure flush
			defer l.Sync()
			// prefer sugared logger and log banner as structured message including server id
			sugar := l.Sugar()
			sugar.Infow("startup",
				"banner", banner(),
				"server.id", config.GetString("server.id"),
				"app.name", config.GetString("app.name"),
				"version", version,
				"config.file", cfgFileUsed,
				"loaded_from_file", res.LoadedFromFile,
			)
		} else {
			// fallback to standard logger: print banner and server id
			log.Printf("%s\nserver.id=%s app.name=%s version=%s config.file=%s loaded_from_file=%v",
				banner(), config.GetString("server.id"), config.GetString("app.name"), version, cfgFileUsed, res.LoadedFromFile)
		}

		// Watch config changes
		config.WatchConfig(func() {
			fmt.Println("config changed")
		})

		// Application demo logic
		name := config.GetString("app.name")
		port := config.GetInt("app.port")
		debug := config.GetBool("debug")

		fmt.Printf("app.name=%s port=%d debug=%v\n", name, port, debug)

		var cfg struct {
			App struct {
				Name string `mapstructure:"name"`
				Port int    `mapstructure:"port"`
			} `mapstructure:"app"`
			Debug bool `mapstructure:"debug"`
		}
		_ = config.Unmarshal(&cfg)
		fmt.Printf("unmarshal cfg: %+v\n", cfg)

		// keep running a short while to allow testing of hot-reload
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		<-ctx.Done()

		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate and print loaded configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration and output a concise summary for validation
		// Load configuration and output a concise summary for validation
		_, err := loadConfig(cfgFile, envPrefix)
		if err != nil {
			return err
		}
		printConfigSummary()
		return nil
	},
}

// loadConfig centralizes config initialization logic used by commands.
// It prefers the provided config file if it exists; otherwise it initializes
// viper to read from environment variables using the provided prefix.
func loadConfig(cfg, env string) (LoadResult, error) {
	var res LoadResult
	if cfg != "" {
		if _, err := os.Stat(cfg); err == nil {
			if err := config.Init(cfg, env); err != nil {
				return res, fmt.Errorf("init config error: %w", err)
			}
			// If a persistent flag was provided, it should override the config file.
			// Bind the CLI flag to viper so command-line value wins over file/env.
			if f := rootCmd.PersistentFlags().Lookup("log-encoding"); f != nil {
				if vv := config.GetViper(); vv != nil {
					_ = vv.BindPFlag("log.encoding", f)

					// Also bind SERVER_ID env var to viper key "server.id" so that
					// environment injection (from Docker) is visible via viper.
					_ = vv.BindEnv("server.id", "SERVER_ID")
				}
			}
			res.LoadedFromFile = true
			return res, nil
		}
		// file specified but not found -> fall back to env
		fmt.Printf("config file not found (%s), will read configuration from environment variables with prefix %s_\n", cfg, env)
		if err := config.Init("", env); err != nil {
			return res, fmt.Errorf("init config error: %w", err)
		}
		if f := rootCmd.PersistentFlags().Lookup("log-encoding"); f != nil {
			if vv := config.GetViper(); vv != nil {
				_ = vv.BindPFlag("log.encoding", f)

				// bind SERVER_ID env var to viper
				_ = vv.BindEnv("server.id", "SERVER_ID")
			}
		}
		res.LoadedFromFile = false
		return res, nil
	}
	// explicit empty cfg -> use environment only
	if err := config.Init("", env); err != nil {
		return res, fmt.Errorf("init config error: %w", err)
	}
	if f := rootCmd.PersistentFlags().Lookup("log-encoding"); f != nil {
		if vv := config.GetViper(); vv != nil {
			_ = vv.BindPFlag("log.encoding", f)

			// bind SERVER_ID env var to viper
			_ = vv.BindEnv("server.id", "SERVER_ID")
		}
	}
	res.LoadedFromFile = false
	return res, nil
}

// printConfigSummary prints a small set of keys useful for quick validation.
func printConfigSummary() {
	v := config.GetViper()
	if v != nil {
		if f := v.ConfigFileUsed(); f != "" {
			fmt.Printf("config file used: %s\n", f)
		} else {
			fmt.Println("config file: <none> (using environment variables)")
		}
	}
	fmt.Printf("app.name=%s\n", config.GetString("app.name"))
	fmt.Printf("server.id=%s\n", config.GetString("server.id"))
	fmt.Printf("app.port=%d\n", config.GetInt("app.port"))
	fmt.Printf("debug=%v\n", config.GetBool("debug"))
	fmt.Printf("log.level=%s\n", config.GetString("log.level"))
	fmt.Printf("log.encoding=%s\n", config.GetString("log.encoding"))
	if v := config.GetViper(); v != nil {
		fmt.Printf("log.output_paths=%v\n", v.GetStringSlice("log.output_paths"))
	}
	fmt.Printf("log.encoder.time_format=%s\n", config.GetString("log.encoder.time_format"))
	fmt.Printf("log.sampling.enabled=%v\n", config.GetBool("log.sampling.enabled"))
}

// banner returns the Bee Message ASCII art used at startup.
func banner() string {
	// concise single-line banner used for structured startup logs
	return "Bee Message: ready"
}

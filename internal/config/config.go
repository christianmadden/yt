package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const defaultConfig = `# yt configuration
# https://github.com/yourusername/yt

# Target video height in pixels (360/480/720/1080/1440/2160)
preferred_quality: 1080

# Warn (but allow) if best available is below this height
min_quality: 720

# Output container format: mp4, mkv, webm
default_format: mp4

# Output directory (empty = current working directory)
default_dir: ""
`

type Config struct {
	PreferredQuality int    `mapstructure:"preferred_quality"`
	MinQuality       int    `mapstructure:"min_quality"`
	DefaultFormat    string `mapstructure:"default_format"`
	DefaultDir       string `mapstructure:"default_dir"`
}

func Load(cfgFile string) (*Config, error) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("could not find home directory: %w", err)
		}
		cfgDir := filepath.Join(home, ".config", "yt")
		cfgPath := filepath.Join(cfgDir, "config.yaml")

		// Create config file on first run
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			if err := os.MkdirAll(cfgDir, 0755); err != nil {
				return nil, fmt.Errorf("could not create config directory: %w", err)
			}
			if err := os.WriteFile(cfgPath, []byte(defaultConfig), 0644); err != nil {
				return nil, fmt.Errorf("could not write default config: %w", err)
			}
			fmt.Printf("Created default config at %s\n\n", cfgPath)
		}

		viper.SetConfigFile(cfgPath)
	}

	viper.SetDefault("preferred_quality", 1080)
	viper.SetDefault("min_quality", 720)
	viper.SetDefault("default_format", "mp4")
	viper.SetDefault("default_dir", "")

	if err := viper.ReadInConfig(); err != nil {
		// If config file doesn't exist, just use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return &cfg, nil
}

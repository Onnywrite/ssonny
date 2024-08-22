package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	//nolint: gochecknoglobals
	ConfigPathFlag = "config-path"
	//nolint: gochecknoglobals
	ConfigPathEnv = "CONFIG_PATH"
)

type Config struct {
	PostgresConn string `yaml:"postgresConn"`

	Https TransportConfig     `yaml:"https"`
	Grpc  GrpcTransportConfig `yaml:"grpc"`

	Tokens TokensConfig `yaml:"tokens"`
}

type TransportConfig struct {
	Port   uint16 `yaml:"port"`
	UseTLS bool   `yaml:"useTls"`
	Cert   string `yaml:"cert"`
	Key    string `yaml:"key"`
}
type GrpcTransportConfig struct {
	Port    uint16        `yaml:"port"`
	UseTLS  bool          `yaml:"useTls"`
	Cert    string        `yaml:"cert"`
	Key     string        `yaml:"key"`
	Timeout time.Duration `yaml:"timeout"`
}

type TokensConfig struct {
	Issuer     string        `yaml:"issuer"`
	SecretPath string        `yaml:"secretPath"`
	PublicPath string        `yaml:"publicPath"`
	AccessTTL  time.Duration `yaml:"accessTtl"`
	IdTTL      time.Duration `yaml:"idTtl"`
	RefreshTTL time.Duration `yaml:"refreshTtl"`
}

func MustLoad(defaultPath string) *Config {
	conf, err := Load(defaultPath)
	if err != nil {
		panic(err)
	}

	return conf
}

func Load(defaultPath string) (*Config, error) {
	var configPath string

	flag.StringVar(&configPath, ConfigPathFlag, "", "config file path")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv(ConfigPathEnv)
	}

	if configPath == "" {
		configPath = defaultPath
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: path '%s'", err, configPath)
	}

	return LoadPath(configPath)
}

func LoadPath(path string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("config could not be loaded: %w", err)
	}

	return &cfg, nil
}

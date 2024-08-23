package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// OSNOVA (BASE):
// 1. flag --config-path specified:
// 		we're on the host machine
// 2. flag doesnt specified, but sso.yaml found in /etc/sso/
// 		we're inside of a container
//
//

type Config struct {
	Containerless Containerless `yaml:"containerless"`
	Http          Http          `yaml:"http"`
	Grpc          Grpc          `yaml:"grpc"`
	Tokens        Tokens        `yaml:"tokens"`
}

type Containerless struct {
	PostgresConn string `yaml:"postgresConn"`
	TlsCertPath  string `yaml:"tlsCertPath"`
	TlsKeyPath   string `yaml:"tlsKeyPath"`
	SecretString string `yaml:"secretString"`
}

type Http struct {
	Port   int  `yaml:"port"`
	UseTls bool `yaml:"useTls"`
}

type Grpc struct {
	Port    int           `yaml:"port"`
	UseTls  bool          `yaml:"useTls"`
	Timeout time.Duration `yaml:"timeout"`
}

type Tokens struct {
	Issuer               string        `yaml:"issuer"`
	AccessTtl            time.Duration `yaml:"accessTtl"`
	IdTtl                time.Duration `yaml:"idTtl"`
	RefreshTtl           time.Duration `yaml:"refreshTtl"`
	EmailVerificationTtl time.Duration `yaml:"emailVerificationTtl"`
}

func MustLoad() *Config {
	filePath, isContainer := findFile()

	file, errr := os.ReadFile(filePath)
	if errr != nil {
		log.Fatal(errr)
	}

	var config Config

	err := yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}

	if isContainer {
		config.Containerless.PostgresConn = os.Getenv("POSTGRES_CONN")
		config.Containerless.TlsCertPath = "/etc/sso/cert.pem"
		config.Containerless.TlsKeyPath = "/etc/sso/key.pem"
		config.Containerless.SecretString = os.Getenv("SECRET_STRING")

		return &config
	}

	if !filepath.IsAbs(config.Containerless.TlsCertPath) {
		config.Containerless.TlsCertPath = filepath.Join(filePath, config.Containerless.TlsCertPath)
	}

	if !filepath.IsAbs(config.Containerless.TlsKeyPath) {
		config.Containerless.TlsKeyPath = filepath.Join(filePath, config.Containerless.TlsKeyPath)
	}

	return &config
}

// true when in a container
// false otherwise
func findFile() (string, bool) {
	var configPath string

	flag.StringVar(&configPath, "config-path", "./sso.yaml", "path to config file")
	flag.Parse()

	if _, err := os.Stat(configPath); err == nil {
		return configPath, false
	}

	if _, err := os.Stat("/etc/sso/sso.yaml"); err != nil {
		log.Fatalf("no config found in neither /etc/sso/sso.yaml nor %s from flag", configPath)
	}

	return "/etc/sso/sso.yaml", true
}

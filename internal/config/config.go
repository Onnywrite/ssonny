package config

import (
	"flag"
	"fmt"
	"os"
)

var (
	ConfigPathFlag = "config-path"
	ConfigPathEnv  = "CONFIG_PATH"

	TlsKeyPathEnv   = "TLS_KEY_PATH"
	TlsCertPathEnv  = "TLS_CERT_PATH"
	PostgresConnEnv = "POSTGRES_CONN"
)

const (
	tlsCertDefaultPath = "/secrets/cert"
	tlsKeyDefaultPath  = "/secrets/key"
	secretTlsCertPath  = "secrets.tlsCertPath"
	secretTlsKeyPath   = "secrets.tlsKeyPath"
)

const (
	SecretPostgresConn = "secrets.postgresConn"
	SecretTlsCert      = "secrets.tlsCert"
	SecretTlsKey       = "secrets.tlsKey"
)

const (
	HttpPort   = "http.port"
	HttpUseTLS = "http.useTls"

	GrpcPort    = "grpc.port"
	GrpcUseTLS  = "grpc.useTls"
	GrpcTimeout = "grpc.timeout"

	TokensIssuer     = "tokens.issuer"
	TokensAccessTtl  = "tokens.accessTtl"
	TokensIdTtl      = "tokens.idTtl"
	TokensRefreshTtl = "tokens.refreshTtl"
)

type Configer interface {
	Get(key string) any
}

func MustLoad(path string) Configer {
	c, err := Load(path)
	if err != nil {
		panic(err)
	}
	return c
}

func Load(path string) (Configer, error) {
	var flagPath string
	flag.StringVar(&flagPath, ConfigPathFlag, "./configs", "path to a config file")
	flag.Parse()

	fmt.Println("config", flagPath, os.Getenv(ConfigPathEnv), path)
	return newViper("config", "yaml", flagPath, os.Getenv(ConfigPathEnv), path)
}

func MustGet[T any](c Configer, key string) T {
	t, err := Get[T](c, key)
	if err != nil {
		panic(err)
	}
	return t
}

func Get[T any](c Configer, key string) (T, error) {
	if t, ok := c.Get(key).(T); ok {
		return t, nil
	}
	empty := *new(T)
	return empty, fmt.Errorf("expected %T type for %s", empty, key)
}

package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	secretTlsCertPath  = "secrets.tls_cert_path"
	secretTlsKeyPath   = "secrets.tls_key_path"
)

const (
	SecretPostgresConn = "secrets.postgres_conn"
	SecretTlsCert      = "secrets.tls_cert"
	SecretTlsKey       = "secrets.tls_key"
)

const (
	HttpPort   = "http.port"
	HttpUseTLS = "http.use_tls"

	GrpcPort    = "grpc.port"
	GrpcUseTLS  = "grpc.use_tls"
	GrpcTimeout = "grpc.timeout"

	TokensIssuer     = "tokens.issuer"
	TokensAccessTtl  = "tokens.access_ttl"
	TokensIdTtl      = "tokens.id_ttl"
	TokensRefreshTtl = "tokens.refresh_ttl"
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
	flag.StringVar(&flagPath, ConfigPathFlag, "./config", "path to a config file")
	flag.Parse()

	filename := strings.TrimSuffix(filepath.Base(path), ".yaml")
	fmt.Println(filename, filepath.Dir(flagPath), filepath.Dir(os.Getenv(ConfigPathEnv)), filepath.Dir(path))
	return newViper(filename, "yaml", filepath.Dir(flagPath), filepath.Dir(os.Getenv(ConfigPathEnv)), filepath.Dir(path))
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

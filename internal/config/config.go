package config

import "time"

type Config struct {
	Tls     Tls     `yaml:"tls"`
	Secrets Secrets `yaml:"secrets"`
	Http    Http    `yaml:"http"`
	Grpc    Grpc    `yaml:"grpc"`
	Tokens  Tokens  `yaml:"tokens"`
}

type Tls struct {
	CertPath string `env:"TLS_CERT_PATH" env-default:"/etc/sso/tls.crt" yaml:"certPath"`
	KeyPath  string `env:"TLS_KEY_PATH"  env-default:"/etc/sso/tls.key" yaml:"keyPath"`
}

type Secrets struct {
	Postgres     Postgres `yaml:"postgres"`
	Redis        Redis    `yaml:"redis"`
	SecretString string   `env:"SECRET_STRING" yaml:"secretString"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" yaml:"host"`
	Port     int    `env:"POSTGRES_PORT" env-default:"5432" yaml:"port"`
	User     string `env:"POSTGRES_USER" yaml:"user"`
	Password string `env:"POSTGRES_PASS" yaml:"password"`
	Database string `env:"POSTGRES_DB"   yaml:"database"`
	SslMode  string `env:"POSTGRES_SSL"  env-default:"disable" yaml:"sslmode"`
}

type Redis struct {
	Host     string `env:"REDIS_HOST" yaml:"host"`
	Port     int    `env:"REDIS_PORT" env-default:"6379" yaml:"port"`
	User     string `env:"REDIS_USER" yaml:"user"`
	Password string `env:"REDIS_PASS" yaml:"password"`
	Db       int    `env:"REDIS_DB"   yaml:"db"`
}

type Http struct {
	Port   int  `env:"HTTP_PORT"    env-default:"8080"  yaml:"port"`
	UseTls bool `env:"HTTP_USE_TLS" env-default:"false" yaml:"useTls"`
}

type Grpc struct {
	Port    int           `env:"GRPC_PORT"    env-default:"5050"  yaml:"port"`
	UseTls  bool          `env:"GRPC_USE_TLS" env-default:"false" yaml:"useTls"`
	Timeout time.Duration `env:"GRPC_TIMEOUT" env-default:"5s"    yaml:"timeout"`
}

type Tokens struct {
	Issuer               string        `env:"TOKENS_ISSUER"                 yaml:"issuer"`
	AccessTtl            time.Duration `env:"TOKENS_ACCESS_TTL"             yaml:"accessTtl"`
	RefreshTtl           time.Duration `env:"TOKENS_REFRESH_TTL"            yaml:"refreshTtl"`
	IdTtl                time.Duration `env:"TOKENS_ID_TTL"                 yaml:"idTtl"`
	EmailVerificationTtl time.Duration `env:"TOKENS_EMAIL_VERIFICATION_TTL" yaml:"emailVerificationTtl"` //nolint: lll
}

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/spf13/viper"
)

type viperConfig struct {
	data map[string]any
	mu   sync.RWMutex
	v    *viper.Viper
}

func newViper(filename, extension string, paths ...string) (c *viperConfig, err error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.SetConfigType(extension)
	for _, path := range paths {
		v.AddConfigPath(path)
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	c = &viperConfig{
		data: make(map[string]any),
		v:    v,
	}
	err = c.checkData()

	return
}

var postgresConnRegex = regexp.MustCompile(`^postgres:\/\/.+:.+@.+:\d{1,5}\/.+(\?.*)?$`)

func (c *viperConfig) checkData() error {
	if postgresConnRegex.MatchString(c.v.GetString(SecretPostgresConn)) {
		c.data[SecretPostgresConn] = c.v.Get(SecretPostgresConn)
	} else if postgresConnRegex.MatchString(os.Getenv(PostgresConnEnv)) {
		c.data[SecretPostgresConn] = os.Getenv(PostgresConnEnv)
	} else {
		return fmt.Errorf("invalid postgres connection string")
	}

	configDir := filepath.Dir(c.v.ConfigFileUsed())
	fmt.Println("cfgDir:", configDir)

	if err := c.findTls(SecretTlsKey, os.Getenv(TlsKeyPathEnv), makeAbsoluteIfNeeded(c.v.GetString(secretTlsKeyPath), configDir), tlsKeyDefaultPath); err != nil {
		return err
	}

	if err := c.findTls(SecretTlsCert, os.Getenv(TlsCertPathEnv), makeAbsoluteIfNeeded(c.v.GetString(secretTlsCertPath), configDir), tlsCertDefaultPath); err != nil {
		return err
	}

	return nil
}

func makeAbsoluteIfNeeded(targetPath, pathBeforeDot string) string {
	if filepath.IsAbs(targetPath) {
		return targetPath
	}
	return filepath.Join(pathBeforeDot, targetPath)
}

// make path related to config as well (maybe)
func (c *viperConfig) findTls(secretTls string, paths ...string) error {
	for _, path := range paths {
		tls, err := os.ReadFile(path)
		fmt.Println(path, err)
		switch {
		case os.IsNotExist(err):
			continue
		case err != nil:
			return err
		default:
			c.data[secretTls] = tls
			return nil
		}
	}
	return fmt.Errorf(secretTls + " not found")
}

func (c *viperConfig) Get(key string) any {
	c.mu.RLock()
	value, ok := c.data[key]
	c.mu.RUnlock()
	if ok {
		return value
	}

	value = c.v.Get(key)
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
	return value
}

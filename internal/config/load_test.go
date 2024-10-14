package config

import (
	"flag"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_relatePath(t *testing.T) {
	tt := []struct {
		name   string
		path   string
		target string
		want   string
	}{
		{
			name:   "from example",
			path:   "./certs/tls.key",
			target: "../configs/conf.yaml",
			want:   "../configs/certs/tls.key",
		},
		{
			name:   "absolute path",
			path:   "/etc/app/tls.crt",
			target: "../configs/conf.yaml",
			want:   "/etc/app/tls.crt",
		},
		{
			name:   "empty path",
			path:   "",
			target: "../configs/conf.yaml",
			want:   "",
		},
		{
			name:   "from parent",
			path:   "certs/tls.key",
			target: "../configs/conf.yaml",
			want:   "../configs/certs/tls.key",
		},
		{
			name:   "absolute target",
			path:   "../../ssl/tls.key",
			target: "/etc/app/configs/conf.yaml",
			want:   "/etc/ssl/tls.key",
		},

		{
			name:   "same dir",
			path:   "tls.key",
			target: "./configs/conf.yaml",
			want:   "configs/tls.key",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := relatePath(tc.path, tc.target)

			assert.Equal(t, tc.want, got)
		})
	}
}

// I haven't found the way to test flag.Parse
func DontTest_getPath(t *testing.T) {
	tt := []struct {
		name     string
		env      string
		flag     string
		paths    []string
		existing []string
		want     string
		wantErr  bool
	}{
		{
			name:     "from flag",
			env:      "env/conf.yaml",
			flag:     "flag/conf.yaml",
			paths:    []string{"path1/conf.yaml", "/etc/app/conf.yaml"},
			existing: []string{"env/conf.yaml", "flag/conf.yaml", "path1/conf.yaml", "/etc/app/conf.yaml"},
			want:     "flag/conf.yaml",
			wantErr:  false,
		},
		{
			name:     "from env",
			env:      "env/conf.yaml",
			flag:     "",
			paths:    []string{"path1/conf.yaml", "/etc/app/conf.yaml"},
			existing: []string{"env/conf.yaml", "flag/conf.yaml", "path1/conf.yaml", "/etc/app/conf.yaml"},
			want:     "env/conf.yaml",
			wantErr:  false,
		},
		{
			name:     "first default",
			env:      "",
			flag:     "",
			paths:    []string{"path1/conf.yaml", "/etc/app/conf.yaml"},
			existing: []string{"env/conf.yaml", "flag/conf.yaml", "path1/conf.yaml", "/etc/app/conf.yaml"},
			want:     "path1/conf.yaml",
			wantErr:  false,
		},
		{
			name:     "second default",
			env:      "",
			flag:     "",
			paths:    []string{"path1/conf.yaml", "/etc/app/conf.yaml"},
			existing: []string{"env/conf.yaml", "flag/conf.yaml", "/etc/app/conf.yaml"},
			want:     "conf.yaml",
			wantErr:  false,
		},
		{
			name:     "not found",
			env:      "",
			flag:     "",
			paths:    []string{"path1/conf.yaml", "/etc/app/conf.yaml"},
			existing: []string{"env/conf.yaml", "flag/conf.yaml"},
			want:     "",
			wantErr:  true,
		},
	}

	mockExist := func(existing []string) existFunc {
		return func(name string) bool {
			return slices.Contains(existing, name)
		}
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := os.Setenv("CONFIG_PATH", tc.env)
			assert.NoError(t, err)

			if tc.flag != "" {
				err = flag.CommandLine.Set("--config-path", tc.flag)
				assert.NoError(t, err)
			}

			got, err := getPath(mockExist(tc.existing), tc.paths...)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, "", got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

package database

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/sbuzas-jwl/go-pkgs/todo/pkg/secrets"
)

type Config struct {
	Secrets secrets.Config

	Path     string `env:"DB_PATH" json:",omitempty"`
	User     string `env:"DB_USER" json:",omitempty"`
	Password string `env:"DB_PASSWORD" json:"-"`

	PoolMinConnections string        `env:"DB_POOL_MIN_CONNS" json:",omitempty"`
	PoolMaxConnections string        `env:"DB_POOL_MAX_CONNS" json:",omitempty"`
	PoolMaxConnLfe     time.Duration `env:"DB_POOL_MAX_CONN_LIFETIME, default=5m" json:",omitempty"`
	PoolMaxConnIdle    time.Duration `env:"DB_POOL_MAX_CONN_IDLE_TIME, default=1m" json:",omitempty"`
}

func (c *Config) DatabaseConfig() *Config {
	return c
}

func (c *Config) SecretManagerConfig() *secrets.Config {
	return &c.Secrets
}

func (c *Config) ConnectionURL() string {
	if c == nil {
		return ":memory:"
	}

	u := &url.URL{
		Scheme: "file",
		Opaque: c.Path,
	}

	params := []string{}
	if c.User != "" || c.Password != "" {
		params = append(params, "_auth")
	}
	if v := c.User; v != "" {
		params = append(params,
			fmt.Sprintf("_auth_user=%s", url.QueryEscape(v)),
		)
	}
	if v := c.Password; v != "" {
		params = append(params,
			fmt.Sprintf("_auth_pass=%s", url.QueryEscape(v)),
		)
	}
	params = append(params, "journal_mode=wal")
	u.RawQuery = strings.Join(params, "&")
	return u.String()
}

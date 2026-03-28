package sqlite

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sbuzas-jwl/go-pkgs/internal/secrets"
)

type Config struct {
	Secrets secrets.Config

	Path     string `env:"DB_PATH" json:",omitempty"`
	User     string `env:"DB_USER" json:",omitempty"`
	Password string `env:"DB_PASSWORD" json:"-"`
}

func (c *Config) DatabaseConfig() *Config {
	return c
}

func (c *Config) SecretManagerConfig() *secrets.Config {
	return &c.Secrets
}

func (c *Config) ConnectionURL() *url.URL {
	if c == nil {
		return &url.URL{
			Opaque: ":memory:",
		}
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
	return u
}

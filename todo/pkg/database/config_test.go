package database

import (
	"testing"

	"github.com/sbuzas-jwl/go-pkgs/todo/internal/project"
	"github.com/sethvargo/go-envconfig"
)

func TestConfig_DatabaseConfig(t *testing.T) {
	t.Parallel()

	cfg1 := &Config{}
	cfg2 := cfg1.DatabaseConfig()

	if cfg1 != cfg2 {
		t.Errorf("expected %#v to be %#v", cfg1, cfg2)
	}
}

func TestConfig_SecretManagerConfig(t *testing.T) {
	t.Parallel()

	cfg := &Config{}

	if smConfig := cfg.SecretManagerConfig(); smConfig == nil {
		t.Errorf("expected SecretManagerConfig to not be nil")
	}
}

func TestConfig_ConnectionURL(t *testing.T) {
	t.Parallel()

	ctx := project.TestContext(t)

	cases := []struct {
		name   string
		config *Config
		want   string
	}{
		{
			name:   "nil",
			config: nil,
			want:   ":memory:",
		},
		{
			name: "host",
			config: &Config{
				Path: "/var/data/mydb.sqlite",
			},
			want: "file:/var/data/mydb.sqlite?journal_mode=wal",
		},
		{
			name: "basic_auth",
			config: &Config{
				User:     "myuser",
				Password: "mypass",
			},
			want: "file:?_auth&_auth_user=myuser&_auth_pass=mypass&journal_mode=wal",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := tc.config
			if cfg != nil {
				if err := envconfig.Process(ctx, cfg); err != nil {
					t.Fatal(err)
				}
			}

			if got, want := cfg.ConnectionURL(), tc.want; got != want {
				t.Errorf("\ngot: %q \nwant: %q", got, want)
			}
		})
	}
}

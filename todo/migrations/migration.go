package migrations

import (
	"embed"
)

//go:embed sqlite
var sqliteMigrations embed.FS

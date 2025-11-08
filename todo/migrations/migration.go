package migrations

import (
	"embed"
)

//go:embed users/sqlite
var usersSqllite embed.FS

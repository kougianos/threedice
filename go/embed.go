// Package threedice embeds the assets the server ships with, so the binary is
// self-contained: no schema files or web root to deploy alongside it.
//
// The embeds live at the module root because //go:embed cannot reach outside
// its own package directory, and both trees are shared: db/migrations is also
// sqlc's schema input, and web/static is a copy of what Spring serves.
package threedice

import (
	"embed"
	"io/fs"
)

//go:embed all:db/migrations
var migrationsFS embed.FS

//go:embed all:web/static
var staticFS embed.FS

// Migrations holds the golang-migrate SQL, rooted so files sit at the top level.
func Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "db/migrations")
	if err != nil {
		panic("embed: db/migrations missing: " + err.Error())
	}
	return sub
}

// Static holds the frontend, rooted so index.html sits at the top level.
func Static() fs.FS {
	sub, err := fs.Sub(staticFS, "web/static")
	if err != nil {
		panic("embed: web/static missing: " + err.Error())
	}
	return sub
}

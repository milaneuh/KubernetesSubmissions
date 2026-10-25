package assets

import (
	"embed"
	"io/fs"
)

//go:embed all:static
var files embed.FS

var StaticFiles = sub(files, "static")

func sub(f embed.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		panic(err)
	}
	return sub
}

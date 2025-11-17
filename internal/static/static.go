package static

import (
	"embed"
	"io/fs"
)

/*
dist/css/* dist/js/* dist/components/* dist/utils/* dist/plugins/* dist/public/*
*/

//go:embed  dist
var DistDir embed.FS

// GetDistFS 返回嵌入的静态文件系统
func GetDistFS() fs.FS {
	distFS, err := fs.Sub(DistDir, "dist")
	if err != nil {
		panic(err)
	}
	return distFS
}

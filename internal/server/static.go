package server

import (
	"io/fs"
	"net/http"

	"github.com/biozz/paste/web"
)

func (h *Web) Static() http.Handler {
	fsys, _ := fs.Sub(web.StaticFS, "static")
	fileServer := http.FileServer(http.FS(fsys))
	if h.conf.Env == "dev" {
		fileServer = http.FileServer(http.Dir("web/static"))
	}
	return fileServer
}

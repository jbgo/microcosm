package web

import (
	"net/http"
	"path"
	"runtime"
)

type WebApp struct {
	Root string
}

func BuildApp() WebApp {
	app := WebApp{Root: findRootPath()}
	return app
}

func findRootPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename))
}

func (app WebApp) AssetsPath() string {
	return path.Join(app.Root, "assets")
}

func (app WebApp) AssetsHandler() http.Handler {
	fileServer := http.FileServer(http.Dir(app.AssetsPath()))
	return http.StripPrefix("/assets/", fileServer)
}

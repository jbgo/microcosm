package web

import (
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
)

type WebApp struct {
	Root string
	Log  *log.Logger
}

func BuildApp() WebApp {
	app := WebApp{Root: findRootPath()}
	app.Log = log.New(os.Stdout, "web ", log.LstdFlags)
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

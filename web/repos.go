package web

import (
	"net/http"
)

func (app WebApp) ListRepos(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, "main", "repos/index", nil)
}

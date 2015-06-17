package web

import (
	"net/http"
)

func (app WebApp) ListServices(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, "main", "services/index", nil)
}

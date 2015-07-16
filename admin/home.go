package admin

import (
	"net/http"
)

func (app WebApp) Home(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, "main", "home", nil)
}

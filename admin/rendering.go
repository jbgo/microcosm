package admin

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
)

var templateFuncs = template.FuncMap{
	"join":     strings.Join,
	"noslash":  func(s string) string { return s[1:] },
	"truncate": func(s string, length int) string { return s[0:length] },
	"printmap": printMap,
}

type LayoutData struct {
	Content template.HTML
}

func (app WebApp) RenderHTML(w http.ResponseWriter, layout, view string, data interface{}) {
	var buffer bytes.Buffer
	viewPath := path.Join(app.Root, "views", view+".html")
	viewTemplate := template.Must(template.New("").Funcs(templateFuncs).ParseFiles(viewPath))
	viewTemplate.ExecuteTemplate(&buffer, path.Base(viewPath), data)
	content := buffer.String()

	layoutPath := path.Join(app.Root, "views", "layouts", layout+".html")
	layoutTemplate := template.Must(template.New("").Funcs(templateFuncs).ParseFiles(layoutPath))
	layoutData := LayoutData{Content: template.HTML(content)}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	layoutTemplate.ExecuteTemplate(w, path.Base(layoutPath), layoutData)
}

func printMap(data map[string]string) string {
	return fmt.Sprintf("%s=%s %s=%s", "service", data["microcosm.service"], "type", data["microcosm.type"])
}

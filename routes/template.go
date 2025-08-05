package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type Template struct {
	template *template.Template
}

func ParseFiles(files ...string) (*Template, error) {
	tpl, err := template.New("").Funcs(template.FuncMap{
		"flash": func() string { return "" },
	}).ParseFiles(files...)
	if err != nil {
		return nil, fmt.Errorf("parsing view files: %w", err)
	}
	return &Template{
		template: tpl,
	}, nil
}

func (t *Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	clone := t.clone()

	cookie, err := r.Cookie("flash")
	if err == nil {
		clone.template = clone.template.Funcs(template.FuncMap{
			"flash": func() string { return cookie.Value },
		})

		expired := http.Cookie{
			Name:     "flash",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(w, &expired)
	}

	err = clone.template.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution failed: "+err.Error(), http.StatusInternalServerError)
	}
}

func (t *Template) clone() *Template {
	return &Template{
		template: t.template,
	}
}

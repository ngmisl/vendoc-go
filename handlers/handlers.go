package handlers

import (
	"html/template"
)

var templates *template.Template

func SetTemplates(tmpl *template.Template) {
	templates = tmpl
}
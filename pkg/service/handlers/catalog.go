package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/sergeygvozdev08101993/golang-local-library/internal/app"
	"github.com/sergeygvozdev08101993/golang-local-library/pkg/service/models"
)

var templateDirPath = ".././pkg/web/templates"

// Index выполняет обязанности пользовательского роутера.
func Index(w http.ResponseWriter, r *http.Request) {

	urlPath := r.URL.Path
	urlParts := strings.Split(urlPath, "/")

	if len(urlParts) == 2 {
		http.Redirect(w, r, "/catalog", http.StatusSeeOther)
	} else {
		if urlParts[2] == "book" && len(urlParts[3]) == 24 {
			GetBook(w, r)
			return
		} else if urlParts[2] == "author" && len(urlParts[3]) == 24 {
			GetAuthor(w, r)
			return
		} else if urlParts[2] == "genre" && len(urlParts[3]) == 24 {
			GetGenre(w, r)
			return
		} else if urlParts[2] == "bookinstance" && len(urlParts[3]) == 24 {
			GetBookInstance(w, r)
			return
		}
	}
}

// Catalog обрабатывает запрос для отображения информации о количестве документов, содержащих в каждой коллекции.
func Catalog(w http.ResponseWriter, r *http.Request) {

	var (
		result []models.Total
		total  models.Total

		err error
	)

	r.ParseForm()

	tmpl, err := template.ParseFiles(templateDirPath+"/index.gohtml", templateDirPath+"/main.gohtml")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse template files: %v", err)
		return
	}

	total, err = models.GetAmountAuthors()
	if err != nil {
		app.ErrLog.Printf("failed to get amount of all documents in authors collection: %v", err)
		renderTemplate(w, tmpl, nil, err.Error())
		return
	}
	result = append(result, total)

	total, err = models.GetAmountBooks()
	if err != nil {
		app.ErrLog.Printf("failed to get amount of all documents in books collection: %v", err)
		renderTemplate(w, tmpl, nil, err.Error())
		return
	}
	result = append(result, total)

	total, err = models.GetAmountGenres()
	if err != nil {
		app.ErrLog.Printf("failed to get amount of all documents in genres collection: %v", err)
		renderTemplate(w, tmpl, nil, err.Error())
		return
	}
	result = append(result, total)

	total, err = models.GetAmountBookInstances()
	if err != nil {
		app.ErrLog.Printf("failed to get amount of all documents in book instances collection: %v", err)
		renderTemplate(w, tmpl, nil, err.Error())
		return
	}
	result = append(result, total)

	total, err = models.GetAmountAvailableBookInstances()
	if err != nil {
		app.ErrLog.Printf("failed to get amount of all documents from available book instances: %v", err)
		renderTemplate(w, tmpl, nil, err.Error())
		return
	}
	result = append(result, total)

	renderTemplate(w, tmpl, result, "")
}

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, result []models.Total, err string) {

	d := models.Data{
		Title:     "Local Library Home",
		TotalList: result,
		Error:     err,
	}

	if err := tmpl.ExecuteTemplate(w, "index", d); err != nil {
		app.ErrLog.Printf("failed to render template file: %v", err)
	}
}

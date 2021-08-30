package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/sergeygvozdev08101993/golang-local-library/internal/app"
	"github.com/sergeygvozdev08101993/golang-local-library/pkg/service/models"
)

// Index выполняет обязанности пользовательского роутера.
func Index(w http.ResponseWriter, r *http.Request) {

	urlPath := r.URL.Path
	urlParts := strings.Split(urlPath, "/")

	if len(urlParts) == 2 {
		http.Redirect(w, r, "/catalog", http.StatusSeeOther)
	} else {
		if urlParts[2] == "book" && len(urlParts[3]) == 24 {

			customRouting(w, r, urlParts, GetBook, DeleteBook, UpdateBook,
				"/book_delete.gohtml", "/book_detail.gohtml")
			return

		} else if urlParts[2] == "author" && len(urlParts[3]) == 24 {

			customRouting(w, r, urlParts, GetAuthor, DeleteAuthor, UpdateAuthor,
				"/author_delete.gohtml", "/author_detail.gohtml")
			return

		} else if urlParts[2] == "genre" && len(urlParts[3]) == 24 {

			customRouting(w, r, urlParts, GetGenre, DeleteGenre, UpdateGenre,
				"/genre_delete.gohtml", "/genre_detail.gohtml")
			return

		} else if urlParts[2] == "bookinstance" && len(urlParts[3]) == 24 {

			customRouting(w, r, urlParts, GetBookInstance, DeleteBookInstance, UpdateBookInstance,
				"/bookinstance_delete.gohtml", "/bookinstance_detail.gohtml")
			return
		}
	}
}

// customRouting содержит реализацию роутинга для конкретного объекта данных.
func customRouting(w http.ResponseWriter, r *http.Request, urlParts []string,
	get models.GetHandler, delete models.DeleteHandler, update models.UpdateHandler,
	deleteTemplate string, detailTemplate string) {

	if len(urlParts) == 5 && urlParts[4] == "delete" {

		if r.Method == "GET" {
			get(w, r, deleteTemplate)
			return
		} else if r.Method == "POST" {
			delete(w, r)
			return
		}
	}

	if len(urlParts) == 5 && urlParts[4] == "update" {
		update(w, r)
		return
	}

	get(w, r, detailTemplate)
	return

}

// Catalog обрабатывает запрос для отображения информации о количестве документов, содержащих в каждой коллекции.
func Catalog(w http.ResponseWriter, r *http.Request) {

	var (
		result []models.Total
		total  models.Total

		err error
	)

	r.ParseForm()

	tmpl, err := template.ParseFiles(models.TemplateDirPath+"/index.gohtml", models.TemplateDirPath+"/main.gohtml")
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

// renderError предназначена для отображения статуса и сообщения об ошибке.
func renderError(w http.ResponseWriter, code int64, body string) {

	d := models.Err{
		Status:  code,
		Message: body,
	}

	tmpl, err := template.ParseFiles(models.TemplateDirPath + "/error.gohtml")
	if err != nil {
		app.ErrLog.Println(err)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "error", d); err != nil {
		app.ErrLog.Println(err)
		return
	}
}

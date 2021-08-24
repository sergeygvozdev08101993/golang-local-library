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

			if len(urlParts) == 5 && urlParts[4] == "delete" {

				if r.Method == "GET" {
					GetBook(w, r, "/book_delete.gohtml")
					return
				} else if r.Method == "POST" {
					DeleteBook(w, r)
					return
				}
			}

			if len(urlParts) == 5 && urlParts[4] == "update" {
				/*GetBook(w, r, "/book_delete.gohtml")
				return*/
			}

			GetBook(w, r, "/book_detail.gohtml")
			return

		} else if urlParts[2] == "author" && len(urlParts[3]) == 24 {

			if len(urlParts) == 5 && urlParts[4] == "delete" {

				if r.Method == "GET" {
					GetAuthor(w, r, "/author_delete.gohtml")
					return
				} else if r.Method == "POST" {
					DeleteAuthor(w, r)
					return
				}
			}

			if len(urlParts) == 5 && urlParts[4] == "update" {
				UpdateAuthor(w, r)
				return
			}

			GetAuthor(w, r, "/author_detail.gohtml")
			return

		} else if urlParts[2] == "genre" && len(urlParts[3]) == 24 {

			if len(urlParts) == 5 && urlParts[4] == "delete" {

				if r.Method == "GET" {
					GetGenre(w, r, "/genre_delete.gohtml")
					return
				} else if r.Method == "POST" {
					DeleteGenre(w, r)
					return
				}
			}

			if len(urlParts) == 5 && urlParts[4] == "update" {
				/*GetGenre(w, r, "/genre_delete.gohtml")
				return*/
			}

			GetGenre(w, r, "/genre_detail.gohtml")
			return

		} else if urlParts[2] == "bookinstance" && len(urlParts[3]) == 24 {

			if len(urlParts) == 5 && urlParts[4] == "delete" {

				if r.Method == "GET" {
					GetBookInstance(w, r, "/bookinstance_delete.gohtml")
					return
				} else if r.Method == "POST" {
					DeleteBookInstance(w, r)
					return
				}
			}

			if len(urlParts) == 5 && urlParts[4] == "update" {
				/*GetBookInstance(w, r, "/bookinstance_delete.gohtml")
				return*/
			}

			GetBookInstance(w, r, "/bookinstance_detail.gohtml")
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

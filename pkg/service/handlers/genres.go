package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/sergeygvozdev08101993/golang-local-library/internal/app"
	"github.com/sergeygvozdev08101993/golang-local-library/pkg/service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetGenre обрабатывает запрос для отображения информации по конкретному жанру.
func GetGenre(w http.ResponseWriter, r *http.Request, contentTemplate string) {

	var (
		genreID primitive.ObjectID
		err     error
	)

	urlPath := r.URL.Path
	urlParts := strings.Split(urlPath, "/")

	genreID, err = primitive.ObjectIDFromHex(urlParts[3])
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get genre Id: %v", err)
		return
	}

	genre, err := models.GetGenreByID(genreID)
	if err != nil && err.Error() == "mongo: no documents in result" {
		renderError(w, http.StatusNotFound, "Not Found")
		app.ErrLog.Printf("failed to get genre by id from database: %v", err)
		return
	}
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get genre by id from database: %v", err)
		return
	}

	books, err := models.GetBooksByGenreID(genreID)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get books by genre id from database: %v", err)
		return
	}

	tmpl, err := template.ParseFiles(templateDirPath+"/index.gohtml", templateDirPath+contentTemplate)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse template files: %v", err)
		return
	}

	d := models.Detail{
		Title: genre.Name,
		Genre: genre,
		Books: books,
	}
	if err = tmpl.ExecuteTemplate(w, "index", d); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to render template file: %v", err)
		return
	}
}

// ListGenres обрабатывает запрос для отображения всех жанров, содержащихся в коллекции.
func ListGenres(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(templateDirPath+"/index.gohtml", templateDirPath+"/genre_list.gohtml")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse template files: %v", err)
		return
	}

	genreList, err := models.GetListAllGenres()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get genre list from database: %v", err)
		return
	}

	d := models.Data{
		Title:     "Genre List",
		GenreList: genreList,
	}

	if err := tmpl.ExecuteTemplate(w, "index", d); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to render template file: %v", err)
		return
	}
}

// CreateGenre обрабатывает GET-запрос для отображения формы по созданию нового жанра,
// а также POST-запрос для обработки полученных данных из формы и добавления их в соответствующую коллекцию ДБ.
func CreateGenre(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		GetCreateGenre(w)
		break
	case "POST":
		PostCreateGenre(w, r)
		break
	}
}

// GetCreateGenre обрабатывает GET-запрос для отображения формы по созданию нового жанра.
func GetCreateGenre(w http.ResponseWriter) {

	var err error

	tmpl, err := template.ParseFiles(templateDirPath+"/index.gohtml", templateDirPath+"/genre_form.gohtml")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse template files: %v", err)
		return
	}

	d := models.Data{
		Title: "Create Genre",
	}
	if err = tmpl.ExecuteTemplate(w, "index", d); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to render template file: %v", err)
		return
	}
}

// PostCreateGenre обрабатывает POST-запрос из HTML-формы по созданию нового жанра,
// записывает данные автора в БД и перенаправляет на страницу жанра.
func PostCreateGenre(w http.ResponseWriter, r *http.Request) {

	var (
		name    string
		genreID string

		err error
	)

	if err = r.ParseForm(); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse genre form: %v", err)
		return
	}

	name = r.FormValue("name")

	name = strings.TrimSpace(name)
	if len(name) == 0 || len(name) < 3 {
		renderError(w, http.StatusBadRequest, "Bad Request")
		app.ErrLog.Println("failed to get genre name parameter: value is empty or amount of symbols is less than 3")
		return
	}

	genreID, err = models.CreateGenre(name)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to create genre: %v", err)
		return
	}

	redirectURL := "/catalog/genre/" + genreID
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// DeleteGenre обрабатывает POST-запрос из HTML-формы по удалению жанра,
// удаляет жанр из БД и перенаправляет на страницу, содержащую список жанров.
func DeleteGenre(w http.ResponseWriter, r *http.Request) {

	var (
		genreIDStr string
		genreID    primitive.ObjectID

		err error
	)

	if err = r.ParseForm(); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse genre delete form: %v", err)
		return
	}

	genreIDStr = r.FormValue("genreId")

	genreID, err = primitive.ObjectIDFromHex(genreIDStr)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get genre ID: %v", err)
		return
	}

	if err = models.DeleteGenreByID(genreID); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to delete genre by ID: %v", err)
		return
	}

	redirectURL := "/catalog/genres"
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/sergeygvozdev08101993/golang-local-library/internal/app"
	"github.com/sergeygvozdev08101993/golang-local-library/pkg/service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetBook обрабатывает запрос для отображения информации по конкретной книге.
func GetBook(w http.ResponseWriter, r *http.Request, contentTemplate string) {

	var (
		bookID primitive.ObjectID
		err    error
	)

	urlPath := r.URL.Path
	urlParts := strings.Split(urlPath, "/")

	bookID, err = primitive.ObjectIDFromHex(urlParts[3])
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get book ID: %v", err)
		return
	}

	book, err := models.GetBookByID(bookID)
	if err != nil && err.Error() == "mongo: no documents in result" {
		renderError(w, http.StatusNotFound, "Not Found")
		app.ErrLog.Printf("failed to get book by ID from database: %v", err)
		return
	}
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get book by ID from database: %v", err)
		return
	}

	bookInstances, err := models.GetBookInstancesByBookID(bookID)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get book instance by ID from database: %v", err)
		return
	}

	tmpl, err := template.ParseFiles(templateDirPath+"/index.gohtml", templateDirPath+contentTemplate)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse template files: %v", err)
		return
	}

	d := models.Detail{
		Title:         book.Name,
		Book:          book,
		BookInstances: bookInstances,
	}
	if err = tmpl.ExecuteTemplate(w, "index", d); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to render template file: %v", err)
		return
	}
}

// ListBooks обрабатывает запрос для отображения всех книг, содержащихся в коллекции.
func ListBooks(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(templateDirPath+"/index.gohtml", templateDirPath+"/book_list.gohtml")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse template files: %v", err)
		return
	}

	bookList, err := models.GetListAllBooks()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get book list from database: %v", err)
		return
	}

	d := models.Data{
		Title:    "Book List",
		BookList: bookList,
	}
	if err := tmpl.ExecuteTemplate(w, "index", d); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to render template file: %v", err)
		return
	}
}

// CreateBook обрабатывает GET-запрос для отображения формы по созданию новой книги,
// а также POST-запрос для обработки полученных данных из формы и добавления их в соответствующую коллекцию ДБ.
func CreateBook(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		GetCreateBook(w)
		break
	case "POST":
		PostCreateBook(w, r)
		break
	}
}

// GetCreateBook обрабатывает GET-запрос для отображения формы по созданию новой книги.
func GetCreateBook(w http.ResponseWriter) {

	var (
		authorList []models.AuthorData
		genreList  []models.GenreData
		err        error
	)

	authorList, err = models.GetListAllAuthors()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get author list from database: %v", err)
		return
	}

	genreList, err = models.GetListAllGenres()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get genre list from database: %v", err)
		return
	}

	tmpl, err := template.ParseFiles(templateDirPath+"/index.gohtml", templateDirPath+"/book_form.gohtml")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse template files: %v", err)
		return
	}

	d := models.Data{
		Title:      "Create Book",
		AuthorList: authorList,
		GenreList:  genreList,
	}
	if err = tmpl.ExecuteTemplate(w, "index", d); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to render template file: %v", err)
		return
	}
}

// PostCreateBook обрабатывает POST-запрос из HTML-формы по созданию новой книги,
// записывает данные книги в БД и перенаправляет на страницу книги.
func PostCreateBook(w http.ResponseWriter, r *http.Request) {

	var (
		title, summary, isbn string
		authorIDStr          string
		genresIDStr          []string
		authorID             primitive.ObjectID
		genreIDs             []primitive.ObjectID
		bookID               string

		err error
	)

	if err = r.ParseForm(); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse book form: %v", err)
		return
	}

	title = r.FormValue("title")
	summary = r.FormValue("summary")
	isbn = r.FormValue("isbn")

	authorIDStr = r.FormValue("author")
	genresIDStr = r.Form["genre"]

	tmpTitle := strings.TrimSpace(title)
	if len(tmpTitle) == 0 {
		renderError(w, http.StatusBadRequest, "Bad Request")
		app.ErrLog.Println("failed to get the book title parameter")
		return
	}

	tmpSummary := strings.TrimSpace(summary)
	if len(tmpSummary) == 0 {
		renderError(w, http.StatusBadRequest, "Bad Request")
		app.ErrLog.Println("failed to get the book summary parameter")
		return
	}

	isbn = strings.TrimSpace(isbn)
	if len(isbn) == 0 {
		renderError(w, http.StatusBadRequest, "Bad Request")
		app.ErrLog.Println("failed to get the book ISBN parameter")
		return
	}

	if len(authorIDStr) == 0 {
		renderError(w, http.StatusBadRequest, "Bad Request")
		app.ErrLog.Println("failed to get author ID parameter")
		return
	}

	if len(genresIDStr) == 0 {
		renderError(w, http.StatusBadRequest, "Bad Request")
		app.ErrLog.Println("failed to get genre IDs parameters")
		return
	}

	authorID, err = primitive.ObjectIDFromHex(authorIDStr)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get author ID: %v", err)
		return
	}

	for _, genreIDStr := range genresIDStr {

		genreID, err := primitive.ObjectIDFromHex(genreIDStr)
		if err != nil {
			renderError(w, http.StatusInternalServerError, "Internal Server Error")
			app.ErrLog.Printf("failed to get genre ID: %v", err)
			break
		}

		genreIDs = append(genreIDs, genreID)
	}

	bookID, err = models.CreateBook(title, summary, isbn, authorID, genreIDs)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to create book: %v", err)
		return
	}

	redirectURL := "/catalog/book/" + bookID
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// renderError предназначена для отображения статуса и сообщения об ошибке.
func renderError(w http.ResponseWriter, code int64, body string) {

	d := models.Err{
		Status:  code,
		Message: body,
	}

	tmpl, err := template.ParseFiles(templateDirPath + "/error.gohtml")
	if err != nil {
		app.ErrLog.Println(err)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "error", d); err != nil {
		app.ErrLog.Println(err)
		return
	}
}

// DeleteBook обрабатывает POST-запрос из HTML-формы по удалению данных книги,
// удаляет данные книги из БД и перенаправляет на страницу, содержащую список книг.
func DeleteBook(w http.ResponseWriter, r *http.Request) {

	var (
		bookIDStr string
		bookID    primitive.ObjectID

		err error
	)

	if err = r.ParseForm(); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse book delete form: %v", err)
		return
	}

	bookIDStr = r.FormValue("bookId")

	bookID, err = primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get book ID: %v", err)
		return
	}

	if err = models.DeleteBookByID(bookID); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to delete book by ID: %v", err)
		return
	}

	redirectURL := "/catalog/books"
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

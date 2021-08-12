package handlers

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/sergeygvozdev08101993/golang-local-library/internal/app"
	"github.com/sergeygvozdev08101993/golang-local-library/pkg/service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetBookInstance обрабатывает запрос для отображения информации по конкретному экземпляру книги.
func GetBookInstance(w http.ResponseWriter, r *http.Request) {

	var (
		bookInstanceID primitive.ObjectID
		err            error
	)

	urlPath := r.URL.Path
	urlParts := strings.Split(urlPath, "/")

	bookInstanceID, err = primitive.ObjectIDFromHex(urlParts[3])
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get book instance ID: %v", err)
		return
	}

	bookInstance, err := models.GetBookInstanceByID(bookInstanceID)
	if err != nil && err.Error() == "mongo: no documents in result" {
		renderError(w, http.StatusNotFound, "Not Found")
		app.ErrLog.Printf("failed to get book instance by ID from database: %v", err)
		return
	}
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get book instance by ID from database: %v", err)
		return
	}

	tmpl, err := template.ParseFiles(templateDirPath+"/index.gohtml", templateDirPath+"/bookinstance_detail.gohtml")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse template files: %v", err)
		return
	}

	d := models.Detail{
		Title:        "Book Instance: " + bookInstance.ID,
		BookInstance: bookInstance,
	}
	if err = tmpl.ExecuteTemplate(w, "index", d); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to render template file: %v", err)
		return
	}
}

// ListBookInstances обрабатывает запрос для отображения всех экземпляров книг, содержащихся в коллекции.
func ListBookInstances(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(templateDirPath+"/index.gohtml", templateDirPath+"/bookinstance_list.gohtml")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse template files: %v", err)
		return
	}

	bookInstanceList, err := models.GetListAllBookInstances()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get book instance list from database: %v", err)
		return
	}

	d := models.Data{
		Title:            "Book Instance List",
		BookInstanceList: bookInstanceList,
	}

	if err := tmpl.ExecuteTemplate(w, "index", d); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to render template file: %v", err)
		return
	}
}

// CreateBookInstance обрабатывает GET-запрос для отображения формы по созданию нового экземпляра книги,
// а также POST-запрос для обработки полученных данных из формы и добавления их в соответствующую коллекцию ДБ.
func CreateBookInstance(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		GetCreateBookInstance(w)
		break
	case "POST":
		PostCreateBookInstance(w, r)
		break
	}
}

// GetCreateBookInstance обрабатывает GET-запрос для отображения формы по созданию нового экземпляра книги.
func GetCreateBookInstance(w http.ResponseWriter) {

	var (
		bookList []models.BookData
		err      error
	)

	bookList, err = models.GetListAllBooks()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get book list from database: %v", err)
		return
	}

	tmpl, err := template.ParseFiles(templateDirPath+"/index.gohtml", templateDirPath+"/bookinstance_form.gohtml")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse template files: %v", err)
		return
	}

	d := models.Data{
		Title:    "Create BookInstance",
		BookList: bookList,
	}
	if err = tmpl.ExecuteTemplate(w, "index", d); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to render template file: %v", err)
		return
	}
}

// PostCreateBookInstance обрабатывает POST-запрос из HTML-формы по нового экземпляра книги,
// записывает данные нового экземпляра книги в БД и перенаправляет на страницу экземпляра книги.
func PostCreateBookInstance(w http.ResponseWriter, r *http.Request) {

	var (
		imprint, status, dueBackStr, bookIDStr string
		bookID                                 primitive.ObjectID
		dueBackTime                            time.Time
		dueBack                                interface{}
		bookInstanceID                         string

		err error
	)

	if err = r.ParseForm(); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse book instance form: %v", err)
		return
	}

	bookIDStr = r.FormValue("book")
	imprint = r.FormValue("imprint")
	status = r.FormValue("status")
	dueBackStr = r.FormValue("due_back")

	if len(bookIDStr) == 0 {
		renderError(w, http.StatusBadRequest, "Bad Request")
		app.ErrLog.Println("failed to get book ID parameter")
		return
	}

	bookID, err = primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get book ID: %v", err)
		return
	}

	tmpImprint := strings.TrimSpace(imprint)
	if len(tmpImprint) == 0 {
		renderError(w, http.StatusBadRequest, "Bad Request")
		app.ErrLog.Println("failed to get the book instance imprint parameter")
		return
	}

	if len(dueBackStr) == 0 {
		dueBack = nil
	} else {
		dueBackTime, err = time.Parse("2006-01-02", dueBackStr)
		if err != nil {
			renderError(w, http.StatusInternalServerError, "Internal Server Error")
			app.ErrLog.Printf("failed to parse book instance due back parameter from bookinstance-form: %v", err)
			return
		}

		dueBack = primitive.NewDateTimeFromTime(dueBackTime)
	}

	if len(status) == 0 {
		renderError(w, http.StatusBadRequest, "Bad Request")
		app.ErrLog.Println("failed to get book instance status parameter")
		return
	}

	bookInstanceID, err = models.CreateBookInstance(imprint, status, bookID, dueBack)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to create book: %v", err)
		return
	}

	redirectURL := "/catalog/bookinstance/" + bookInstanceID
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)

}

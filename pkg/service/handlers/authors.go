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

// GetAuthor обрабатывает запрос для отображения информации по конкретному автору.
func GetAuthor(w http.ResponseWriter, r *http.Request, contentTemplate string) {

	var (
		authorID primitive.ObjectID
		err      error
	)

	urlPath := r.URL.Path
	urlParts := strings.Split(urlPath, "/")

	authorID, err = primitive.ObjectIDFromHex(urlParts[3])
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get author ID: %v", err)
		return
	}

	author, err := models.GetAuthorByID(authorID)
	if err != nil && err.Error() == "mongo: no documents in result" {
		renderError(w, http.StatusNotFound, "Not Found")
		app.ErrLog.Printf("failed to get author by ID from database: %v", err)
		return
	}
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get author by ID from database: %v", err)
		return
	}

	books, err := models.GetBooksByAuthorID(authorID)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get books by author ID from database: %v", err)
		return
	}

	tmpl, err := template.ParseFiles(templateDirPath+"/index.gohtml", templateDirPath+contentTemplate)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse template files: %v", err)
		return
	}

	if len(author.DateBirth) != 0 {
		dateBirthTmp, err := time.Parse("2006-01-02", author.DateBirth)
		if err != nil {
			renderError(w, http.StatusInternalServerError, "Internal Server Error")
			app.ErrLog.Printf("failed to parse the author date of birth: %v", err)
			return
		}
		author.DateBirth = dateBirthTmp.Format("2 Jan, 2006")
	}

	if len(author.DateDeath) != 0 {
		dateDeathTmp, err := time.Parse("2006-01-02", author.DateDeath)
		if err != nil {
			renderError(w, http.StatusInternalServerError, "Internal Server Error")
			app.ErrLog.Printf("failed to parse the author date of death: %v", err)
			return
		}
		author.DateDeath = dateDeathTmp.Format("2 Jan, 2006")
	}

	d := models.Detail{
		Title:  author.FamilyName + ", " + author.FirstName,
		Author: author,
		Books:  books,
	}
	if err = tmpl.ExecuteTemplate(w, "index", d); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to render template file: %v", err)
		return
	}
}

// ListAuthors обрабатывает запрос для отображения всех авторов, содержащихся в коллекции.
func ListAuthors(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(templateDirPath+"/index.gohtml", templateDirPath+"/author_list.gohtml")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse template files: %v", err)
		return
	}

	authorList, err := models.GetListAllAuthors()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get author list from database: %v", err)
		return
	}

	d := models.Data{
		Title:      "Author List",
		AuthorList: authorList,
	}

	if err := tmpl.ExecuteTemplate(w, "index", d); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to render template file: %v", err)
		return
	}
}

// CreateAuthor обрабатывает GET-запрос для отображения формы по созданию нового автора,
// а также POST-запрос для обработки полученных данных из формы и добавления их в соответствующую коллекцию БД.
func CreateAuthor(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		GetCreateAuthor(w)
		break
	case "POST":
		PostCreateAuthor(w, r)
		break
	}
}

// GetCreateAuthor обрабатывает GET-запрос для отображения формы по созданию нового автора.
func GetCreateAuthor(w http.ResponseWriter) {

	tmpl, err := template.ParseFiles(templateDirPath+"/index.gohtml", templateDirPath+"/author_form.gohtml")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse template files: %v", err)
		return
	}

	d := models.Data{
		Title: "Create Author",
	}
	if err = tmpl.ExecuteTemplate(w, "index", d); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to render template file: %v", err)
		return
	}
}

// PostCreateAuthor обрабатывает POST-запрос из HTML-формы по созданию нового автора,
// записывает данные автора в БД и перенаправляет на страницу автора.
func PostCreateAuthor(w http.ResponseWriter, r *http.Request) {

	var (
		firstName, familyName, birthDateStr, deathDateStr string
		birthDateTime, deathDateTime                      time.Time
		birthDate, deathDate                              interface{}
		authorID                                          string

		err error
	)

	if err = r.ParseForm(); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse author create form: %v", err)
		return
	}

	firstName = r.FormValue("first_name")
	familyName = r.FormValue("family_name")
	birthDateStr = r.FormValue("date_of_birth")
	deathDateStr = r.FormValue("date_of_death")

	firstName = strings.TrimSpace(firstName)
	if len(firstName) == 0 {
		renderError(w, http.StatusBadRequest, "Bad Request")
		app.ErrLog.Println("failed to get author first name parameter")
		return
	}

	familyName = strings.TrimSpace(familyName)
	if len(familyName) == 0 {
		renderError(w, http.StatusBadRequest, "Bad Request")
		app.ErrLog.Println("failed to get author family name parameter")
		return
	}

	if len(birthDateStr) == 0 {
		birthDate = nil
	} else {
		birthDateTime, err = time.Parse("2006-01-02", birthDateStr)
		if err != nil {
			renderError(w, http.StatusInternalServerError, "Internal Server Error")
			app.ErrLog.Printf("failed to parse author birth date parameter from author-form: %v", err)
			return
		}

		birthDate = primitive.NewDateTimeFromTime(birthDateTime)
	}

	if len(deathDateStr) == 0 {
		deathDate = nil
	} else {

		deathDateTime, err = time.Parse("2006-01-02", deathDateStr)
		if err != nil {
			renderError(w, http.StatusInternalServerError, "Internal Server Error")
			app.ErrLog.Printf("failed to parse author death date parameter from author-form: %v", err)
			return
		}

		deathDate = primitive.NewDateTimeFromTime(deathDateTime)
	}

	authorID, err = models.CreateAuthor(firstName, familyName, birthDate, deathDate)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to create author: %v", err)
		return
	}

	redirectURL := "/catalog/author/" + authorID
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// DeleteAuthor обрабатывает POST-запрос из HTML-формы по удалению данных автора,
// удаляет данные автора из БД и перенаправляет на страницу, содержащую список авторов.
func DeleteAuthor(w http.ResponseWriter, r *http.Request) {

	var (
		authorIDStr string
		authorID    primitive.ObjectID

		err error
	)

	if err = r.ParseForm(); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse author delete form: %v", err)
		return
	}

	authorIDStr = r.FormValue("authorId")

	authorID, err = primitive.ObjectIDFromHex(authorIDStr)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get author ID: %v", err)
		return
	}

	if err = models.DeleteAuthorByID(authorID); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to delete author by ID: %v", err)
		return
	}

	redirectURL := "/catalog/authors"
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// UpdateAuthor обрабатывает GET-запрос для отображения формы по обновлению данных автора,
// а также POST-запрос для обновления полученных данных в БД.
func UpdateAuthor(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		GetUpdateAuthor(w, r)
		break
	case "POST":
		PostUpdateAuthor(w, r)
		break
	}
}

// GetUpdateAuthor обрабатывает GET-запрос для отображения формы по обновлению данных автора.
func GetUpdateAuthor(w http.ResponseWriter, r *http.Request) {

	var (
		authorID primitive.ObjectID
		err      error
	)

	urlPath := r.URL.Path
	urlParts := strings.Split(urlPath, "/")

	authorID, err = primitive.ObjectIDFromHex(urlParts[3])
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get author ID: %v", err)
		return
	}

	author, err := models.GetAuthorByID(authorID)
	if err != nil && err.Error() == "mongo: no documents in result" {
		renderError(w, http.StatusNotFound, "Not Found")
		app.ErrLog.Printf("failed to get author by ID from database: %v", err)
		return
	}
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get author by ID from database: %v", err)
		return
	}

	tmpl, err := template.ParseFiles(templateDirPath+"/index.gohtml", templateDirPath+"/author_form.gohtml")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse template files: %v", err)
		return
	}

	d := models.Detail{
		Title:  "Update Author",
		Author: author,
	}
	if err = tmpl.ExecuteTemplate(w, "index", d); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to render template file: %v", err)
		return
	}
}

// PostUpdateAuthor обрабатывает POST-запрос из HTML-формы по обновлению данных автора,
// обновляет данные в БД и перенаправляет на страницу автора.
func PostUpdateAuthor(w http.ResponseWriter, r *http.Request) {

	var (
		firstName, familyName, birthDateStr, deathDateStr string
		birthDateTime, deathDateTime                      time.Time
		birthDate, deathDate                              interface{}
		authorIDStr                                       string
		authorID                                          primitive.ObjectID

		err error
	)

	urlPath := r.URL.Path
	urlParts := strings.Split(urlPath, "/")
	authorIDStr = urlParts[3]

	authorID, err = primitive.ObjectIDFromHex(authorIDStr)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to get author ID: %v", err)
		return
	}

	if err = r.ParseForm(); err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to parse author create form: %v", err)
		return
	}

	firstName = r.FormValue("first_name")
	familyName = r.FormValue("family_name")
	birthDateStr = r.FormValue("date_of_birth")
	deathDateStr = r.FormValue("date_of_death")

	firstName = strings.TrimSpace(firstName)
	if len(firstName) == 0 {
		renderError(w, http.StatusBadRequest, "Bad Request")
		app.ErrLog.Println("failed to get author first name parameter")
		return
	}

	familyName = strings.TrimSpace(familyName)
	if len(familyName) == 0 {
		renderError(w, http.StatusBadRequest, "Bad Request")
		app.ErrLog.Println("failed to get author family name parameter")
		return
	}

	if len(birthDateStr) == 0 {
		birthDate = nil
	} else {
		birthDateTime, err = time.Parse("2006-01-02", birthDateStr)
		if err != nil {
			renderError(w, http.StatusInternalServerError, "Internal Server Error")
			app.ErrLog.Printf("failed to parse author birth date parameter from author-form: %v", err)
			return
		}

		birthDate = primitive.NewDateTimeFromTime(birthDateTime)
	}

	if len(deathDateStr) == 0 {
		deathDate = nil
	} else {

		deathDateTime, err = time.Parse("2006-01-02", deathDateStr)
		if err != nil {
			renderError(w, http.StatusInternalServerError, "Internal Server Error")
			app.ErrLog.Printf("failed to parse author death date parameter from author-form: %v", err)
			return
		}

		deathDate = primitive.NewDateTimeFromTime(deathDateTime)
	}

	err = models.UpdateAuthorByID(authorID, firstName, familyName, birthDate, deathDate)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error")
		app.ErrLog.Printf("failed to update author: %v", err)
		return
	}

	redirectURL := "/catalog/author/" + authorIDStr
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

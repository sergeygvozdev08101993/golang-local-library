package models

import "net/http"

// TemplateDirPath содержит путь к директории, где хранятся шаблоны.
const TemplateDirPath = ".././pkg/web/templates"

// GetHandler описывает тип обработчика, предоставляющего пользователю данные.
type GetHandler func(http.ResponseWriter, *http.Request, string)

// DeleteHandler описывает тип обработчика, удаляющего данные из БД.
type DeleteHandler func(http.ResponseWriter, *http.Request)

// UpdateHandler описывает тип обработчика, обновляющего данные в БД.
type UpdateHandler func(http.ResponseWriter, *http.Request)

// Data описывает структуру данных, предназначенную для
// вывода списка всех документов и их общего числа в каждой коллекции.
type Data struct {
	Title            string
	AuthorList       interface{}
	BookList         interface{}
	BookInstanceList interface{}
	GenreList        interface{}
	TotalList        interface{}
	Error            string
}

// Detail описывает структуру данных, предназначенную для
//	вывода информации по конкретному документу.
type Detail struct {
	Title         string
	Book          BookData
	Books         []BookData
	BookInstance  BookInstanceData
	BookInstances []BookInstanceData
	Author        interface{}
	Genre         interface{}
}

// Total описывает структуру данных, предназначенной для
// хранения информации о количестве документов в каждой коллекции.
type Total struct {
	DocName string
	Amount  int64
}

// Err описывает структуру данных полученной ошибки.
type Err struct {
	Status  int64
	Message string
}

// AuthorData описывает структуру данных коллекции authors.
type AuthorData struct {
	ID         string
	FirstName  string
	FamilyName string
	DateBirth  string
	DateDeath  string
}

// BookData описывает структуру данных коллекции books.
type BookData struct {
	ID      string
	Genre   interface{}
	Name    string
	Author  interface{}
	Summary string
	Isbn    string
}

// GenreData описывает структуру данных коллекции genres.
type GenreData struct {
	ID   string
	Name string
}

// BookInstanceData описывает структуру данных коллекции bookinstances.
type BookInstanceData struct {
	ID       string
	Status   string
	BookName string
	BookID   string
	Imprint  string
	DueBack  string
}

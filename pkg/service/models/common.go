package models

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
	Genre         GenreData
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

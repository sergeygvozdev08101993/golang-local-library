package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sergeygvozdev08101993/golang-local-library/internal/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetBookByID получает данные книги по ее ID.
func GetBookByID(bookID primitive.ObjectID) (BookData, error) {

	var (
		book    BookData
		tmpBook bson.D
		genres  []GenreData

		err error
	)

	libDB := app.ClientDB.Database("local_library")
	booksCollection := libDB.Collection("books")

	ctx := context.TODO()
	err = booksCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: bookID}}).Decode(&tmpBook)
	if err != nil {
		return book, err
	}
	if err == mongo.ErrNoDocuments {
		return book, err
	}

	for _, genreID := range tmpBook[1].Value.(primitive.A) {

		genre, err := GetGenreByID(genreID.(primitive.ObjectID))
		if err != nil && err.Error() == "mongo: no documents in result" {
			break
		}
		if err != nil {
			break
		}

		genres = append(genres, GenreData{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}

	switch "author" {
	case tmpBook[3].Key:
		authorID := tmpBook[3].Value.(primitive.ObjectID)
		summary := tmpBook[4].Value.(string)

		author, err := GetAuthorByID(authorID)
		if err != nil && err.Error() == "mongo: no documents in result" {
			return book, err
		}
		if err != nil {
			return book, err
		}

		book = getResultBook(tmpBook, bookID, author, summary, genres)
		break

	case tmpBook[4].Key:
		authorID := tmpBook[4].Value.(primitive.ObjectID)
		summary := tmpBook[3].Value.(string)

		author, err := GetAuthorByID(authorID)
		if err != nil && err.Error() == "mongo: no documents in result" {
			return book, err
		}
		if err != nil {
			return book, err
		}

		book = getResultBook(tmpBook, bookID, author, summary, genres)
		break
	}

	return book, nil
}

func getResultBook(tmpBook bson.D, bookID primitive.ObjectID, author AuthorData, summary string, genres []GenreData) BookData {

	return BookData{
		ID:   bookID.Hex(),
		Name: tmpBook[2].Value.(string),
		Author: AuthorData{
			ID:   author.ID,
			Name: author.Name,
		},
		Summary: summary,
		Isbn:    tmpBook[5].Value.(string),
		Genre:   genres,
	}
}

// GetBooksByAuthorID получает список книг по ID автора.
func GetBooksByAuthorID(authorID primitive.ObjectID) ([]BookData, error) {

	var (
		tmpBooks []bson.M
		books    []BookData

		err error
	)

	libDB := app.ClientDB.Database("local_library")
	booksCollection := libDB.Collection("books")

	ctx := context.TODO()
	cursor, err := booksCollection.Find(ctx, bson.D{primitive.E{Key: "author", Value: authorID}})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &tmpBooks); err != nil {
		return nil, err
	}

	for _, tmpBook := range tmpBooks {

		books = append(books, BookData{
			ID:      tmpBook["_id"].(primitive.ObjectID).Hex(),
			Name:    tmpBook["title"].(string),
			Summary: tmpBook["summary"].(string),
		})
	}

	return books, nil
}

// GetBooksByGenreID получает список книг по ID жанра.
func GetBooksByGenreID(genreID primitive.ObjectID) ([]BookData, error) {

	var (
		tmpBooks []bson.M
		books    []BookData

		err error
	)

	libDB := app.ClientDB.Database("local_library")
	booksCollection := libDB.Collection("books")

	ctx := context.TODO()
	cursor, err := booksCollection.Find(ctx, bson.D{primitive.E{Key: "genre", Value: genreID}})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &tmpBooks); err != nil {
		return nil, err
	}

	for _, tmpBook := range tmpBooks {

		books = append(books, BookData{
			ID:      tmpBook["_id"].(primitive.ObjectID).Hex(),
			Name:    tmpBook["title"].(string),
			Summary: tmpBook["summary"].(string),
		})
	}

	return books, nil
}

// GetListAllBooks получает список всех книг, отсортированный по названию в алфавитном порядке.
func GetListAllBooks() ([]BookData, error) {

	var (
		tmpBooks []bson.D
		books    []BookData

		err error
	)

	libDB := app.ClientDB.Database("local_library")
	booksCollection := libDB.Collection("books")

	ctx := context.TODO()
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "title", Value: 1}})

	cursor, err := booksCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &tmpBooks); err != nil {
		return nil, err
	}

	authorList, err := GetListAllAuthors()
	if err != nil {
		return nil, err
	}

	for _, book := range tmpBooks {
		for _, author := range authorList {

			var authorIDFromBooksCollection string

			id := book[0].Value.(primitive.ObjectID).Hex()
			name := book[2].Value.(string)

			switch "author" {
			case book[3].Key:
				authorIDFromBooksCollection = book[3].Value.(primitive.ObjectID).Hex()
				break
			case book[4].Key:
				authorIDFromBooksCollection = book[4].Value.(primitive.ObjectID).Hex()
				break
			}

			authorIDFromAuthorsCollection := author.ID
			if authorIDFromBooksCollection == authorIDFromAuthorsCollection {

				books = append(books, BookData{
					ID:     id,
					Name:   name,
					Author: author.Name,
				})

				break
			}
		}
	}

	return books, nil
}

//GetAmountBooks получает число книг, содержащихся в коллекции.
func GetAmountBooks() (Total, error) {

	var (
		result Total
		amount int64

		err error
	)

	libDB := app.ClientDB.Database("local_library")
	ctx := context.TODO()

	amount, err = libDB.Collection("books").CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, fmt.Errorf("GetAmountBooks: %v", err)
	}

	result = Total{DocName: "Books", Amount: amount}

	return result, nil
}

// CreateBook записывает данные новой книги в БД.
func CreateBook(title, summary, isbn string, authorID primitive.ObjectID, genreIDs []primitive.ObjectID) (string, error) {

	var (
		res    *mongo.InsertOneResult
		bookID string

		err error
	)

	libDB := app.ClientDB.Database("local_library")
	booksCollection := libDB.Collection("books")

	ctx := context.TODO()
	res, err = booksCollection.InsertOne(ctx, bson.D{
		primitive.E{Key: "genre", Value: genreIDs},
		primitive.E{Key: "title", Value: title},
		primitive.E{Key: "summary", Value: summary},
		primitive.E{Key: "author", Value: authorID},
		primitive.E{Key: "isbn", Value: isbn},
	})
	if err != nil {
		return bookID, err
	}

	bookID = res.InsertedID.(primitive.ObjectID).Hex()
	return bookID, nil
}

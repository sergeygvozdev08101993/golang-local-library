package models

import (
	"context"
	"fmt"

	"github.com/sergeygvozdev08101993/golang-local-library/internal/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetBookInstanceByID получает данные экземпляра книги по его ID.
func GetBookInstanceByID(bookInstanceID primitive.ObjectID) (BookInstanceData, error) {

	var (
		tmpBookInstance bson.M
		bookInstance    BookInstanceData
		dueBack         string

		err error
	)

	libDB := app.ClientDB.Database("local_library")
	bookInstancesCollection := libDB.Collection("bookinstances")

	ctx := context.TODO()
	err = bookInstancesCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: bookInstanceID}}).Decode(&tmpBookInstance)
	if err != nil {
		return bookInstance, err
	}
	if err == mongo.ErrNoDocuments {
		return bookInstance, err
	}

	bookID := tmpBookInstance["book"].(primitive.ObjectID)
	book, err := GetBookByID(bookID)
	if err != nil && err.Error() == "mongo: no documents in result" {
		return bookInstance, err
	}
	if err != nil {
		return bookInstance, err
	}

	if tmpBookInstance["due_back"] != nil {
		dueBack = tmpBookInstance["due_back"].(primitive.DateTime).Time().Format("2 Jan, 2006")
	} else {
		dueBack = ""
	}

	bookInstance = BookInstanceData{
		ID:       tmpBookInstance["_id"].(primitive.ObjectID).Hex(),
		Status:   tmpBookInstance["status"].(string),
		BookID:   book.ID,
		BookName: book.Name,
		Imprint:  tmpBookInstance["imprint"].(string),
		DueBack:  dueBack,
	}

	return bookInstance, nil
}

// GetBookInstancesByBookID получает данные экземпляра книги по ID книги.
func GetBookInstancesByBookID(bookID primitive.ObjectID) ([]BookInstanceData, error) {

	var (
		bookInstances    []BookInstanceData
		tmpBookInstances []bson.M

		err error
	)

	libDB := app.ClientDB.Database("local_library")
	bookInstancesCollection := libDB.Collection("bookinstances")

	ctx := context.TODO()
	cursor, err := bookInstancesCollection.Find(ctx, bson.D{primitive.E{Key: "book", Value: bookID}})
	if err != nil {
		return bookInstances, err
	}

	err = cursor.All(ctx, &tmpBookInstances)
	if err != nil {
		return bookInstances, err
	}

	for _, tmpBookInstance := range tmpBookInstances {

		var dueBack string

		if tmpBookInstance["due_back"] != nil {
			dueBack = tmpBookInstance["due_back"].(primitive.DateTime).Time().Format("2 Jan, 2006")
		} else {
			dueBack = ""
		}

		bookInstances = append(bookInstances, BookInstanceData{
			ID:      tmpBookInstance["_id"].(primitive.ObjectID).Hex(),
			Status:  tmpBookInstance["status"].(string),
			Imprint: tmpBookInstance["imprint"].(string),
			DueBack: dueBack,
		})
	}

	return bookInstances, nil
}

// GetListAllBookInstances получает список экземпляров книг.
func GetListAllBookInstances() ([]BookInstanceData, error) {

	var (
		tmpBookInstances []bson.D
		bookInstances    []BookInstanceData
		err              error
	)

	libDB := app.ClientDB.Database("local_library")
	bookInstancesCollection := libDB.Collection("bookinstances")

	ctx := context.TODO()
	cursor, err := bookInstancesCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &tmpBookInstances)
	if err != nil {
		return nil, err
	}

	bookList, err := GetListAllBooks()
	if err != nil {
		return nil, err
	}

	for _, bookInstance := range tmpBookInstances {
		for _, book := range bookList {

			var dueBack string

			id := bookInstance[0].Value.(primitive.ObjectID).Hex()
			status := bookInstance[1].Value.(string)
			bookIDFromBookInstanceCollection := bookInstance[2].Value.(primitive.ObjectID).Hex()
			imprint := bookInstance[3].Value.(string)

			if bookInstance[4].Value != nil {
				dueBack = bookInstance[4].Value.(primitive.DateTime).Time().Format("2 Jan, 2006")
			} else {
				dueBack = ""
			}

			bookIDFromBookCollection := book.ID
			if bookIDFromBookInstanceCollection == bookIDFromBookCollection {

				bookInstances = append(bookInstances, BookInstanceData{
					ID:       id,
					Status:   status,
					BookName: book.Name,
					Imprint:  imprint,
					DueBack:  dueBack,
				})

				break
			}
		}
	}

	return bookInstances, nil
}

// GetAmountBookInstances получает число экземпляров книг, содержащихся в коллекции.
func GetAmountBookInstances() (Total, error) {

	var (
		result Total
		amount int64
		err    error
	)

	libDB := app.ClientDB.Database("local_library")
	ctx := context.TODO()

	amount, err = libDB.Collection("bookinstances").CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, fmt.Errorf("GetAmountBookInstances: %v", err)
	}

	result = Total{DocName: "Copies", Amount: amount}

	return result, nil
}

// GetAmountAvailableBookInstances получает число доступных экземпляров книг, содержащихся в коллекции.
func GetAmountAvailableBookInstances() (Total, error) {

	var (
		result Total
		amount int64
		err    error
	)

	libDB := app.ClientDB.Database("local_library")
	ctx := context.TODO()

	amount, err = libDB.Collection("bookinstances").CountDocuments(ctx, bson.D{primitive.E{Key: "status", Value: "Available"}})
	if err != nil {
		return result, fmt.Errorf("GetAmountAvailableBookInstances: %v", err)
	}

	result = Total{DocName: "Copies available", Amount: amount}

	return result, nil
}

// CreateBookInstance записывает данные нового экземпляра книги в БД.
func CreateBookInstance(imprint, status string, bookID primitive.ObjectID, dueBack interface{}) (string, error) {

	var (
		res            *mongo.InsertOneResult
		bookInstanceID string

		err error
	)

	libDB := app.ClientDB.Database("local_library")
	bookInstancesCollection := libDB.Collection("bookinstances")

	ctx := context.TODO()
	res, err = bookInstancesCollection.InsertOne(ctx, bson.D{
		primitive.E{Key: "status", Value: status},
		primitive.E{Key: "book", Value: bookID},
		primitive.E{Key: "imprint", Value: imprint},
		primitive.E{Key: "due_back", Value: dueBack},
	})
	if err != nil {
		return bookInstanceID, err
	}

	bookInstanceID = res.InsertedID.(primitive.ObjectID).Hex()
	return bookInstanceID, nil
}

// DeleteBookInstanceByID удаляет данные экземпляра книги из БД.
func DeleteBookInstanceByID(bookInstanceID primitive.ObjectID) error {

	var err error

	libDB := app.ClientDB.Database("local_library")
	bookinstancesCollection := libDB.Collection("bookinstances")

	ctx := context.TODO()
	_, err = bookinstancesCollection.DeleteOne(ctx, bson.D{primitive.E{
		Key: "_id", Value: bookInstanceID,
	}})
	if err != nil {
		return err
	}

	return nil
}

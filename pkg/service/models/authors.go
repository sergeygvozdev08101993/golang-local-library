package models

import (
	"context"
	"fmt"
	"reflect"

	"github.com/sergeygvozdev08101993/golang-local-library/internal/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAuthorByID получает данные автора по его ID.
func GetAuthorByID(authorID primitive.ObjectID) (AuthorData, error) {

	var (
		tmpAuthor            bson.M
		author               AuthorData
		dateBirth, dateDeath string

		err error
	)

	libDB := app.ClientDB.Database("local_library")
	authorsCollection := libDB.Collection("authors")

	err = authorsCollection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: authorID}}).Decode(&tmpAuthor)
	if err != nil {
		return author, err
	}
	if err == mongo.ErrNoDocuments {
		return author, err
	}

	if tmpAuthor["date_of_birth"] == nil {
		dateBirth = ""
	} else {
		dateBirth = tmpAuthor["date_of_birth"].(primitive.DateTime).Time().Format("2006-01-02")
	}

	if tmpAuthor["date_of_death"] == nil {
		dateDeath = ""
	} else {
		dateDeath = tmpAuthor["date_of_death"].(primitive.DateTime).Time().Format("2006-01-02")
	}

	author = AuthorData{
		ID:         tmpAuthor["_id"].(primitive.ObjectID).Hex(),
		FirstName:  tmpAuthor["first_name"].(string),
		FamilyName: tmpAuthor["family_name"].(string),
		DateBirth:  dateBirth,
		DateDeath:  dateDeath,
	}

	return author, nil
}

// GetListAllAuthors получает список авторов, отсортированный по фамилии в алфавитном порядке.
func GetListAllAuthors() ([]AuthorData, error) {

	var (
		tmpAuthorList []bson.D
		authorList    []AuthorData

		err error
	)

	libDB := app.ClientDB.Database("local_library")
	authorsCollection := libDB.Collection("authors")

	ctx := context.TODO()
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "family_name", Value: 1}})

	cursor, err := authorsCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &tmpAuthorList); err != nil {
		return nil, err
	}

	for _, author := range tmpAuthorList {

		var (
			dateBirth string
			dateDeath string
		)

		id := author[0]
		firstName := author[1]
		familyName := author[2]

		if author[3].Key == "date_of_birth" {
			dateBirth = getDateStr(author[3].Value)
		}

		if (len(author) == 5 || len(author) == 6) && author[4].Key == "date_of_birth" {
			dateBirth = getDateStr(author[4].Value)
		}

		if (len(author) == 5 || len(author) == 6) && author[4].Key == "date_of_death" {
			dateDeath = getDateStr(author[4].Value)
		}

		authorList = append(authorList, AuthorData{
			ID:         id.Value.(primitive.ObjectID).Hex(),
			FirstName:  firstName.Value.(string),
			FamilyName: familyName.Value.(string),
			DateBirth:  dateBirth,
			DateDeath:  dateDeath,
		})
	}

	return authorList, nil
}

func getDateStr(authorValue interface{}) string {

	var dateStr string

	if reflect.TypeOf(authorValue) != nil {
		dateStr = authorValue.(primitive.DateTime).Time().Format("2 Jan, 2006")
	} else {
		dateStr = ""
	}

	return dateStr
}

// GetAmountAuthors получает число авторов, содержащихся в коллекции.
func GetAmountAuthors() (Total, error) {

	var (
		result Total
		amount int64

		err error
	)

	libDB := app.ClientDB.Database("local_library")
	ctx := context.TODO()

	amount, err = libDB.Collection("authors").CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, fmt.Errorf("GetAmountAuthors: %v", err)
	}

	result = Total{DocName: "Authors", Amount: amount}

	return result, nil
}

// CreateAuthor записывает данные нового автора в БД.
func CreateAuthor(firstName, familyName string, birthDate, deathDate interface{}) (string, error) {

	var (
		res      *mongo.InsertOneResult
		authorID string

		err error
	)

	libDB := app.ClientDB.Database("local_library")
	authorsCollection := libDB.Collection("authors")

	ctx := context.TODO()
	res, err = authorsCollection.InsertOne(ctx, bson.D{
		primitive.E{Key: "first_name", Value: firstName},
		primitive.E{Key: "family_name", Value: familyName},
		primitive.E{Key: "date_of_birth", Value: birthDate},
		primitive.E{Key: "date_of_death", Value: deathDate},
	})
	if err != nil {
		return authorID, err
	}

	authorID = res.InsertedID.(primitive.ObjectID).Hex()
	return authorID, nil
}

// DeleteAuthorByID удаляет данные автора из БД.
func DeleteAuthorByID(authorID primitive.ObjectID) error {

	var err error

	libDB := app.ClientDB.Database("local_library")
	authorsCollection := libDB.Collection("authors")

	ctx := context.TODO()
	_, err = authorsCollection.DeleteOne(ctx, bson.D{primitive.E{
		Key: "_id", Value: authorID,
	}})
	if err != nil {
		return err
	}

	return nil
}

// UpdateAuthorByID обновляет данные автора в БД.
func UpdateAuthorByID(authorID primitive.ObjectID, firstName, familyName string, birthDate, deathDate interface{}) error {

	var err error

	libDB := app.ClientDB.Database("local_library")
	authorsCollection := libDB.Collection("authors")

	ctx := context.TODO()
	_, err = authorsCollection.UpdateOne(
		ctx,
		bson.M{"_id": authorID},
		bson.D{
			{Key: "$set", Value: bson.D{primitive.E{Key: "first_name", Value: firstName}}},
			{Key: "$set", Value: bson.D{primitive.E{Key: "family_name", Value: familyName}}},
			{Key: "$set", Value: bson.D{primitive.E{Key: "date_of_birth", Value: birthDate}}},
			{Key: "$set", Value: bson.D{primitive.E{Key: "date_of_death", Value: deathDate}}},
		})
	if err != nil {
		return err
	}

	return nil
}

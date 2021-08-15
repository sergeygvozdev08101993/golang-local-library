package models

import (
	"context"
	"fmt"

	"github.com/sergeygvozdev08101993/golang-local-library/internal/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetGenreByID получает данные жанра по его ID.
func GetGenreByID(genreID primitive.ObjectID) (GenreData, error) {

	var (
		tmpGenre bson.M
		genre    GenreData
		err      error
	)

	libDB := app.ClientDB.Database("local_library")
	genresCollection := libDB.Collection("genres")

	err = genresCollection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: genreID}}).Decode(&tmpGenre)
	if err != nil {
		return genre, err
	}
	if err == mongo.ErrNoDocuments {
		return genre, err
	}

	genre = GenreData{
		ID:   tmpGenre["_id"].(primitive.ObjectID).Hex(),
		Name: tmpGenre["name"].(string),
	}

	return genre, nil
}

// GetListAllGenres получает список всех жанров, отсортированный по названию в алфавитном порядке.
func GetListAllGenres() ([]GenreData, error) {

	var (
		tmpGenres []bson.D
		genres    []GenreData
		err       error
	)

	libDB := app.ClientDB.Database("local_library")
	genresCollection := libDB.Collection("genres")

	ctx := context.TODO()
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "name", Value: 1}})

	cursor, err := genresCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &tmpGenres); err != nil {
		return nil, err
	}

	for _, genre := range tmpGenres {

		id := genre[0]
		name := genre[1]

		genres = append(genres, GenreData{
			ID:   id.Value.(primitive.ObjectID).Hex(),
			Name: name.Value.(string),
		})
	}

	return genres, nil
}

// GetAmountGenres получает число жанров, содержащихся в коллекции.
func GetAmountGenres() (Total, error) {

	var (
		result Total
		amount int64
		err    error
	)

	libDB := app.ClientDB.Database("local_library")
	ctx := context.TODO()

	amount, err = libDB.Collection("genres").CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, fmt.Errorf("GetAmountGenres: %v", err)
	}

	result = Total{DocName: "Genres", Amount: amount}

	return result, nil
}

// CreateGenre записывает данные нового жанра в БД.
func CreateGenre(name string) (string, error) {

	var (
		res     *mongo.InsertOneResult
		genreID string

		err error
	)

	libDB := app.ClientDB.Database("local_library")
	genresCollection := libDB.Collection("genres")

	ctx := context.TODO()
	res, err = genresCollection.InsertOne(ctx, bson.D{
		primitive.E{Key: "name", Value: name},
	})
	if err != nil {
		return genreID, err
	}

	genreID = res.InsertedID.(primitive.ObjectID).Hex()
	return genreID, nil
}

// DeleteGenreByID удаляет жанр из БД.
func DeleteGenreByID(genreID primitive.ObjectID) error {

	var err error

	libDB := app.ClientDB.Database("local_library")
	genresCollection := libDB.Collection("genres")

	ctx := context.TODO()
	_, err = genresCollection.DeleteOne(ctx, bson.D{primitive.E{
		Key: "_id", Value: genreID,
	}})
	if err != nil {
		return err
	}

	return nil
}

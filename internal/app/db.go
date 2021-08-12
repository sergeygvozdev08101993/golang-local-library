package app

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const dbURI = "mongodb+srv://cooluser:coolpassword@cluster0.a9azn.mongodb.net/local_library?retryWrites=true"

// ClientDB описывает структуру,
// представляющей пул подключений к БД.
var ClientDB *mongo.Client

// ConnectAndPingToClientDB создает нового клиента для подключения его к БД,
// устанавливает это соединение и проверяет его живость.
func ConnectAndPingToClientDB() (*mongo.Client, error) {

	var err error

	client, err := mongo.NewClient(options.Client().ApplyURI(dbURI))
	if err != nil {
		return nil, err
	}

	if err = client.Connect(context.Background()); err != nil {
		return nil, err
	}

	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

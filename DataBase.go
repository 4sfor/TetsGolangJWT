package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

// структура документа в базе
type Doc struct {
	ID      string `bson:"_id"`
	Name    string `bson:"name"`
	Refresh string `bson:"refresh"`
}

// подключение к базе
func ConnectDb() (*mongo.Collection, error) {

	// Устанавливаем контекст и таймаут
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Создаем подключение
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	// Получаем коллекцию
	collection := client.Database("Users").Collection("reg")
	fmt.Println(collection)
	return collection, err
}

// функция проверки наличия пользователя в базе, если он есть происходит запись токена обновления
func CheckGUID(guid string, refreshToken string) bool {
	collection, err := ConnectDb()
	filter := bson.D{{"_id", guid}}
	update := bson.D{
		{"$set", bson.D{
			{"refresh", refreshToken},
		}},
	}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Println(result)
	return true
}

// функция ищет пользователя в базе и проверяет соответсвует ли его токен обновдения переданному
func CheckRefresh(guid string, refreshTokenFromCookie string) bool {
	collection, err := ConnectDb()
	if err != nil {
		log.Fatal(err)
	}
	var doc Doc
	filter := bson.D{{"_id", guid}}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	tokenHash := doc.Refresh
	fmt.Println(tokenHash)
	result := bcrypt.CompareHashAndPassword([]byte(tokenHash), []byte(refreshTokenFromCookie))
	fmt.Println(result)
	return result == nil
}

package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Users struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
	// Другие поля
}

func ConnectDb() {

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
}

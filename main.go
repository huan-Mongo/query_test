package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Doc struct {
	ID bson.RawValue `bson:"_id"`
}

func main() {
	// MongoDB connection URI
	uri := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Get a handle to the users collection
	collection := client.Database("mydb").Collection("test")

	// Insert a doc with regexps string inside {"_id", "/123/"}
	collection.InsertOne(context.Background(), bson.D{{"_id", "/123/"}})

	// The only doc found is {"_id", "/123/"}
	var foundDoc Doc
	collection.FindOne(context.Background(), bson.D{}).Decode(&foundDoc)

	// Insert another doc {"_id", "123"}
	collection.InsertOne(context.Background(), bson.D{{"_id", "123"}})

	cursor, _ := collection.Find(context.Background(), bson.D{{"_id",
		bson.D{{"$in", bson.A{foundDoc.ID}}}}})
	// We shall only find one doc which is {"_id", "/123/"}
	var docs []Doc
	cursor.All(context.Background(), &docs)
	fmt.Println(len(docs))
	for _, doc := range docs {
		fmt.Println(doc.ID.String())
	}
}

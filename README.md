package main

import (
"context"
"fmt"
"log"
"regexp"
"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define a struct to represent your MongoDB document
type User struct {
ID    string `bson:"_id,omitempty"`
Name  string `bson:"name,omitempty"`
Email string `bson:"email,omitempty"`
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
	collection := client.Database("mydb").Collection("users")

	// Regular expression pattern
	pattern := "example.*"

	// Compile the regular expression pattern
	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal(err)
	}

	// Construct the query using the $in operator and the regular expression
	filter := bson.M{"name": bson.M{"$in": bson.A{regex}}}

	// Perform the find operation with the filter
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and print the results
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %s, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}
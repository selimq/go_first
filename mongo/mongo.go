package mon

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/selimq/go1/lib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Baglan s
func Baglan() (*mongo.Collection, context.Context) {
	// Replace the uri string with your MongoDB deployment's connection string.
	uri := "mongodb+srv://selim:selim123@cluster0.d5a5q.mongodb.net/wordstore?retryWrites=true&w=majority"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("wordstore").Collection("words")
	return collection, ctx
}

//CreateWord ..
func CreateWord(ctx context.Context, collection mongo.Collection) {
	print("123s")
	word := lib.Word{ID: 3123}
	ins, err := collection.InsertOne(context.TODO(), word)
	if err != nil {
		log.Fatal(err)
	}
	println(ins.InsertedID)
}

//GetClient is for connect mongodb
func GetClient() *mongo.Client {
	uri := "mongodb+srv://selim:selim123@cluster0.d5a5q.mongodb.net/wordstore?retryWrites=true&w=majority"
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected")
	return client
}

//InsertNewWord is
func InsertNewWord(client *mongo.Client, word lib.Word) (interface{}, error) {
	collection := client.Database("wordstore").Collection("words")
	insertResult, err := collection.InsertOne(context.TODO(), word)
	if err != nil {
		log.Fatalln("Error on inserting new Word", err)
		return nil, err
	}
	log.Println("Word inserted..")
	return insertResult.InsertedID, nil
}

//DeleteWord ..
func DeleteWord(client *mongo.Client, filter bson.M) (interface{}, error) {
	collection := client.Database("wordstore").Collection(("words"))
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatalln("Error on delete Word", err)
		return nil, err
	}
	log.Println("Delete succeed")

	return deleteResult.DeletedCount, nil
}

//ReturnWords ..
func ReturnWords(client *mongo.Client, filter bson.M) []*lib.Word {
	var words []*lib.Word
	collection := client.Database("wordstore").Collection("words")
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal("Error on Finding all the documents", err)
	}
	for cur.Next(context.TODO()) {
		var hero lib.Word
		err = cur.Decode(&hero)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		words = append(words, &hero)
	}
	return words
}

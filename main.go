package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/selimq/go1/lib"
	"github.com/selimq/go1/lib/crud"
	"github.com/selimq/go1/lib/redis"
	mon "github.com/selimq/go1/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//enableCors is for api requests
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
}

func params(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	userID := -1
	var err error
	if val, ok := pathParams["userID"]; ok {
		userID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	commentID := -1
	if val, ok := pathParams["commentID"]; ok {
		commentID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	query := r.URL.Query()
	location := query.Get("location")

	w.Write([]byte(fmt.Sprintf(`{"userID": %d, "commentID": %d, "location": "%s" }`, userID, commentID, location)))
} /*

func examp() {
	resp, err := http.Get("http://authserver.somee.com/api/values")
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	fmt.Print(string(body)[2])
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		panic(err)
	}
	println(body)

}*/

//Persons list of lib:Person
var Persons []lib.Person

func onePerson(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	personID := -1
	var err error
	if val, ok := pathParams["personID"]; ok {
		personID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}
	var personName string = "NoN"
	for i := range Persons {
		if Persons[i].ID == personID {
			personName = Persons[i].Ad
			break
		}
	}
	w.Write([]byte(fmt.Sprintf(`{"personID": %d, "personName": %s }`, personID, personName)))

}
func postPerson(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()
	if err != nil {
		print(err)
	}

	var person lib.Person
	json.Unmarshal(body, &person)
	Persons = append(Persons, person)

	json.NewEncoder(w).Encode(person)
}
func listAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Persons)
}
func deletePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	pathParams := mux.Vars(r)

	personID := -1
	var err error
	if val, ok := pathParams["personID"]; ok {
		personID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a ID"}`))
			return
		}
	}
	for i := range Persons {
		if Persons[i].ID == personID {
			Persons = RemoveIndex(Persons, personID)
			break
		}
	}

	w.Write([]byte("Deleted"))
}

// RemoveIndex for delete of item in slice
func RemoveIndex(s []lib.Person, index int) []lib.Person {
	return append(s[:index], s[index+1:]...)
}
func postWord(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var word lib.Word
	json.Unmarshal(body, &word)
	insertedID, err := mon.InsertNewWord(client, word)
	if err != nil {
		log.Fatal("Error")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	log.Println(insertedID)

}
func deleteWord(w http.ResponseWriter, r *http.Request) {

	pathParams := mux.Vars(r)

	str := pathParams["str"]
	println(str)
	result, err := mon.DeleteWord(client, bson.M{"Text": str})
	if err != nil {
		log.Fatal(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	log.Println(result)

}
func listAllWords(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	log.Println("Listing All Words..")
	words := mon.ReturnWords(client, bson.M{})
	print(words[3].TranslatedText)
	json.NewEncoder(w).Encode(words)

}
func handleRequests() {
	println("localhost:10000 ...")
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("", crud.Get).Methods(http.MethodGet)
	api.HandleFunc("", crud.Post).Methods(http.MethodPost)
	api.HandleFunc("", crud.Put).Methods(http.MethodPut)
	api.HandleFunc("", crud.Delete).Methods(http.MethodDelete)

	//create word
	api.HandleFunc("/word", postWord).Methods((http.MethodPost))
	//delete word
	api.HandleFunc("/word/{str}", deleteWord).Methods((http.MethodDelete))
	//list all
	api.HandleFunc("/word", listAllWords).Methods((http.MethodGet))
	/*//create
	api.HandleFunc("/persons/create", postPerson).Methods((http.MethodPost))
	//list all
	api.HandleFunc("/persons", listAll).Methods((http.MethodGet))
	//one person list
	api.HandleFunc("/persons/{personID}", onePerson).Methods(http.MethodGet)
	//delete
	api.HandleFunc("/persons/delete/{personID}", deletePerson).Methods((http.MethodDelete))

	api.HandleFunc("/user/{userID}/comment/{commentID}", params).Methods(http.MethodGet)*/

	//Heroku vb sitelere yüklendiginde otomatik port seçmesi için
	//eğer yok ise yani lokal ise 10000 olarak atanır
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	log.Println(port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
func init() {
	println("Init state")
}

var client *mongo.Client

func main() {
	//	client = mon.GetClient()
	/*words := mon.ReturnWords(client, bson.M{"Text": "door"})
	for _, word := range words {
		log.Println(word.Text)
	}*/

	ctx := context.Background()

	database, err := redis.NewDatabase("127.0.0.1:6379")
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	pipe := database.Client.TxPipeline()
	pipe.Set(ctx, "key", "value", 0)
	pipe.Set(ctx, "isim", "esse", 0)
	pipe.Exec(ctx)

	deger := pipe.Get(ctx, "esra")
	print(deger)

	//!
	//handleRequests()

	//  collection, context :=  mongo.Baglan()
	/*
		ali := lib.Person{ID: 2, Ad: "Ali"}
		a := lib.Person{ID: 1, Ad: "AA"}
		Persons = append(Persons, ali, a)
		println(Persons)
		//go handleRequests()
		//	v := lib.Word{ID: 54, Text: "droor", TranslatedText: "kapı", CreatedAt: time.Now(), Language: "en_US"}
		//	mongo.CreateWord(&v)
		mongo.CreateWord(context, *collection)*

	*/

}

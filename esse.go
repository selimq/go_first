package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/selimq/go1/lib"
	mon "github.com/selimq/go1/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "get called"}`))
}
func get1(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
	writer.Write(lib.Read()) //read
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))

}

func put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "put called"}`))
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "delete called"}`))
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
	}
	log.Println(insertedID)

}
func handleRequests() {
	println("localhost:8080 ...")
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("", get).Methods(http.MethodGet)
	api.HandleFunc("", post).Methods(http.MethodPost)
	api.HandleFunc("", put).Methods(http.MethodPut)
	api.HandleFunc("", delete).Methods(http.MethodDelete)

	//create words
	api.HandleFunc("/word", postWord).Methods((http.MethodPost))
	/*//create
	api.HandleFunc("/persons/create", postPerson).Methods((http.MethodPost))
	//list all
	api.HandleFunc("/persons", listAll).Methods((http.MethodGet))
	//one person list
	api.HandleFunc("/persons/{personID}", onePerson).Methods(http.MethodGet)
	//delete
	api.HandleFunc("/persons/delete/{personID}", deletePerson).Methods((http.MethodDelete))

	api.HandleFunc("/user/{userID}/comment/{commentID}", params).Methods(http.MethodGet)*/

	log.Fatal(http.ListenAndServe(":8080", r))

}
func init() {
	println("İnit")
}

var client *mongo.Client

func main() {
	client = mon.GetClient()
	words := mon.ReturnWords(client, bson.M{"Text": "door"})
	for _, word := range words {
		log.Println(word.Text)
	}
	handleRequests()
	/* collection, context :=  mongo.Baglan()

	ali := lib.Person{ID: 2, Ad: "Ali"}
	a := lib.Person{ID: 1, Ad: "AA"}
	Persons = append(Persons, ali, a)
	println(Persons)
	//go handleRequests()
	//	v := lib.Word{ID: 54, Text: "droor", TranslatedText: "kapı", CreatedAt: time.Now(), Language: "en_US"}
	//	mongo.CreateWord(&v)
	mongo.CreateWord(context, *collection)*/
}

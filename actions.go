package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MutantsTransactions_cll = MeliContextDB().Collection("MutantsTransactions")

var collection = getSession().DB("MELI").C("MutantsTransactions")

const uri = "mongodb+srv://root:admin@atlascluster.fewo3js.mongodb.net/?retryWrites=true&w=majority"

func MeliContextDB() *mongo.Database {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database("MELI")
}

func getSession() *mgo.Session {
	session, err := mgo.Dial("mongodb://localhost:27017")

	if err != nil {
		panic(err)
	}
	return session
}

func Response(w http.ResponseWriter, status int, result Movie) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(result)
}

func ResponseList(w http.ResponseWriter, status int, result Movies) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(result)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola mundo desde mi servidor GO")
}

func Contacto(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Esta es la pagina de contacto")
}

func MovieList(w http.ResponseWriter, r *http.Request) {
	var result []Movie
	err := collection.Find(nil).Sort("-_id").All(&result)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Results :", result)
	}
	ResponseList(w, 200, result)
}
func MovieShow(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	movieId := params["id"]

	if !bson.IsObjectIdHex(movieId) {
		w.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(movieId)
	result := Movie{}
	err := collection.FindId(oid).One(&result)

	if err != nil {
		w.WriteHeader(404)
		return
	} else {
		Response(w, 200, result)
	}
}

func MovieAdd(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var movie_data Movie

	err := decoder.Decode(&movie_data)

	if err != nil {
		panic(err)
	}

	defer r.Body.Close()
	result, err := MutantsTransactions_cll.InsertOne(context.TODO(), movie_data)

	fmt.Println(result, err)

	Response(w, 200, movie_data)
}

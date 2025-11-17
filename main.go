package main

import (
	"context"
	"bufio"
	"strings"
	"os"
	"time"
	"net/http"
	"encoding/json"
	"encoding/base64"
	"crypto/rand"
	"log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)
var coll *mongo.Collection
var max_TTL time.Duration = 48

type URLMapping struct{
	Id bson.ObjectID `bson:"_id"`
	Link string `bson:"link"`
	Shorted string `bson:"shorted"`	
	Expires_at time.Time `bson:"expiresAt, omitempty"`
	Created_at time.Time `bson:"createAt"`
}


func get_uri() string{
	var uri string;
	file, err := os.Open(".env")
	if err != nil{
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	uri = strings.SplitN(line, "=", 2)[1]
	uri = uri[1:len(uri) - 1]
	return uri
}


func initMongo() {
	clientOps := options.Client().ApplyURI(get_uri())
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	client, err := mongo.Connect(clientOps) 
	if err != nil{
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil{
		log.Fatal(err)
	}
	log.Println("Connexted to Mongo. yay!")
	coll = client.Database("shortener").Collection("links")
	index := mongo.IndexModel{
		Keys : bson.M{"expiresAt": 1},
		Options : options.Index().SetExpireAfterSeconds(0),
	}
	if _, err := coll.Indexes().CreateOne(ctx, index); err != nil{
		log.Fatal("Error in TTL Index creation: ", err)
	}
}


func main(){
	initMongo()
	router := http.NewServeMux()
	router.Handle("/", http.FileServer(http.Dir("./static")))
	router.HandleFunc("POST /shorten", shorten)
	router.HandleFunc("GET /{shortened}", redirecter)
	http.ListenAndServe(":8080", router)
}

//Function to redirect the shortened links with the original ones
func redirecter (w http.ResponseWriter, r *http.Request){
	short := r.PathValue("shortened")
	log.Printf("The path value is ", short)
	var doc URLMapping
	err := coll.FindOne(r.Context(), bson.M{"shorted": short}).Decode(&doc)
	if err != nil{
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, doc.Link, http.StatusFound)
}

func shorten(w http.ResponseWriter, r *http.Request){
	var data struct{
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil{
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		log.Printf("The error is : ", err)
		return
	}
	defer r.Body.Close()
	log.Printf("Recieved: Link:%s", data.URL)
	if valid_url := check("link", data.URL, true); valid_url != ""{
		resp := map[string]string{
			"shorted": valid_url,
		}
		jsonResponse(w, resp, http.StatusCreated)
		return
	}
	code := code_gen()
	log.Printf("The value of code is: ", code)
	doc := URLMapping{
		Id : bson.NewObjectID(),
		Link: data.URL,
		Shorted: code,
		Expires_at: time.Now().Add(time.Hour * max_TTL),
		Created_at : time.Now(),
	}
	jsonResponse(w, map[string]string{"shorted_value" : code}, http.StatusCreated)
	if _, err := coll.InsertOne(r.Context(), doc); err != nil{
		http.Error(w, "Insertion Failed", http.StatusInternalServerError)
		log.Printf("Error insertion: ", err)
		return
	}
}


func jsonResponse(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, `{"error":"failed to encode json"}`, http.StatusInternalServerError)
    }
}


//Function to check the URL is already shorted
func check(key string, value string, is_update bool) string {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()
	
	filter := bson.M{ key : value}

    // FindOneAndUpdate returns SingleResult
    res := coll.FindOne(ctx, filter)

    var doc URLMapping
    err := res.Decode(&doc)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            // document not found, handle accordingly
            log.Println("No document found for key:", key)
            return ""
        } else {
            log.Fatal(err)
        }
    }    
	if is_update {
	update := bson.M{"$set": bson.M{"expiresAt": time.Now().Add(max_TTL * time.Hour)}}
	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil{
		log.Fatal(err)
		}
	}
    return doc.Shorted
}


func code_gen() string {
	var n int = 6
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil{
		return ""
	}
	code := base64.RawURLEncoding.EncodeToString(b)
	if check("shorted", code, false) != ""{
		return code_gen()
	}
	return code
}
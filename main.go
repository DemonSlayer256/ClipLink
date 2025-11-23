package main

import (
	"ClipLink/middleware"
	m "ClipLink/models"
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var (
	user_coll     *mongo.Collection
	link_coll     *mongo.Collection
	max_TTL       time.Duration = 48
	max_url_limit int           = 5
	signingKey                  = []byte(LoadEnv("SECURE_KEY"))
)

func hasher(pass string) string {
	hash := sha256.Sum256([]byte(pass))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func register(w http.ResponseWriter, r *http.Request) {
    var newUser m.User

    // Decode JSON safely
    if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
        jsonResponse(w, map[string]string{"message": "Invalid JSON data"}, http.StatusBadRequest)
        log.Println("Error decoding JSON:", err)
        return
    }
    defer r.Body.Close()

    // Validate required fields
    if strings.TrimSpace(newUser.User) == "" || strings.TrimSpace(newUser.Pass) == "" {
        jsonResponse(w, map[string]string{"message": "Username and password are required"}, http.StatusBadRequest)
        return
    }

    // Ensure the user does not already exist
    var existingUser m.User
    err := user_coll.FindOne(r.Context(), bson.M{"user": newUser.User}).Decode(&existingUser)

    if err != nil && err != mongo.ErrNoDocuments {
        jsonResponse(w, map[string]string{"message": "Error checking for existing user"}, http.StatusInternalServerError)
        log.Println("Error checking user:", err)
        return
    }

    if err == nil {
        // user already exists
        jsonResponse(w, map[string]string{"message": "User already exists"}, http.StatusConflict)
        return
    }

    // Now safe to create the new user
    newUser.Id = primitive.NewObjectID()
    newUser.Pass = hasher(newUser.Pass) // hash password
    newUser.Left = max_url_limit

    _, err = user_coll.InsertOne(r.Context(), newUser)
    if err != nil {
        jsonResponse(w, map[string]string{"message": "Error inserting user"}, http.StatusInternalServerError)
        log.Println("Error inserting user:", err)
        return
    }

    jsonResponse(w, map[string]string{"message": "User registered successfully"}, http.StatusCreated)
}


func LoadEnv(keys ...string) string {
	env := make(map[string]string)

	file, err := os.Open(".env")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}
		// Split line into key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip invalid lines
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) || (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}
		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return ""
	}

	// Filter only the requested keys
	var required string
	for _, key := range keys {
		if value, exists := env[key]; exists {
			required = value
			break
		}
	}

	return required
}
func initMongo() {
	clientOps := options.Client().ApplyURI(LoadEnv("MONGO_URI"))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(clientOps)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connexted to Mongo. yay!")
	link_coll = client.Database("shortener").Collection("links")
	if err := link_coll.Drop(ctx); err != nil {
		log.Fatal(err)
	}
	user_coll = client.Database("shortener").Collection("users")
	if err := user_coll.Drop(ctx); err != nil {
		log.Fatal(err)
	}
	index := mongo.IndexModel{
		Keys:    bson.M{"expiresAt": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	if _, err := link_coll.Indexes().CreateOne(ctx, index); err != nil {
		log.Fatal("Error in TTL Index creation: ", err)
	}
}
func generateToken(user string) (string, error) {
	claims := m.JWT{
		Username: user,
		Exp:      time.Now().Add(1 * time.Hour).Unix(),
		Iat:      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func login(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Username string `json:"user"`
		Password string `json:"pass"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		jsonResponse(w, map[string]string{"message" : "Invalid JSON data"}, http.StatusBadRequest)
		log.Println("Error decoding json: ", err)
		return
	}
	defer r.Body.Close()
	var stored m.User
	err := user_coll.FindOne(r.Context(), bson.M{"user": data.Username}).Decode(&stored)
	if err != nil {
		jsonResponse(w, map[string]string{"message" : "Invalid credentials"}, http.StatusUnauthorized)
		return
	}
	if hasher(data.Password) != stored.Pass {
		jsonResponse(w, map[string]string{"message" : "Invalid credentials"}, http.StatusUnauthorized)
		return
	}

	token, err := generateToken(data.Username)
	if err != nil {
		jsonResponse(w, map[string]string{"message" : "Token generation error"}, http.StatusInternalServerError)
		log.Println("Error in token generation: ", err)
		return
	}

	jsonResponse(w, map[string]string{"token": token}, http.StatusOK)
}

func main() {
	initMongo()
	router := http.NewServeMux()
	router.HandleFunc("/", renderer)
	// The below shorten can be changed to use GET and POST methods for a different route named Delete_url. However
	// Since this is an API and not a webpage, DELETE is used but it may be subject to change in future
	router.Handle("GET /links", middleware.Auth(http.HandlerFunc(send_all)))
	router.Handle("POST /delete", middleware.Auth(http.HandlerFunc(delete_link)))
	router.Handle("POST /shorten", middleware.Auth(http.HandlerFunc(shorten)))
	router.HandleFunc("POST /register", register)
	router.HandleFunc("POST /login", login)
	router.Handle("GET /user/{route}", http.HandlerFunc(renderer))
	router.HandleFunc("GET /{shortened}", redirecter)
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	server.ListenAndServe()
}

func send_all(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    // Retrieve the username from the context
    user, ok := r.Context().Value("username").(string)
    if !ok {
        jsonResponse(w, map[string]string{"Error": "Something happened"}, http.StatusInternalServerError)
        return
    }

    // Create a filter to find all links belonging to the user
    filter := bson.M{"user": user}

    // Find all links for the user
    cursor, err := link_coll.Find(ctx, filter)
    if err != nil {
        jsonResponse(w, map[string]string{"Error": "Error fetching links"}, http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx) // Ensure the cursor is closed when done

    // Prepare a slice to hold the links
    var links []m.URLMapping

    // Iterate through the cursor and decode each link into the slice
    for cursor.Next(ctx) {
        var link m.URLMapping
        if err := cursor.Decode(&link); err != nil {
            jsonResponse(w, map[string]string{"Error": "Error decoding links"}, http.StatusInternalServerError)
            return
        }
		link.Link =  "https://" + link.Link
        links = append(links, link)
    }

    // Check for cursor errors
    if err := cursor.Err(); err != nil {
        jsonResponse(w, map[string]string{"Error": "Error during cursor iteration"}, http.StatusInternalServerError)
        return
    }

    // Return the list of links as a JSON response
    jsonResponse(w, links, http.StatusOK)
}


func renderer(w http.ResponseWriter, r *http.Request) {
		path := r.PathValue("route")
		log.Println(r.URL.Path)
		if path == "" || path == "/" {
			fmt.Printf("The path is : %s", r.URL.EscapedPath())
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}
		if strings.Contains(path, ".."){
			jsonResponse(w, map[string]string{"message" : "Forbidden path"}, http.StatusForbidden)
			return
		}
		if !strings.HasSuffix(path, ".html") {
			path = path + ".html"
		}
		http.ServeFile(w, r, "./static/" + path)
}

// Function to redirect the shortened links with the original ones
func redirecter(w http.ResponseWriter, r *http.Request) {
	short := r.PathValue("shortened")
	var doc m.URLMapping
	err := link_coll.FindOne(r.Context(), bson.M{"shorted": short}).Decode(&doc)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	//Since the normalized link does not contain scheme, we default it to http
	var link string = "https://" + doc.Link
	http.Redirect(w, r, link, http.StatusFound)
}

func delete_link(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	username, ok := r.Context().Value("username").(string)
	if !ok {
		jsonResponse(w, map[string]string{"Error": "Something happened"}, http.StatusInternalServerError)
	}
	var data struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		jsonResponse(w, map[string]string{"message" : "Invalid JSON data"}, http.StatusBadRequest)
		log.Println("The Error is : ", err, r.Body)
		return
	}
	defer r.Body.Close()
	if data.Code == "" {
		jsonResponse(w, map[string]string{"message" : "Code should be present"}, http.StatusBadRequest)
		return
	}
	if _, err := link_coll.DeleteOne(ctx, bson.M{"shorted": data.Code, "user": username}); err != nil {
		jsonResponse(w, map[string]string{"message" : "The link is not found"}, http.StatusBadRequest)
		log.Println("An error in deletion: ", err)
		return
	}
	update := bson.M{"$inc": bson.M{"left": 1}}
	if _, err := user_coll.UpdateOne(ctx, bson.M{"user": username}, update); err != nil {
		jsonResponse(w, map[string]string{"message" : "Error updating remaining urls"}, http.StatusInternalServerError)
		log.Println("Error updating remaining: ", err)
		return
	}
	jsonResponse(w, map[string]string{"Message": "Successfully deleted"}, http.StatusFound)
}

func normalize(link string) string {
	parsed, err := url.Parse(link)
	if err != nil {
		log.Println("Error in parsing url : ", err)
		return ""
	}
	normalized := parsed.Host + parsed.Path
	if parsed.RawQuery != "" {
		normalized += "?" + parsed.RawQuery
	}
	if parsed.Fragment != "" {
		normalized += "#" + parsed.Fragment
	}
	return normalized
}

func within_limit(user string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var data m.User
	err := user_coll.FindOne(ctx, bson.M{"user": user}).Decode(&data)
	if err != nil {
		log.Println("Error finding user: ", err)
		return false
	}
	if data.User == user && data.Left > 0 {
		update := bson.M{"$inc": bson.M{"left": -1}}
		if _, err := user_coll.UpdateOne(ctx, bson.M{"user": user}, update); err != nil {
			log.Println("Error updating remaining: ", err)
			return false
		}
		return true
	}
	return false
}

func shorten(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		jsonResponse(w, map[string]string{"message" : "Invalid JSON data"}, http.StatusBadRequest)
		log.Println("The Error is : ", err)
		return
	}
	username, ok := r.Context().Value("username").(string)
	if !ok {
		jsonResponse(w, map[string]string{"Error": "Something happened"}, http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	log.Printf("Recieved: Link: %s", data.URL)
	normalized := normalize(data.URL)
	if normalized == "" {
		jsonResponse(w, map[string]string{"Error": "Invalid URL. Expected {\"url\": \"https://example.com\"} but recieved something else"}, http.StatusBadRequest)
		return
	}
	if valid_url := check("link", normalized, true, username); valid_url != "" {
		resp := map[string]string{
			"shorted_value": valid_url,
		}
		jsonResponse(w, resp, http.StatusCreated)
		return
	}
	if !within_limit(username) {
		jsonResponse(w, map[string]string{"Error": "You cannot have more than 5 active links."}, http.StatusBadRequest)
		return
	}
	code := code_gen()
	log.Printf("The value of code is: %s", code)
	doc := m.URLMapping{
		Id:         primitive.NewObjectID(),
		Link:       normalized,
		Shorted:    code,
		Expires_at: time.Now().Add(time.Hour * max_TTL),
		Created_at: time.Now(),
		User:       username,
	}

	if _, err := link_coll.InsertOne(r.Context(), doc); err != nil {
		jsonResponse(w, map[string]string{"message" : "Insertion Failed"}, http.StatusInternalServerError)
		log.Println("Error insertion: ", err)
		return
	}
	jsonResponse(w, map[string]string{"shorted_value": code}, http.StatusCreated)
}

func jsonResponse(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")

    b, err := json.Marshal(data)
    if err != nil {
        log.Println("Failed to encode JSON response:", err)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`{"error":"failed to encode json"}`))
        return
    }

    w.WriteHeader(status)
    w.Write(b)
}

// Function to check the URL is already shorted
func check(key string, value string, is_update bool, user string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{key: value}

	// FindOneAndUpdate returns SingleResult
	res := link_coll.FindOne(ctx, filter)

	var doc m.URLMapping
	err := res.Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// document not found, handle accordingly
			return ""
		} else {
			log.Fatal(err)
		}
	}
	if is_update && doc.User == user {
		update := bson.M{"$set": bson.M{"expiresAt": time.Now().Add(max_TTL * time.Hour)}}
		_, err := link_coll.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
		}
		return doc.Shorted
	}
	return ""
}

func code_gen() string {
	var n int = 6
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	code := base64.RawURLEncoding.EncodeToString(b)
	if check("shorted", code, false, "") != "" {
		return code_gen()
	}
	return code
}

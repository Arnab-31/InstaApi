package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"context"
	"net/url"
	"strconv"

	
	"net/http"
	"strings"
	"sync"
	"time"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"


	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)


//struct for Post and Users
type Post struct {
	Caption          string `json:"caption"`
	ImageURL         string `json:"imageUrl"`
	ID               string `bson:"id"`
	UserID           string `bson:"userID`
	Timestamp        time.Time `json:"timestamp"`
}

type User struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	ID               string `bson:"id"`
	Password         string `json:"password`
}


//Collections for posts and users
var collection *mongo.Collection
var userCollection *mongo.Collection
var ctx = context.TODO()



//Function for Connecting to Database
func dbConnect(){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("Enter MONGO URL"))

	
	if err != nil {
		fmt.Println("Error connecting to database")
	}else{
		fmt.Println("Database Connected!")
	}

	collection = client.Database("InstaApi").Collection("Posts")
	userCollection = client.Database("InstaApi").Collection("Users")
}


type coasterHandlers struct {
	sync.Mutex
	store map[string]Post
}

//function for redirecting get and route request to their respective functions
func (h *coasterHandlers) posts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

//function to get all posts using pagination. Each page gives 5 posts. /posts?page=1
func (h *coasterHandlers) get(w http.ResponseWriter, r *http.Request) {

	filter := bson.D{{}}
	u, err := url.Parse(r.URL.String())
	q := u.Query()
	fmt.Println(q["page"][0])

	cur, err := collection.Find(ctx, filter)

	if err != nil {
        fmt.Println("Error")
    }

	var posts []*Post
	for cur.Next(ctx) {
        var t Post
        err := cur.Decode(&t)
        if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
        }
	
        posts = append(posts, &t)
    }

	i, _ := strconv.Atoi(q["page"][0])

	var low = (i-1) * 5
	var high = low + 5


	if low > len(posts){
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}else if high > len(posts){
		high = len(posts)
	}

	posts = posts[low:high]

	jsonBytes2, err := json.Marshal(posts)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes2)

	fmt.Println(cur)
	
}


//function to get a specific post using its id. /posts/:id
func (h *coasterHandlers) getPost(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	filter := bson.D{{"id", parts[2]}}

	var Result *Post
	err := collection.FindOne(ctx, filter).Decode(&Result)
	if err != nil {
        fmt.Println("Error")
    }

	jsonBytes2, _ := json.Marshal(Result)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes2)

	fmt.Println(Result)
}

// function to get all posts having a specific user id.  /posts/users/:id
func (h *coasterHandlers) getUserPosts(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	filter := bson.D{{}}
	cur, _ := collection.Find(ctx, filter)
	var posts []*Post
	fmt.Println(parts[3])

	for cur.Next(ctx) {
        var t Post
        err := cur.Decode(&t)
        if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			fmt.Println("Error")

        }
		if(t.UserID == parts[3]){
			posts = append(posts, &t)
		}
       
    }

	jsonBytes2,_ := json.Marshal(posts)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes2)
	fmt.Println(cur)
}


//function to upload a post.  /posts
func (h *coasterHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var coaster Post
	err = json.Unmarshal(bodyBytes, &coaster)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	coaster.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[coaster.ID] = coaster
	doc, err := collection.InsertOne(ctx, coaster)
	fmt.Println(doc) 
	jsonBytes2,_ := json.Marshal(doc)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes2)
	defer h.Unlock()
}


func newCoasterHandlers() *coasterHandlers {
	return &coasterHandlers{
		store: map[string]Post{},
	}

}


//function for ecnrypting passwords
func encrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}



//function for creating new users. /users
func (h *coasterHandlers) createUser(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var coaster User
	err = json.Unmarshal(bodyBytes, &coaster)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}


	coaster.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()

	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
	var a = encrypt(coaster.Password, key)
	coaster.Password = string(a)
	doc, err := userCollection.InsertOne(ctx, coaster)
	fmt.Println(doc) 

	jsonBytes2,_ := json.Marshal(coaster.ID)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes2)
	
	defer h.Unlock()
}


//function to get a speciifc user using its id.  /users/:id
func (h *coasterHandlers) getUser(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	filter := bson.D{{"id", parts[2]}}

	var Result *User
	err := userCollection.FindOne(ctx, filter).Decode(&Result)
	if err != nil {
        fmt.Println(err)
    }

	jsonBytes2, _ := json.Marshal(Result)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes2)

	fmt.Println(Result)

}





func main() {
	
	dbConnect()
	coasterHandlers := newCoasterHandlers()
	http.HandleFunc("/posts", coasterHandlers.posts)
	http.HandleFunc("/posts/", coasterHandlers.getPost)
	http.HandleFunc("/posts/users/", coasterHandlers.getUserPosts)
	http.HandleFunc("/users", coasterHandlers.createUser)
	http.HandleFunc("/users/", coasterHandlers.getUser)

	
	err := http.ListenAndServe(":8080", nil)    //listening to port: 8080
	if err != nil {
		panic(err)
	}
}
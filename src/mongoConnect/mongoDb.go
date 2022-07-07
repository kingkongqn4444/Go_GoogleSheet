package mongoConnect

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var userCollection *mongo.Collection
var categoryCollection *mongo.Collection

type Person struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Age         int                `json:"age,omitempty" bson:"age,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Phone 		string 			   `json:"phone",omitempty bson:"phone:omitempty"`
}

type User struct {
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"password"`
	Token     string `json:"token"`
}

type Category struct {
	ID  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Category_Name string `json:"categoryName"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

func ConnectMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clienOptions := options.Client().ApplyURI("mongodb://localhost:27000")
	client, _ := mongo.Connect(ctx, clienOptions)
	database := client.Database("BE_Wallet")


	userCollection = database.Collection("user")
	categoryCollection = database.Collection("categories")
	fmt.Println("Connect Success !!")
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("file", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func GetPersonByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	id, _ := primitive.ObjectIDFromHex(mux.Vars(request)["id"])
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := userCollection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(person)
}

func UpdatePerson(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	id, _ := primitive.ObjectIDFromHex(mux.Vars(request)["id"])
	json.NewDecoder(request.Body).Decode(&person)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := userCollection.UpdateOne(ctx, Person{ID: id}, bson.M{"$set": person})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}

func DeletePerson(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	id, _ := primitive.ObjectIDFromHex(mux.Vars(request)["id"])
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := userCollection.DeleteOne(ctx, Person{ID: id})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}


// Tạo tài khoản user
func CreatePerson(c *gin.Context){
	c.Header("content-type", "application/json")
	var person Person
	json.NewDecoder(c.Request.Body).Decode(&person)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _ := userCollection.InsertOne(ctx, person)
	id := result.InsertedID
	person.ID = id.(primitive.ObjectID)
	json.NewEncoder(c.Writer).Encode(person)
}

// Đăng nhập tài khoản
func LoginHandler(c *gin.Context) {
	c.Header("content-type", "application/json")
	var user User
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	var result User
	var res ResponseResult
	err = userCollection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
	if err != nil {
		res.Error = "Invalid username"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		res.Error = "Invalid password"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  result.Username,
		"firstname": result.FirstName,
		"lastname":  result.LastName,
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		res.Error = "Error while generating token,Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	result.Token = tokenString
	result.Password = ""
	json.NewEncoder(c.Writer).Encode(result)
}

// get profile current user
func ProfileHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	var result User
	var res ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Username = claims["username"].(string)
		result.FirstName = claims["firstname"].(string)
		result.LastName = claims["lastname"].(string)

		json.NewEncoder(c.Writer).Encode(result)
		return
	} else {
		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
}

// Lấy danh sách categories
func GetListAllCategory(c *gin.Context) {
	c.Header("Content-Type","application/json")
	var listCategories []Category
	cur, _ := categoryCollection.Find(context.TODO(), bson.D{{}})
	for cur.Next(context.TODO()) {
		var category Category
		cur.Decode(&category)
		listCategories = append(listCategories, category)
	}
	cur.Close(context.TODO())
	json.NewEncoder(c.Writer).Encode(listCategories)
}

func RegisterHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var user User
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &user)
	var res ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	var result User
	err = userCollection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
			if err != nil {
				res.Error = "Error While Hashing Password, Try Again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}
			user.Password = string(hash)
			_, err = userCollection.InsertOne(context.TODO(), user)
			if err != nil {
				res.Error = "Error While Creating User, Try Again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}
			res.Result = "Registration Successful"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}
		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Result = "Username already Exists!!"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

// Tạo mới categories
func CreateCategory(c *gin.Context)  {
	c.Header("Content-Type","application/json")
	var category Category
	var res ResponseResult
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &category)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = categoryCollection.FindOne(context.TODO(), bson.D{{"category_name",category.Category_Name}}).Decode(&category)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			_, err := categoryCollection.InsertOne(ctx, category)
			if err != nil {
				res.Error = "Error While Creating User, Try Again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}
			res.Result = "Create category successful !!"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}
		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Result = "Category already Exists!!"
	json.NewEncoder(c.Writer).Encode(res)
}

func Login(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	id, _ := primitive.ObjectIDFromHex(mux.Vars(request)["id"])
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := userCollection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(person)
}
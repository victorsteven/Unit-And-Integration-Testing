package tests

import (
	"bytes"
	"database/sql"
	"efficient-api/controllers"
	"efficient-api/domain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

const (
	queryGetMessage     = "SELECT id, title, body, created_at FROM messages WHERE id=?;"
	queryInsertMessage  = "INSERT INTO messages(title, body, created_at) VALUES(?, ?, ?);"
	queryUpdateMessage  = "UPDATE messages SET title=?, body=? WHERE id=?;"
	queryDeleteMessage  = "DELETE FROM messages WHERE id=?;"
	queryGetAllMessages = "SELECT id, title, body, created_at FROM messages;"
)

var (
	tm     = time.Now()
	server = domain.Server{}
	dbNow  *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("./../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}


	os.Exit(m.Run())
}

// func dbConn()   {
// 	var err error
// 	dbdriver := os.Getenv("DBDRIVER_TEST")
// 	username := os.Getenv("USERNAME_TEST")
// 	password := os.Getenv("PASSWORD_TEST")
// 	host := os.Getenv("HOST_TEST")
// 	database := os.Getenv("DATABASE_TEST")
// 	port := os.Getenv("PORT_TEST")

// 	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database)

// 	server.DB, err = sql.Open(dbdriver, DBURL)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// }

func Database() {
	dbDriver := os.Getenv("DBDRIVER_TEST")
	username := os.Getenv("USERNAME_TEST")
	password := os.Getenv("PASSWORD_TEST")
	host := os.Getenv("HOST_TEST")
	database := os.Getenv("DATABASE_TEST")
	port := os.Getenv("PORT_TEST")

	dbNow = domain.MessageRepo.Initialize(dbDriver, username, password, port, host, database)
}

func refreshUserTable() error {

	stmt, err := dbNow.Prepare("TRUNCATE TABLE messages;")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("Error truncating messages table: %s", err)
	}
	return nil
}

//func TestAde(t *testing.T) {
//	refreshUserTable()
//	//seedOneMessage()
//	//fmt.Println("this is the seed: ", msg)
//}

func seedOneMessage() (domain.Message, error) {
	msg := domain.Message{
		Title:     "the title",
		Body:      "the body",
		CreatedAt: time.Now(),
	}
	stmt, err := dbNow.Prepare(queryInsertMessage)
	if err != nil {
		panic(err.Error())
	}
	insertResult, createErr := stmt.Exec(msg.Title, msg.Body, msg.CreatedAt)
	if createErr != nil {
		log.Fatalf("Error creating message: %s", createErr)
	}
	msgId, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("Error creating message: %s", createErr)
	}
	msg.Id = msgId
	return msg, nil
}

//func TestHello(t *testing.T) {
//	message, err := seedOneMessage()
//	if err != nil {
//		fmt.Println("error creating message: ", err)
//	}
//	fmt.Println("this is the message created: ", message)
//}

//func seedMessages() ([]domain.Message, error) {
//
//	var err error
//	if err != nil {
//		return nil, err
//	}
//	messages := []domain.Message{
//		domain.Message{
//			Title:     "first title",
//			Body:      "first the body",
//			CreatedAt: time.Now(),
//		},
//		domain.Message{
//			Title:     "second title",
//			Body:      "second body",
//			CreatedAt: time.Now(),
//		},
//	}
//
//	for i, _ := range messages {
//		//err := server.DB.Model(&models.User{}).Create(&users[i]).Error
//		//if err != nil {
//		//	return []models.User{}, err
//		//}
//	}
//	return messages, nil
//}

//func Add(x, y int64) int64 {
//	//fmt.Println("adding")
//	return x + y
//}

//func TestAdd(t *testing.T) {
//	ans := Add(2, 4)
//
//	fmt.Println("this is the ans", ans)
//}

func TestCreateMessage(t *testing.T) {

	Database()

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON  string
		statusCode int
		title      string
		body       string
		errMessage string
	}{
		{
			inputJSON:  `{"title":"the title", "body": "the body"}`,
			statusCode: 201,
			title:      "the title",
			body:       "the body",
			errMessage: "",
		},
		{
			inputJSON:  `{"title":"the title", "body": "the body"}`,
			statusCode: 500,
			errMessage: "title already taken",
		},
		{
			inputJSON:  `{"title":"", "body": "the body"}`,
			statusCode: 422,
			errMessage: "Please enter a valid title",
		},
		{
			inputJSON:  `{"title":"the title", "body": ""}`,
			statusCode: 422,
			errMessage: "Please enter a valid body",
		},
	}

	for _, v := range samples {
		r := gin.Default()
		r.POST("/messages", controllers.CreateMessage)
		req, err := http.NewRequest(http.MethodPost, "/messages", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		fmt.Println("this is the response data: ", responseMap)
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			//casting the interface to map:
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["body"], v.body)
		}
		if v.statusCode == 400 || v.statusCode == 422 || v.statusCode == 500 && v.errMessage != "" {
			assert.Equal(t, responseMap["message"], v.errMessage)
		}
	}
}

func TestGetMessageByID(t *testing.T) {

	Database()

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	message, err := seedOneMessage()
	if err != nil {
		t.Errorf("Error while seeding table: %s", err)
	}
	samples := []struct {
		id         string
		statusCode int
		title      string
		body       string
		errMessage string
	}{
		{
			id:         strconv.Itoa(int(message.Id)),
			statusCode: 200,
			title:   message.Title,
			body:     message.Body,
			errMessage: "",
		},
		{
			id:         "unknwon",
			statusCode: 400,
			errMessage: "message id should be a number",
		},
		{
			id:         strconv.Itoa(12322), //an id that does not exist
			statusCode: 404,
			errMessage: "no record matching given id",
		},
	}
	for _, v := range samples {
		r := gin.Default()
		r.GET("/messages/:message_id", controllers.GetMessage)
		req, err := http.NewRequest(http.MethodGet, "/messages/"+v.id, nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		//fmt.Println("this is the response data: ", responseMap)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			//casting the interface to map:
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["body"], v.body)
		}
		if v.statusCode == 400 || v.statusCode == 422 && v.errMessage != "" {
			assert.Equal(t, responseMap["message"], v.errMessage)
		}
	}
}

func TestUpdateMessage(t *testing.T) {

	Database()

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	message, err := seedOneMessage()
	fmt.Println("this is the error: ", message)
	if err != nil {
		t.Errorf("Error while seeding table: %s", err)
	}

	samples := []struct {
		id          string
		inputJSON  string
		statusCode int
		title      string
		body       string
		errMessage string
	}{
		{
			id:          strconv.Itoa(int(message.Id)),
			inputJSON:  `{"title":"update title", "body": "update body"}`,
			statusCode: 200,
			title:      "update title",
			body:       "update body",
			errMessage: "",
		},
		{
			//Empty title
			id:          strconv.Itoa(int(message.Id)),
			inputJSON:  `{"title":"", "body": "update body"}`,
			statusCode: 422,
			errMessage: "Please enter a valid title",
		},
		{
			//Empty title
			id:          strconv.Itoa(int(message.Id)),
			inputJSON:  `{"title":"the title", "body": ""}`,
			statusCode: 422,
			errMessage: "Please enter a valid body",
		},
		{
			id:         "unknwon",
			statusCode: 400,
			errMessage: "message id should be a number",
		},
		{
			id:         strconv.Itoa(12322), //an id that does not exist
			inputJSON:  `{"title":"the title", "body": "the body"}`,
			statusCode: 404,
			errMessage: "no record matching given id",
		},
	}

	for _, v := range samples {
		r := gin.Default()
		r.PUT("/messages/:message_id", controllers.UpdateMessage)
		req, err := http.NewRequest(http.MethodPut, "/messages/"+v.id, bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		fmt.Println("this is the response data: ", responseMap)
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			//casting the interface to map:
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["body"], v.body)
		}
		if v.statusCode == 400 || v.statusCode == 422 || v.statusCode == 500 && v.errMessage != "" {
			assert.Equal(t, responseMap["message"], v.errMessage)
		}
	}
}

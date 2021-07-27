package integration__tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SemmiDev/lets-tests/controllers"
	"github.com/SemmiDev/lets-tests/domain"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCreateChat(t *testing.T) {
	database()
	err := refreshChatsTable()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		inputJSON  string
		statusCode int
		sender     string
		receiver   string
		body       string
		errMessage string
	}{
		{
			inputJSON:  `{"sender": "+6282323231", "receiver": "+6282323232", "body": "hello"}`,
			statusCode: 201,
			sender:     "+6282323231",
			receiver:   "+6282323232",
			body:       "hello",
			errMessage: "",
		},
		{
			inputJSON:  `{"sender": "", "receiver": "+6282387325972", "body": "hello"}`,
			statusCode: 422,
			errMessage: "Required Sender",
		},
		{
			inputJSON:  `{"sender": "+6282387325971", "receiver": "", "body": "hello"}`,
			statusCode: 422,
			errMessage: "Required Receiver",
		},
		{
			inputJSON:  `{"sender": "+6282387325971", "receiver": "+6282387325972", "body": ""}`,
			statusCode: 422,
			errMessage: "Required Body",
		},
		{
			inputJSON:  `{"sender": 1231231, "receiver": "+6282387325972", "body": "hello"}`,
			statusCode: 422,
			errMessage: "invalid json body",
		},
		{
			inputJSON:  `{"sender": "+6282387325971", "receiver": 123das, "body": "hello"}`,
			statusCode: 422,
			errMessage: "invalid json body",
		},
		{
			inputJSON:  `{"sender": "+6282387325971", "receiver": "+6282387325972", "body": 123dsada"}`,
			statusCode: 422,
			errMessage: "invalid json body",
		},
	}

	for _, v := range samples {
		r := chi.NewRouter()
		r.Post("/api/v1/chats", controllers.CreateChat)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/chats", bytes.NewBufferString(v.inputJSON))
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
			assert.Equal(t, responseMap["sender"], v.sender)
			assert.Equal(t, responseMap["receiver"], v.receiver)
			assert.Equal(t, responseMap["body"], v.body)
		}
		if v.statusCode == 400 || v.statusCode == 422 || v.statusCode == 500 && v.errMessage != "" {
			assert.Equal(t, responseMap["message"], v.errMessage)
		}
	}
}

func TestGetChatByID(t *testing.T) {
	database()
	err := refreshChatsTable()
	if err != nil {
		log.Fatal(err)
	}
	message, err := seedOneChat()
	if err != nil {
		t.Errorf("Error while seeding table: %s", err)
	}

	samples := []struct {
		id         string
		statusCode int
		sender     string
		receiver   string
		body       string
		errMessage string
	}{
		{
			id:         strconv.Itoa(int(message.Id)),
			statusCode: 200,
			sender:     message.Sender,
			receiver:   message.Receiver,
			body:       message.Body,
			errMessage: "",
		},
		{
			id:         "unknwon",
			statusCode: 400,
			errMessage: "chat id should be a number",
		},
		{
			id:         strconv.Itoa(12322), //an id that does not exist
			statusCode: 404,
			errMessage: "no record matching given id",
		},
	}

	for _, v := range samples {
		r := chi.NewRouter()
		r.Get("/api/v1/chats/{chat_id}", controllers.GetChat)
		req, err := http.NewRequest(http.MethodGet, "/api/v1/chats/"+v.id, nil)
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
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			//casting the interface to map:
			log.Println("expected")
			log.Println(responseMap["sender"])
			log.Println(responseMap["receiver"])
			log.Println(responseMap["body"])
			log.Println("actual")
			log.Println(v.sender)
			log.Println(v.receiver)
			log.Println(v.body)
			assert.Equal(t, responseMap["sender"], v.sender)
			assert.Equal(t, responseMap["receiver"], v.receiver)
			assert.Equal(t, responseMap["body"], v.body)
		}
		if v.statusCode == 400 || v.statusCode == 422 && v.errMessage != "" {
			assert.Equal(t, responseMap["message"], v.errMessage)
		}
	}
}

func TestUpdateChat(t *testing.T) {
	database()
	err := refreshChatsTable()
	if err != nil {
		log.Fatal(err)
	}
	messages, err := seedChats()
	if err != nil {
		t.Errorf("Error while seeding table: %s", err)
	}

	//Get only the first message id
	firstId := messages[0].Id

	samples := []struct {
		id         string
		inputJSON  string
		statusCode int
		body       string
		errMessage string
	}{
		{
			id:         strconv.Itoa(int(firstId)),
			inputJSON:  `{"body": "update body"}`,
			statusCode: 200,
			body:       "update body",
			errMessage: "",
		},
		{
			//Empty body
			id:         strconv.Itoa(int(firstId)),
			inputJSON:  `{"body": ""}`,
			statusCode: 422,
			errMessage: "Required Body",
		},
		{
			//when an integer is used like a string for body
			id:         strconv.Itoa(int(firstId)),
			inputJSON:  `{"body": 123453}`,
			statusCode: 422,
			errMessage: "invalid json body",
		},
		{
			id:         "unknwon",
			statusCode: 400,
			errMessage: "chat id should be a number",
		},
		{
			id:         strconv.Itoa(12322), //an id that does not exist
			inputJSON:  `{"body": "the body"}`,
			statusCode: 404,
			errMessage: "no record matching given id",
		},
	}
	for _, v := range samples {
		r := chi.NewRouter()
		r.Put("/api/v1/chats/{chat_id}", controllers.UpdateChat)
		req, err := http.NewRequest(http.MethodPut, "/api/v1/chats/"+v.id, bytes.NewBufferString(v.inputJSON))
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
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			//casting the interface to map:
			assert.Equal(t, "", v.errMessage)
			assert.Equal(t, responseMap["body"], v.body)
		}
		if v.statusCode == 400 || v.statusCode == 422 || v.statusCode == 500 && v.errMessage != "" {
			assert.Equal(t, responseMap["message"], v.errMessage)
		}
	}
}

func TestGetAllMessage(t *testing.T) {
	database()
	err := refreshChatsTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedChats()
	if err != nil {
		t.Errorf("Error while seeding table: %s", err)
	}
	r := chi.NewRouter()
	r.Get("/api/v1/chats", controllers.GetAllChats)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/chats", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var msgs []domain.Chat

	err = json.Unmarshal(rr.Body.Bytes(), &msgs)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(msgs), 2)
}

func TestDeleteMessage(t *testing.T) {
	database()
	err := refreshChatsTable()
	if err != nil {
		log.Fatal(err)
	}
	message, err := seedOneChat()
	if err != nil {
		t.Errorf("Error while seeding table: %s", err)
	}
	samples := []struct {
		id         string
		statusCode int
		status     string
		errMessage string
	}{
		{
			id:         strconv.Itoa(int(message.Id)),
			statusCode: 200,
			status:     "deleted",
			errMessage: "",
		},
		{
			id:         "unknwon",
			statusCode: 400,
			errMessage: "chat id should be a number",
		},
		{
			id:         strconv.Itoa(12322), //an id that does not exist
			statusCode: 404,
			errMessage: "no record matching given id",
		},
	}
	for _, v := range samples {
		r := chi.NewRouter()
		r.Delete("/api/v1/chats/{chat_id}", controllers.DeleteChat)
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/chats/"+v.id, nil)
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
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			//casting the interface to map:
			assert.Equal(t, responseMap["status"], v.status)
		}
		if v.statusCode == 400 || v.statusCode == 422 && v.errMessage != "" {
			assert.Equal(t, responseMap["message"], v.errMessage)
		}
	}
}

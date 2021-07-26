package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/SemmiDev/lets-tests/domain"
	"github.com/SemmiDev/lets-tests/services"
	"github.com/SemmiDev/lets-tests/utils"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	getChatService    func(chatId int64) (*domain.Chat, utils.ChatErr)
	createChatService func(message *domain.Chat) (*domain.Chat, utils.ChatErr)
	updateChatService func(message *domain.Chat) (*domain.Chat, utils.ChatErr)
	deleteChatService func(chatId int64) utils.ChatErr
	getAllChatService func() ([]domain.Chat, utils.ChatErr)
)

type serviceMock struct{}

func (sm *serviceMock) GetChat(chatId int64) (*domain.Chat, utils.ChatErr) {
	return getChatService(chatId)
}
func (sm *serviceMock) CreateChat(message *domain.Chat) (*domain.Chat, utils.ChatErr) {
	return createChatService(message)
}
func (sm *serviceMock) UpdateChat(message *domain.Chat) (*domain.Chat, utils.ChatErr) {
	return updateChatService(message)
}
func (sm *serviceMock) DeleteChat(chatId int64) utils.ChatErr {
	return deleteChatService(chatId)
}
func (sm *serviceMock) GetAllChats() ([]domain.Chat, utils.ChatErr) {
	return getAllChatService()
}

func TestGetChat_Success(t *testing.T) {
	services.ChatsService = &serviceMock{}

	sender := utils.RandomSender()
	receiver := utils.RandomReceiver()
	body := utils.RandomBody()
	now := time.Now()

	getChatService = func(chatId int64) (*domain.Chat, utils.ChatErr) {
		return &domain.Chat{
			Id:        1,
			Sender:    sender,
			Receiver:  receiver,
			Body:      body,
			CreatedAt: now,
		}, nil
	}

	chatId := "1" //this has to be a string, because is passed through the url
	r := chi.NewRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/chats/"+chatId, nil)
	rr := httptest.NewRecorder()
	r.Get("/api/v1/chats/{chat_id}", GetChat)
	r.ServeHTTP(rr, req)

	var message domain.Chat
	err := json.Unmarshal(rr.Body.Bytes(), &message)

	log.Println("---------------------------")
	log.Println(message)
	log.Println("---------------------------")
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 1, message.Id)
	assert.EqualValues(t, sender, message.Sender)
	assert.EqualValues(t, receiver, message.Receiver)
	assert.EqualValues(t, body, message.Body)
}

func TestGetChat_Invalid_Id(t *testing.T) {
	chatId := "abc" //this has to be a string, because is passed through the url
	r := chi.NewRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/chats/"+chatId, nil)
	rr := httptest.NewRecorder()
	r.Get("/api/v1/chats/{chat_id}", GetChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "chat id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func TestGet_Chat_Not_Found(t *testing.T) {
	services.ChatsService = &serviceMock{}
	getChatService = func(chatId int64) (*domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.NotFoundError, "chat not found")
	}

	chatId := "1" //valid id
	r := chi.NewRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/chats/"+chatId, nil)
	rr := httptest.NewRecorder()
	r.Get("/api/v1/chats/{chat_id}", GetChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusNotFound, apiErr.Status())
	assert.EqualValues(t, "chat not found", apiErr.Message())
	assert.EqualValues(t, "not_found", apiErr.Error())
}

func TestGetChat_Chat_Database_Error(t *testing.T) {
	services.ChatsService = &serviceMock{}
	getChatService = func(chatId int64) (*domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.InternalServerError, "database error")
	}
	chatId := "1" //valid id
	r := chi.NewRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/chats/"+chatId, nil)
	rr := httptest.NewRecorder()
	r.Get("/api/v1/chats/{chat_id}", GetChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "database error", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
}

func TestCreateChat_Success(t *testing.T) {
	services.ChatsService = &serviceMock{}

	sender := utils.RandomSender()
	receiver := utils.RandomReceiver()
	body := utils.RandomBody()
	now := time.Now()

	createChatService = func(message *domain.Chat) (*domain.Chat, utils.ChatErr) {
		return &domain.Chat{
			Id:        1,
			Sender:    sender,
			Receiver:  receiver,
			Body:      body,
			CreatedAt: now,
		}, nil
	}

	jsonBody := `{"sender": "+6282323232", "receiver": "+6282323232", "body": "hello"}`
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/chats", bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Post("/api/v1/chats", CreateChat)
	r.ServeHTTP(rr, req)

	var message domain.Chat
	err = json.Unmarshal(rr.Body.Bytes(), &message)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.EqualValues(t, http.StatusCreated, rr.Code)
	assert.EqualValues(t, 1, message.Id)
	assert.EqualValues(t, sender, message.Sender)
	assert.EqualValues(t, receiver, message.Receiver)
	assert.EqualValues(t, body, message.Body)
}

func TestCreateChat_Invalid_Json(t *testing.T) {
	inputJson := `{"sender": 1234, "receiver": "+1231231231", "body": "hello"}`
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/chats", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Post("/api/v1/chats", CreateChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "invalid json body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

//This test is not really necessary here, because it has been handled in the service test
func TestCreateChat_Empty_Body(t *testing.T) {
	services.ChatsService = &serviceMock{}
	createChatService = func(message *domain.Chat) (*domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.UnprocessableEntityError, "Required Body")
	}

	inputJson := `{"sender": "+6282323232", "receiver": "+6282323232", "body": ""}`
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/chats", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Post("/api/v1/chats", CreateChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "Required Body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

func TestCreateChat_Empty_Sender(t *testing.T) {
	services.ChatsService = &serviceMock{}
	createChatService = func(message *domain.Chat) (*domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.UnprocessableEntityError, "Required Sender")
	}

	inputJson := `{"sender": "", "receiver": "+6282323232", "body": "hello"}`
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/chats", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Post("/api/v1/chats", CreateChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "Required Sender", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

func TestCreateChat_Empty_Receiver(t *testing.T) {
	services.ChatsService = &serviceMock{}
	createChatService = func(message *domain.Chat) (*domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.UnprocessableEntityError, "Required Receiver")
	}

	inputJson := `{"sender": "+6282323232", "receiver": "", "body": "hello"}`
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/chats", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Post("/api/v1/chats", CreateChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "Required Receiver", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

func TestCreateChat_Same_Sender_Receiver(t *testing.T) {
	services.ChatsService = &serviceMock{}
	createChatService = func(message *domain.Chat) (*domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.UnprocessableEntityError, "Sender and Receiver must different")
	}

	inputJson := `{"sender": "+6213131312", "receiver": "+6213131312", "body": "hello"}`
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/chats", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Post("/api/v1/chats", CreateChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "Sender and Receiver must different", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

func TestCreateChat_Not_Valid_Sender(t *testing.T) {
	services.ChatsService = &serviceMock{}
	createChatService = func(message *domain.Chat) (*domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.UnprocessableEntityError, "Invalid Sender Phone Number")
	}

	inputJson := `{"sender": "xxxx", "receiver": "+6213131312", "body": "hello"}`
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/chats", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Post("/api/v1/chats", CreateChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "Invalid Sender Phone Number", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

func TestCreateChat_Not_Valid_Receiver(t *testing.T) {
	services.ChatsService = &serviceMock{}
	createChatService = func(message *domain.Chat) (*domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.UnprocessableEntityError, "Invalid Receiver Phone Number")
	}

	inputJson := `{"sender": "+6213131312", "receiver": "xxxx", "body": "hello"}`
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/chats", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Post("/api/v1/chats", CreateChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "Invalid Receiver Phone Number", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

func TestUpdateChat_Success(t *testing.T) {
	services.ChatsService = &serviceMock{}

	sender := utils.RandomSender()
	receiver := utils.RandomReceiver()

	updateChatService = func(message *domain.Chat) (*domain.Chat, utils.ChatErr) {
		return &domain.Chat{
			Id:       1,
			Sender:   sender,
			Receiver: receiver,
			Body:     "update body",
		}, nil
	}

	jsonBody := `{"body": "update body"}`
	r := chi.NewRouter()
	id := "1"
	req, err := http.NewRequest(http.MethodPut, "/api/v1/chats/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Put("/api/v1/chats/{chat_id}", UpdateChat)
	r.ServeHTTP(rr, req)

	var message domain.Chat
	err = json.Unmarshal(rr.Body.Bytes(), &message)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 1, message.Id)
	assert.EqualValues(t, "update body", message.Body)
}

//We dont need to mock the service method here, because we wont call it
func TestUpdateChat_Invalid_Id(t *testing.T) {
	jsonBody := `{"body": "update body"}`
	r := chi.NewRouter()
	id := "abc"
	req, err := http.NewRequest(http.MethodPut, "/api/v1/chats/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Put("/api/v1/chats/{chat_id}", UpdateChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "chat id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

//When for instance an integer is provided instead of a string
func TestUpdateChat_Invalid_Json(t *testing.T) {
	inputJson := `{"body": 21231}`
	r := chi.NewRouter()
	id := "1"
	req, err := http.NewRequest(http.MethodPut, "/api/v1/chats/"+id, bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Put("/api/v1/chats/{chat_id}", UpdateChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "invalid json body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

//This test is not really necessary here, because it has been handled in the service test
func TestUpdateChat_Empty_Body(t *testing.T) {
	services.ChatsService = &serviceMock{}
	updateChatService = func(message *domain.Chat) (*domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.UnprocessableEntityError, "Required Body")
	}
	inputJson := `{"body": ""}`
	id := "1"
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPut, "/api/v1/chats/"+id, bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Put("/api/v1/chats/{chat_id}", UpdateChat)
	r.ServeHTTP(rr, req)
	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "Required Body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

//Other errors can happen when we try to update the message
func TestUpdateChat_Error_Updating(t *testing.T) {
	services.ChatsService = &serviceMock{}
	updateChatService = func(message *domain.Chat) (*domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.InternalServerError, "error when updating chat")
	}
	jsonBody := `{"body": "update body"}`
	r := chi.NewRouter()
	id := "1"
	req, err := http.NewRequest(http.MethodPut, "/api/v1/chats/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Put("/api/v1/chats/{chat_id}", UpdateChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)

	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "error when updating chat", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
}

func TestDeleteChat_Success(t *testing.T) {
	services.ChatsService = &serviceMock{}
	deleteChatService = func(msg int64) utils.ChatErr {
		return nil
	}
	r := chi.NewRouter()
	id := "1"
	req, err := http.NewRequest(http.MethodDelete, "/api/v1/chats/"+id, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Delete("/api/v1/chats/{chat_id}", DeleteChat)
	r.ServeHTTP(rr, req)

	var response = make(map[string]string)
	theErr := json.Unmarshal(rr.Body.Bytes(), &response)
	if theErr != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, response["status"], "deleted")
}

func TestDeleteChat_Invalid_Id(t *testing.T) {
	r := chi.NewRouter()
	id := "abc"
	req, err := http.NewRequest(http.MethodDelete, "/api/v1/chats/"+id, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Delete("/api/v1/chats/{chat_id}", DeleteChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "chat id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func TestDeleteChat_Failure(t *testing.T) {
	services.ChatsService = &serviceMock{}
	deleteChatService = func(msg int64) utils.ChatErr {
		return utils.ErrorKind(utils.InternalServerError, "error deleting chat")
	}
	r := chi.NewRouter()
	id := "1"
	req, err := http.NewRequest(http.MethodDelete, "/api/v1/chats/"+id, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Delete("/api/v1/chats/{chat_id}", DeleteChat)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "error deleting chat", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
}

func TestGetAllChats_Success(t *testing.T) {
	services.ChatsService = &serviceMock{}

	sender1 := utils.RandomSender()
	receiver1 := utils.RandomReceiver()
	body1 := utils.RandomBody()

	sender2 := utils.RandomSender()
	receiver2 := utils.RandomReceiver()
	body2 := utils.RandomBody()

	getAllChatService = func() ([]domain.Chat, utils.ChatErr) {
		return []domain.Chat{
			{
				Id:       1,
				Sender:   sender1,
				Receiver: receiver1,
				Body:     body1,
			},
			{
				Id:       2,
				Sender:   sender2,
				Receiver: receiver2,
				Body:     body2,
			},
		}, nil
	}
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/chats/", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Get("/api/v1/chats/", GetAllChats)
	r.ServeHTTP(rr, req)

	var messages []domain.Chat
	theErr := json.Unmarshal(rr.Body.Bytes(), &messages)
	if theErr != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	assert.Nil(t, err)
	assert.NotNil(t, messages)

	assert.EqualValues(t, messages[0].Id, 1)
	assert.EqualValues(t, messages[0].Sender, sender1)
	assert.EqualValues(t, messages[0].Receiver, receiver1)
	assert.EqualValues(t, messages[0].Body, body1)

	assert.EqualValues(t, messages[1].Id, 2)
	assert.EqualValues(t, messages[1].Sender, sender2)
	assert.EqualValues(t, messages[1].Receiver, receiver2)
	assert.EqualValues(t, messages[1].Body, body2)
}

//For any reason we could not get the messages
func TestGetAllChats_Failure(t *testing.T) {
	services.ChatsService = &serviceMock{}
	getAllChatService = func() ([]domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.InternalServerError, "error getting chats")
	}

	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/chats/", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Get("/api/v1/chats/", GetAllChats)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, "error getting chats", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
}

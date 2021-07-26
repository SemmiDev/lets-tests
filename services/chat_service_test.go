package services

import (
	"database/sql"
	"fmt"
	"github.com/SemmiDev/lets-tests/domain"
	"github.com/SemmiDev/lets-tests/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var (
	sender   = utils.RandomSender()
	receiver = utils.RandomReceiver()
	body     = utils.RandomBody()
	now      = time.Now()

	getChatDomain     func(chatId int64) (*domain.Chat, utils.ChatErr)
	createChatDomain  func(msg *domain.Chat) (*domain.Chat, utils.ChatErr)
	updateChatDomain  func(msg *domain.Chat) (*domain.Chat, utils.ChatErr)
	deleteChatDomain  func(chatId int64) utils.ChatErr
	getAllChatsDomain func() ([]domain.Chat, utils.ChatErr)
)

type getDBMock struct{}

func (m *getDBMock) Get(chatId int64) (*domain.Chat, utils.ChatErr) {
	return getChatDomain(chatId)
}
func (m *getDBMock) Create(msg *domain.Chat) (*domain.Chat, utils.ChatErr) {
	return createChatDomain(msg)
}
func (m *getDBMock) Update(msg *domain.Chat) (*domain.Chat, utils.ChatErr) {
	return updateChatDomain(msg)
}
func (m *getDBMock) Delete(chatId int64) utils.ChatErr {
	return deleteChatDomain(chatId)
}
func (m *getDBMock) GetAll() ([]domain.Chat, utils.ChatErr) {
	return getAllChatsDomain()
}
func (m *getDBMock) Initialize(string, string, string, string, string, string) *sql.DB {
	return nil
}

func TestChatsService_GetChat_Success(t *testing.T) {
	domain.ChatRepo = &getDBMock{} //this is where we swapped the functionality

	getChatDomain = func(chatId int64) (*domain.Chat, utils.ChatErr) {
		return &domain.Chat{
			Id:        1,
			Sender:    sender,
			Receiver:  receiver,
			Body:      body,
			CreatedAt: now,
		}, nil
	}

	msg, err := ChatsService.GetChat(1)

	fmt.Println("this is the chat: ", msg)
	assert.NotNil(t, msg)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, msg.Id)
	assert.EqualValues(t, sender, msg.Sender)
	assert.EqualValues(t, receiver, msg.Receiver)
	assert.EqualValues(t, body, msg.Body)
	assert.EqualValues(t, now, msg.CreatedAt)
}

func TestChatsService_GetChat_NotFoundID(t *testing.T) {
	domain.ChatRepo = &getDBMock{}

	getChatDomain = func(chatId int64) (*domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.NotFoundError, "the id is not found")
	}

	msg, err := ChatsService.GetChat(1)
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
	assert.EqualValues(t, "the id is not found", err.Message())
	assert.EqualValues(t, "not_found", err.Error())
}

func TestChatsService_CreateChat_Success(t *testing.T) {
	domain.ChatRepo = &getDBMock{}

	createChatDomain = func(msg *domain.Chat) (*domain.Chat, utils.ChatErr) {
		return &domain.Chat{
			Id:        1,
			Sender:    sender,
			Receiver:  receiver,
			Body:      body,
			CreatedAt: now,
		}, nil
	}

	request := &domain.Chat{
		Sender:    sender,
		Receiver:  receiver,
		Body:      body,
		CreatedAt: now,
	}

	msg, err := ChatsService.CreateChat(request)
	fmt.Println("this is the chat: ", msg)
	assert.NotNil(t, msg)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, msg.Id)
	assert.EqualValues(t, sender, msg.Sender)
	assert.EqualValues(t, receiver, msg.Receiver)
	assert.EqualValues(t, body, msg.Body)
	assert.EqualValues(t, now, msg.CreatedAt)
}

func TestChatsService_CreateChat_Invalid_Request(t *testing.T) {
	tests := []struct {
		request    *domain.Chat
		statusCode int
		errMsg     string
		errErr     string
	}{
		{
			request: &domain.Chat{
				Sender:    "",
				Receiver:  receiver,
				Body:      body,
				CreatedAt: now,
			},
			statusCode: http.StatusUnprocessableEntity,
			errMsg:     "Required Sender",
			errErr:     "invalid_request",
		},
		{
			request: &domain.Chat{
				Sender:    sender,
				Receiver:  "",
				Body:      body,
				CreatedAt: now,
			},
			statusCode: http.StatusUnprocessableEntity,
			errMsg:     "Required Receiver",
			errErr:     "invalid_request",
		},
		{
			request: &domain.Chat{
				Sender:    sender,
				Receiver:  receiver,
				Body:      "",
				CreatedAt: now,
			},
			statusCode: http.StatusUnprocessableEntity,
			errMsg:     "Required Body",
			errErr:     "invalid_request",
		},
		{
			request: &domain.Chat{
				Sender:    sender,
				Receiver:  "hemhemhem",
				Body:      body,
				CreatedAt: now,
			},
			statusCode: http.StatusUnprocessableEntity,
			errMsg:     "Invalid Receiver Phone Number",
			errErr:     "invalid_request",
		},
		{
			request: &domain.Chat{
				Sender:    "hemhemhem",
				Receiver:  receiver,
				Body:      body,
				CreatedAt: now,
			},
			statusCode: http.StatusUnprocessableEntity,
			errMsg:     "Invalid Sender Phone Number",
			errErr:     "invalid_request",
		},
	}
	for _, tt := range tests {
		msg, err := ChatsService.CreateChat(tt.request)
		assert.Nil(t, msg)
		assert.NotNil(t, err)
		assert.EqualValues(t, tt.errMsg, err.Message())
		assert.EqualValues(t, tt.statusCode, err.Status())
		assert.EqualValues(t, tt.errErr, err.Error())
	}
}

func TestChatsService_UpdateChat_Success(t *testing.T) {
	domain.ChatRepo = &getDBMock{}
	getChatDomain = func(chatId int64) (*domain.Chat, utils.ChatErr) {
		return &domain.Chat{
			Id:        1,
			Sender:    sender,
			Receiver:  receiver,
			Body:      body,
			CreatedAt: now,
		}, nil
	}

	newBody := utils.RandomBody()
	updateChatDomain = func(msg *domain.Chat) (*domain.Chat, utils.ChatErr) {
		return &domain.Chat{
			Id:        1,
			Sender:    sender,
			Receiver:  receiver,
			Body:      newBody,
			CreatedAt: now,
		}, nil
	}

	request := &domain.Chat{
		Body: newBody,
	}

	msg, err := ChatsService.UpdateChat(request)
	assert.NotNil(t, msg)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, msg.Id)
	assert.EqualValues(t, newBody, msg.Body)
}

func TestChatsService_UpdateMessage_Empty_Body(t *testing.T) {
	tests := []struct {
		request    *domain.Chat
		statusCode int
		errMsg     string
		errErr     string
	}{
		{
			request: &domain.Chat{
				Body: "",
			},
			statusCode: http.StatusUnprocessableEntity,
			errMsg:     "Required Body",
			errErr:     "invalid_request",
		},
	}
	for _, tt := range tests {
		msg, err := ChatsService.UpdateChat(tt.request)
		assert.Nil(t, msg)
		assert.NotNil(t, err)
		assert.EqualValues(t, tt.statusCode, err.Status())
		assert.EqualValues(t, tt.errMsg, err.Message())
		assert.EqualValues(t, tt.errErr, err.Error())
	}
}

func TestChatsService_UpdateMessage_Failure_Getting_Former_Chat(t *testing.T) {
	domain.ChatRepo = &getDBMock{}
	getChatDomain = func(chatId int64) (*domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.InternalServerError, "error getting chat")
	}

	request := &domain.Chat{
		Sender:    sender,
		Receiver:  receiver,
		Body:      utils.RandomBody(),
		CreatedAt: time.Time{},
	}

	msg, err := ChatsService.UpdateChat(request)
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.EqualValues(t, "error getting chat", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "server_error", err.Error())
}

func TestChatsService_UpdateMessage_Failure_Updating_Chat(t *testing.T) {
	domain.ChatRepo = &getDBMock{}

	getChatDomain = func(chatId int64) (*domain.Chat, utils.ChatErr) {
		return &domain.Chat{
			Id:       1,
			Sender:   sender,
			Receiver: receiver,
		}, nil
	}

	updateChatDomain = func(msg *domain.Chat) (*domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.InternalServerError, "error updating message")
	}

	request := &domain.Chat{
		Body: utils.RandomBody(),
	}

	msg, err := ChatsService.UpdateChat(request)
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.EqualValues(t, "error updating message", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "server_error", err.Error())
}

func TestChatsService_DeleteChat_Success(t *testing.T) {
	domain.ChatRepo = &getDBMock{}
	getChatDomain = func(chatId int64) (*domain.Chat, utils.ChatErr) {
		return &domain.Chat{
			Id:       1,
			Sender:   sender,
			Receiver: receiver,
		}, nil
	}

	deleteChatDomain = func(chatId int64) utils.ChatErr {
		return nil
	}

	err := ChatsService.DeleteChat(1)
	assert.Nil(t, err)
}

func TestChatsService_DeleteMessage_Error_Getting_Chat(t *testing.T) {
	domain.ChatRepo = &getDBMock{}
	getChatDomain = func(chatId int64) (*domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.InternalServerError, "Something went wrong getting chat")
	}

	err := ChatsService.DeleteChat(1)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Something went wrong getting chat", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "server_error", err.Error())
}

func TestChatsService_DeleteMessage_Error_Deleting_Chat(t *testing.T) {
	domain.ChatRepo = &getDBMock{}
	getChatDomain = func(chatId int64) (*domain.Chat, utils.ChatErr) {
		return &domain.Chat{
			Id:       1,
			Sender:   sender,
			Receiver: receiver,
			Body:     body,
		}, nil
	}

	deleteChatDomain = func(chatId int64) utils.ChatErr {
		return utils.ErrorKind(utils.InternalServerError, "error deleting chat")
	}

	err := ChatsService.DeleteChat(1)
	assert.NotNil(t, err)
	assert.EqualValues(t, "error deleting chat", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "server_error", err.Error())
}

func TestChatsService_GetAllChats(t *testing.T) {
	domain.ChatRepo = &getDBMock{}

	newSender := utils.RandomSender()
	newReceiver := utils.RandomReceiver()
	newBody := utils.RandomBody()

	getAllChatsDomain = func() ([]domain.Chat, utils.ChatErr) {
		return []domain.Chat{
			{
				Id:       1,
				Sender:   sender,
				Receiver: receiver,
				Body:     body,
			},
			{
				Id:       2,
				Sender:   newSender,
				Receiver: newReceiver,
				Body:     newBody,
			},
		}, nil
	}

	messages, err := ChatsService.GetAllChats()
	assert.Nil(t, err)
	assert.NotNil(t, messages)
	assert.EqualValues(t, messages[0].Id, 1)
	assert.EqualValues(t, messages[0].Sender, sender)
	assert.EqualValues(t, messages[0].Receiver, receiver)
	assert.EqualValues(t, messages[0].Body, body)

	assert.EqualValues(t, messages[1].Id, 2)
	assert.EqualValues(t, messages[1].Sender, newSender)
	assert.EqualValues(t, messages[1].Receiver, newReceiver)
	assert.EqualValues(t, messages[1].Body, newBody)
}

func TestChatsService_GetAllChats_Error_Getting_Chats(t *testing.T) {
	domain.ChatRepo = &getDBMock{}
	getAllChatsDomain = func() ([]domain.Chat, utils.ChatErr) {
		return nil, utils.ErrorKind(utils.InternalServerError, "error getting chats")
	}

	messages, err := ChatsService.GetAllChats()
	assert.NotNil(t, err)
	assert.Nil(t, messages)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error getting chats", err.Message())
	assert.EqualValues(t, "server_error", err.Error())
}

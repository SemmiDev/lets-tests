package services

import (
	"github.com/SemmiDev/lets-tests/domain"
	"github.com/SemmiDev/lets-tests/utils"
	"time"
)

var (
	ChatsService chatServiceInterface = &chatsService{}
)

type chatsService struct{}

type chatServiceInterface interface {
	GetChat(int64) (*domain.Chat, utils.ChatErr)
	CreateChat(*domain.Chat) (*domain.Chat, utils.ChatErr)
	UpdateChat(*domain.Chat) (*domain.Chat, utils.ChatErr)
	DeleteChat(int64) utils.ChatErr
	GetAllChats() ([]domain.Chat, utils.ChatErr)
}

func (c *chatsService) GetChat(id int64) (*domain.Chat, utils.ChatErr) {
	message, err := domain.ChatRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (c *chatsService) CreateChat(chat *domain.Chat) (*domain.Chat, utils.ChatErr) {
	if err := chat.Validate(""); err != nil {
		return nil, err
	}
	chat.CreatedAt = time.Now()
	chat, err := domain.ChatRepo.Create(chat)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (c *chatsService) UpdateChat(chat *domain.Chat) (*domain.Chat, utils.ChatErr) {
	if err := chat.Validate("update"); err != nil {
		return nil, err
	}
	current, err := domain.ChatRepo.Get(chat.Id)
	if err != nil {
		return nil, err
	}

	current.Body = chat.Body
	updateMsg, err := domain.ChatRepo.Update(current)
	if err != nil {
		return nil, err
	}
	return updateMsg, nil
}

func (c *chatsService) DeleteChat(chatId int64) utils.ChatErr {
	msg, err := domain.ChatRepo.Get(chatId)
	if err != nil {
		return err
	}
	deleteErr := domain.ChatRepo.Delete(msg.Id)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

func (c *chatsService) GetAllChats() ([]domain.Chat, utils.ChatErr) {
	chats, err := domain.ChatRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return chats, nil
}

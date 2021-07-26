package controllers

import (
	"encoding/json"
	"github.com/SemmiDev/lets-tests/domain"
	"github.com/SemmiDev/lets-tests/services"
	"github.com/SemmiDev/lets-tests/utils"
	"net/http"
)

func CreateChat(w http.ResponseWriter, r *http.Request) {
	var chat domain.Chat
	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		theErr := utils.ErrorKind(utils.UnprocessableEntityError, "invalid json body")
		MarshalError(w, theErr.Status(), theErr)
		return
	}

	res, theErr := services.ChatsService.CreateChat(&chat)
	if theErr != nil {
		MarshalError(w, theErr.Status(), theErr)
		return
	}

	MarshallSuccess(w, http.StatusCreated, "CREATED", res)
	return
}

func GetChat(w http.ResponseWriter, r *http.Request) {
	chatId, err := GetUrlPathInt64(r, "chat_id")
	if err != nil {
		MarshalError(w, err.Status(), err)
		return
	}

	chat, getErr := services.ChatsService.GetChat(chatId)
	if getErr != nil {
		MarshalError(w, getErr.Status(), getErr)
		return
	}

	MarshallSuccess(w, http.StatusOK, "OK", chat)
	return
}

func GetAllChats(w http.ResponseWriter, _ *http.Request) {
	chats, getErr := services.ChatsService.GetAllChats()
	if getErr != nil {
		MarshalError(w, getErr.Status(), getErr)
		return
	}

	MarshallSuccess(w, http.StatusOK, "OK", chats)
	return
}

func UpdateChat(w http.ResponseWriter, r *http.Request) {
	chatId, err := GetUrlPathInt64(r, "chat_id")
	if err != nil {
		MarshalError(w, err.Status(), err)
		return
	}

	var req domain.UpdateChatRequest
	reqErr := json.NewDecoder(r.Body).Decode(&req)
	if reqErr != nil {
		theErr := utils.ErrorKind(utils.UnprocessableEntityError, "invalid json body")
		MarshalError(w, theErr.Status(), theErr)
		return
	}

	chat := domain.Chat{
		Id:   chatId,
		Body: req.Body,
	}
	update, theErr := services.ChatsService.UpdateChat(&chat)
	if theErr != nil {
		MarshalError(w, theErr.Status(), theErr)
		return
	}

	MarshallSuccess(w, http.StatusOK, "OK", update)
	return
}

func DeleteChat(w http.ResponseWriter, r *http.Request) {
	chatId, err := GetUrlPathInt64(r, "chat_id")
	if err != nil {
		MarshalError(w, err.Status(), err)
		return
	}

	err = services.ChatsService.DeleteChat(chatId)
	if err != nil {
		MarshalError(w, err.Status(), err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "deleted",
	})
	return
}

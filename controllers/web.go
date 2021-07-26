package controllers

import (
	"encoding/json"
	"github.com/SemmiDev/lets-tests/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Error  bool        `json:"error"`
	Data   interface{} `json:"data"`
}

func GetUrlPathInt64(r *http.Request, key string) (int64, utils.ChatErr) {
	chatId, err := strconv.ParseInt(chi.URLParam(r, key), 10, 64)
	if err != nil {
		return 0, utils.ErrorKind(utils.BadRequestError, "chat id should be a number")
	}
	return chatId, nil
}

func MarshalError(w http.ResponseWriter, code int, err utils.ChatErr) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(err)
}

func MarshallSuccess(w http.ResponseWriter, code int, status string, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

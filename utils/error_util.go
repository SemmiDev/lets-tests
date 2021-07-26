package utils

import (
	"encoding/json"
	"net/http"
)

type ChatErr interface {
	Message() string
	Status() int
	Error() string
}

type chatErr struct {
	ErrMessage string `json:"message"`
	ErrStatus  int    `json:"status"`
	ErrError   string `json:"error"`
}

func (e *chatErr) Error() string {
	return e.ErrError
}

func (e *chatErr) Message() string {
	return e.ErrMessage
}

func (e *chatErr) Status() int {
	return e.ErrStatus
}

type ErrKind string

const (
	NotFoundError            ErrKind = "NotFoundError"
	BadRequestError          ErrKind = "BadRequestError"
	UnprocessableEntityError ErrKind = "UnprocessableEntityError"
	InternalServerError      ErrKind = "InternalServerError"
)

func ErrorKind(errKind ErrKind, chat string) ChatErr {
	switch errKind {
	case NotFoundError:
		return notFound(chat)
	case BadRequestError:
		return badRequest(chat)
	case UnprocessableEntityError:
		return unprocessableEntity(chat)
	case InternalServerError:
		return internalServer(chat)
	}
	return nil
}

func notFound(msg string) ChatErr {
	return &chatErr{
		ErrMessage: msg,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "not_found",
	}
}

func badRequest(msg string) ChatErr {
	return &chatErr{
		ErrMessage: msg,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   "bad_request",
	}
}

func unprocessableEntity(msg string) ChatErr {
	return &chatErr{
		ErrMessage: msg,
		ErrStatus:  http.StatusUnprocessableEntity,
		ErrError:   "invalid_request",
	}
}

func internalServer(msg string) ChatErr {
	return &chatErr{
		ErrMessage: msg,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "server_error",
	}
}

func NewApiErrFromBytes(body []byte) (ChatErr, error) {
	var result chatErr
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

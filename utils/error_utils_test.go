package utils

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestErrsKind(t *testing.T) {
	var tests = []struct {
		Name       string
		ErrKind    ErrKind
		ErrMessage string
		ErrStatus  int
		ErrError   string
	}{
		{
			Name:       "Not Found Error",
			ErrKind:    NotFoundError,
			ErrMessage: "not found",
			ErrStatus:  http.StatusNotFound,
			ErrError:   "not_found",
		},
		{
			Name:       "Bad Request Error",
			ErrKind:    BadRequestError,
			ErrMessage: "bad request",
			ErrStatus:  http.StatusBadRequest,
			ErrError:   "bad_request",
		},
		{
			Name:       "Unprocessable Entity Error",
			ErrKind:    UnprocessableEntityError,
			ErrMessage: "invalid request",
			ErrStatus:  http.StatusUnprocessableEntity,
			ErrError:   "invalid_request",
		},
		{
			Name:       "Internal Server Error",
			ErrKind:    InternalServerError,
			ErrMessage: "server error",
			ErrStatus:  http.StatusInternalServerError,
			ErrError:   "server_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := ErrorKind(tt.ErrKind, tt.ErrMessage)
			assert.Equal(t, tt.ErrStatus, got.Status())
			assert.Equal(t, tt.ErrMessage, got.Message())
			assert.Equal(t, tt.ErrError, got.Error())
		})
	}
}

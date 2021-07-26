package domain

import (
	"database/sql"
	"github.com/SemmiDev/lets-tests/utils"
	"regexp"
	"strings"
	"time"
)

type Chat struct {
	Id        int64     `json:"id"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateChatRequest struct {
	Body string `json:"body"`
}

func (m *Chat) Validate(kind interface{}) utils.ChatErr {
	kind = kind.(string)

	if kind == "update" {
		m.Body = strings.TrimSpace(m.Body)
		if m.Body == "" {
			return utils.ErrorKind(utils.UnprocessableEntityError, "Required Body")
		}
		return nil
	}

	m.Sender = strings.TrimSpace(m.Sender)
	m.Receiver = strings.TrimSpace(m.Receiver)
	m.Body = strings.TrimSpace(m.Body)

	if m.Sender == "" {
		return utils.ErrorKind(utils.UnprocessableEntityError, "Required Sender")
	}
	if m.Receiver == "" {
		return utils.ErrorKind(utils.UnprocessableEntityError, "Required Receiver")
	}
	if m.Body == "" {
		return utils.ErrorKind(utils.UnprocessableEntityError, "Required Body")
	}

	phoneRegexp := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

	if from := phoneRegexp.MatchString(m.Sender); from == false {
		return utils.ErrorKind(utils.UnprocessableEntityError, "Invalid Sender Phone Number")
	}
	if to := phoneRegexp.MatchString(m.Receiver); to == false {
		return utils.ErrorKind(utils.UnprocessableEntityError, "Invalid Receiver Phone Number")
	}

	if m.Sender == m.Receiver {
		return utils.ErrorKind(utils.UnprocessableEntityError, "Sender and Receiver must different")
	}

	return nil
}

type chatRepoInterface interface {
	Get(Id int64) (*Chat, utils.ChatErr)
	Create(chat *Chat) (*Chat, utils.ChatErr)
	Update(chat *Chat) (*Chat, utils.ChatErr)
	Delete(Id int64) utils.ChatErr
	GetAll() ([]Chat, utils.ChatErr)
	Initialize(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) *sql.DB
}

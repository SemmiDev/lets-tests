package domain

import (
	"database/sql"
	"fmt"
	. "github.com/SemmiDev/lets-tests/utils"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	queryInsertChat  = `INSERT INTO chats(id, sender, receiver, body, created_at) VALUES (?,?,?,?);`
	queryGetChat     = `SELECT id, sender, receiver, body, created_at FROM chats WHERE id=?;`
	queryUpdateChat  = `UPDATE chats SET body=? WHERE id=?;`
	queryDeleteChat  = `DELETE FROM chats WHERE id=?;`
	queryGetAllChats = `SELECT id, sender, receiver, body, created_at FROM chats;`
)

type chatRepo struct {
	db *sql.DB
}

var ChatRepo chatRepoInterface = &chatRepo{}

func NewChatRepository(db *sql.DB) chatRepoInterface {
	return &chatRepo{db: db}
}

func (m *chatRepo) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) *sql.DB {
	var err error
	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	m.db, err = sql.Open(Dbdriver, DBURL)
	if err != nil {
		log.Fatal("This is the error connecting to the database:", err)
	}
	log.Printf("We are connected to the %s database", Dbdriver)

	return m.db
}

func (m *chatRepo) Get(chatId int64) (*Chat, ChatErr) {
	stmt, err := m.db.Prepare(queryGetChat)
	if err != nil {
		return nil, ErrorKind(InternalServerError, fmt.Sprintf("Error when trying to prepare chat: %s", err.Error()))
	}
	defer stmt.Close()

	var msg Chat
	result := stmt.QueryRow(chatId)
	if getError := result.Scan(&msg.Id, &msg.Sender, &msg.Receiver, &msg.Body, &msg.CreatedAt); getError != nil {
		fmt.Println("this is the error man: ", getError)
		return nil, ParseError(getError)
	}
	log.Println(result)
	return &msg, nil
}

func (m *chatRepo) Create(msg *Chat) (*Chat, ChatErr) {
	stmt, err := m.db.Prepare(queryInsertChat)
	if err != nil {
		return nil, ErrorKind(InternalServerError, fmt.Sprintf("error when trying to prepare user to save: %s", err.Error()))
	}
	defer stmt.Close()
	insertResult, createErr := stmt.Exec(msg.Sender, msg.Receiver, msg.Body, msg.CreatedAt)
	if createErr != nil {
		return nil, ParseError(createErr)
	}
	msgId, err := insertResult.LastInsertId()
	if err != nil {
		return nil, ErrorKind(InternalServerError, fmt.Sprintf("error when trying to save chat: %s", err.Error()))
	}
	msg.Id = msgId
	return msg, nil
}

func (m *chatRepo) Update(msg *Chat) (*Chat, ChatErr) {
	stmt, err := m.db.Prepare(queryUpdateChat)
	if err != nil {
		return nil, ErrorKind(InternalServerError, fmt.Sprintf("error when trying to prepare user to update: %s", err.Error()))
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(msg.Body, msg.Id)
	if updateErr != nil {
		return nil, ParseError(updateErr)
	}
	return msg, nil
}

func (m *chatRepo) Delete(msgId int64) ChatErr {
	stmt, err := m.db.Prepare(queryDeleteChat)
	if err != nil {
		return ErrorKind(InternalServerError, fmt.Sprintf("error when trying to delete chat: %s", err.Error()))
	}
	defer stmt.Close()

	if _, err := stmt.Exec(msgId); err != nil {
		return ErrorKind(InternalServerError, fmt.Sprintf("error when trying to delete chat %s", err.Error()))
	}
	return nil
}

func (m *chatRepo) GetAll() ([]Chat, ChatErr) {
	stmt, err := m.db.Prepare(queryGetAllChats)
	if err != nil {
		return nil, ErrorKind(InternalServerError, fmt.Sprintf("Error when trying to prepare all chats: %s", err.Error()))
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, ParseError(err)
	}
	defer rows.Close()

	results := make([]Chat, 0)

	for rows.Next() {
		var msg Chat
		if getError := rows.Scan(&msg.Id, &msg.Sender, &msg.Receiver, &msg.Body, &msg.CreatedAt); getError != nil {
			return nil, ErrorKind(InternalServerError, fmt.Sprintf("Error when trying to get chat: %s", getError.Error()))
		}
		results = append(results, msg)
	}
	if len(results) == 0 {
		return nil, ErrorKind(NotFoundError, "no records found")
	}
	return results, nil
}

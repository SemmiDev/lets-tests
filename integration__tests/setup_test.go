package integration__tests

import (
	"database/sql"
	"github.com/SemmiDev/lets-tests/domain"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
	"time"
)

const (
	queryTruncateChat = "TRUNCATE TABLE chats;"
	queryInsertChat   = "INSERT INTO chats(sender,receiver, body, created_at) VALUES(?, ?, ?, ?);"
	queryGetAllChats  = "SELECT id, sender, receiver, body, created_at FROM chats;"
)

var dbConn *sql.DB

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("./../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	os.Exit(m.Run())
}

func database() {
	dbDriver := os.Getenv("DBDRIVER_TEST")
	username := os.Getenv("USERNAME_TEST")
	password := os.Getenv("PASSWORD_TEST")
	host := os.Getenv("HOST_TEST")
	database := os.Getenv("DATABASE_TEST")
	port := os.Getenv("PORT_TEST")

	dbConn = domain.ChatRepo.Initialize(dbDriver, username, password, port, host, database)
}

func refreshChatsTable() error {
	stmt, err := dbConn.Prepare(queryTruncateChat)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("Error truncating messages table: %s", err)
	}
	return nil
}

func seedOneChat() (domain.Chat, error) {
	msg := domain.Chat{
		Sender:    "+6282387325971",
		Receiver:  "+6282387325972",
		Body:      "hello",
		CreatedAt: time.Now(),
	}

	stmt, err := dbConn.Prepare(queryInsertChat)
	if err != nil {
		panic(err.Error())
	}
	insertResult, createErr := stmt.Exec(msg.Sender, msg.Receiver, msg.Body, msg.CreatedAt)
	if createErr != nil {
		log.Fatalf("Error creating message: %s", createErr)
	}
	msgId, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("Error creating message: %s", createErr)
	}
	msg.Id = msgId
	return msg, nil
}

func seedChats() ([]domain.Chat, error) {
	msgs := []domain.Chat{
		{
			Sender:    "+6282387325960",
			Receiver:  "+6282387325961",
			Body:      "hello",
			CreatedAt: time.Now(),
		},
		{
			Sender:    "+6282387325980",
			Receiver:  "+6282387325981",
			Body:      "hello",
			CreatedAt: time.Now(),
		},
	}
	stmt, err := dbConn.Prepare(queryInsertChat)
	if err != nil {
		panic(err.Error())
	}
	for i := range msgs {
		_, createErr := stmt.Exec(msgs[i].Sender, msgs[i].Receiver, msgs[i].Body, msgs[i].CreatedAt)
		if createErr != nil {
			return nil, createErr
		}
	}

	getStmt, err := dbConn.Prepare(queryGetAllChats)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := getStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]domain.Chat, 0)

	for rows.Next() {
		var msg domain.Chat
		if getError := rows.Scan(&msg.Id, &msg.Sender, &msg.Receiver, &msg.Body, &msg.CreatedAt); getError != nil {
			return nil, err
		}
		results = append(results, msg)
	}
	return results, nil
}

package domain

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SemmiDev/lets-tests/utils"
	"log"
	"reflect"
	"testing"
	"time"
)

var sender = utils.RandomSender()
var receiver = utils.RandomReceiver()
var body = utils.RandomBody()
var createdAt = time.Now()

func TestMessageRepo_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewChatRepository(db)

	tests := []struct {
		name    string
		s       chatRepoInterface
		msgId   int64
		mock    func()
		want    *Chat
		wantErr bool
	}{
		{
			//When everything works as expected
			name:  "OK",
			s:     s,
			msgId: 1,
			mock: func() {
				//We added one row
				rows := sqlmock.NewRows([]string{"Id", "Sender", "Receiver", "Body", "CreatedAt"}).AddRow(1, sender, receiver, body, createdAt)
				mock.ExpectPrepare("SELECT (.+) FROM chats").ExpectQuery().WithArgs(1).WillReturnRows(rows)
			},
			want: &Chat{
				Id:        1,
				Sender:    sender,
				Receiver:  receiver,
				Body:      body,
				CreatedAt: createdAt,
			},
		},
		{
			//When the role tried to access is not found
			name:  "Not Found",
			s:     s,
			msgId: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"Id", "Sender", "Receiver", "Body", "CreatedAt"}) //observe that we didnt add any role here
				mock.ExpectPrepare("SELECT (.+) FROM chats").ExpectQuery().WithArgs(1).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			//When invalid statement is provided, ie the SQL syntax is wrong(in this case, we provided a wrong database)
			name:  "Invalid Prepare",
			s:     s,
			msgId: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"Id", "Sender", "Receiver", "Body", "CreatedAt"})
				mock.ExpectPrepare("SELECT (.+) FROM wrong_table").ExpectQuery().WithArgs(1).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Get(tt.msgId)
			log.Println(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
//func TestChatRepo_Create(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database", err)
//	}
//	defer db.Close()
//	s := NewChatRepository(db)
//	tm := time.Now()
//
//	tests := []struct {
//		name    string
//		s       chatRepoInterface
//		request *Chat
//		mock    func()
//		want    *Chat
//		wantErr bool
//	}{
//		{
//			name: "OK",
//			s:    s,
//			request: &Chat{
//				Sender:    sender,
//				Receiver:  receiver,
//				Body:      body,
//				CreatedAt: tm,
//			},
//			mock: func() {
//				mock.ExpectPrepare("INSERT INTO chats").ExpectExec().WithArgs("sender", "receiver", "body", tm).WillReturnResult(sqlmock.NewResult(1, 1))
//			},
//			want: &Chat{
//				Id:        1,
//				Sender:    sender,
//				Receiver:  receiver,
//				Body:      body,
//				CreatedAt: tm,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.mock()
//			got, err := tt.s.Create(tt.request)
//			if (err != nil) != tt.wantErr {
//				fmt.Println("this is the error message: ", err.Message())
//				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if err == nil && !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Create() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//When the right number of arguments are passed
//This test is just to improve coverage
func TestChatRepo_Initialize(t *testing.T) {
	dbdriver := "mysql"
	username := "root"
	password := ""
	host := "localhost"
	database := "chats"
	port := "5432"
	dbConnect := ChatRepo.Initialize(dbdriver, username, password, port, host, database)
	fmt.Println("this is the pool: ", dbConnect)
}

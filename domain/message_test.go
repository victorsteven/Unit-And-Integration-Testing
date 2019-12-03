package domain

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
	"time"
)

var created_at = time.Now()

func TestMessageRepo_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewMessageRepository(db)

	tests := []struct {
		name    string
		s       messageRepoInterface
		msgId   int64
		mock    func()
		want    *Message
		wantErr bool
	}{
		{
			//When everything works as expected
			name:  "OK",
			s:     s,
			msgId: 1,
			mock: func() {
				//We added one row
				rows := sqlmock.NewRows([]string{"Id", "Title", "Body", "CreatedAt"}).AddRow(1, "title", "body", created_at)
				mock.ExpectPrepare("SELECT (.+) FROM messages").ExpectQuery().WithArgs(1).WillReturnRows(rows)
			},
			want: &Message{
				Id:        1,
				Title:     "title",
				Body:      "body",
				CreatedAt: created_at,
			},
		},
		{
			//When the role tried to access is not found
			name:  "Not Found",
			s:     s,
			msgId: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"Id", "Title", "Body", "CreatedAt"}) //observe that we didnt add any role here
				mock.ExpectPrepare("SELECT (.+) FROM messages").ExpectQuery().WithArgs(1).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			//When invalid statement is provided, ie the SQL syntax is wrong(in this case, we provided a wrong database)
			name:  "Invalid Prepare",
			s:     s,
			msgId: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"Id", "Title", "Body", "CreatedAt"}).AddRow(1, "title", "body", created_at)
				mock.ExpectPrepare("SELECT (.+) FROM wrong_table").ExpectQuery().WithArgs(1).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Get(tt.msgId)
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

func TestMessageRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database", err)
	}
	defer db.Close()
	s := NewMessageRepository(db)
	tm := time.Now()

	tests := []struct {
		name    string
		s       messageRepoInterface
		request *Message
		mock    func()
		want    *Message
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			request: &Message{
				Title:     "title",
				Body:      "body",
				CreatedAt: tm,
			},
			mock: func() {
				mock.ExpectPrepare("INSERT INTO messages").ExpectExec().WithArgs("title", "body", tm).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: &Message{
				Id:        1,
				Title:     "title",
				Body:      "body",
				CreatedAt: tm,
			},
		},
		{
			name: "Empty title",
			s: s,
			request: &Message{
				Title:     "",
				Body:      "body",
				CreatedAt: tm,
			},
			mock: func(){
				mock.ExpectPrepare("INSERT INTO messages").ExpectExec().WithArgs("title", "body", tm).WillReturnError(errors.New("empty title"))
			},
			wantErr: true,
		},
		{
			name: "Empty body",
			s: s,
			request: &Message{
				Title:     "title",
				Body:      "",
				CreatedAt: tm,
			},
			mock: func(){
				mock.ExpectPrepare("INSERT INTO messages").ExpectExec().WithArgs("title", "body", tm).WillReturnError(errors.New("empty body"))
			},
			wantErr: true,
		},
		{
			name: "Invalid SQL query",
			s: s,
			request: &Message{
				Title:     "title",
				Body:      "body",
				CreatedAt: tm,
			},
			mock: func(){
				//Instead of using "INSERT", we used "INSETER"
				mock.ExpectPrepare("INSERT INTO wrong_table").ExpectExec().WithArgs("title", "body", tm).WillReturnError( errors.New("invalid sql query"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Create(tt.request)
			if (err != nil) != tt.wantErr {
				fmt.Println("this is the error message: ", err.Message())
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageRepo_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database", err)
	}
	defer db.Close()
	s := NewMessageRepository(db)

	tests := []struct {
		name    string
		s       messageRepoInterface
		request *Message
		mock    func()
		want    *Message
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			request: &Message{
				Id: 1,
				Title:     "update title",
				Body:      "update body",
			},
			mock: func() {
				mock.ExpectPrepare("UPDATE messages").ExpectExec().WithArgs("update title", "update body", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			want: &Message{
				Id:        1,
				Title:     "update title",
				Body:      "update body",
			},
		},
		{
			name: "Invalid SQL Query",
			s:    s,
			request: &Message{
				Id: 1,
				Title:     "update title",
				Body:      "update body",
			},
			mock: func() {
				mock.ExpectPrepare("UPDATER messages").ExpectExec().WithArgs("update title", "update body", 1).WillReturnError(errors.New("error in sql query statement"))
			},
			wantErr: true,
		},
		{
			name: "Invalid Query Id",
			s:    s,
			request: &Message{
				Id: 0,
				Title:     "update title",
				Body:      "update body",
			},
			mock: func() {
				mock.ExpectPrepare("UPDATE messages").ExpectExec().WithArgs("update title", "update body", 0).WillReturnError(errors.New("invalid update id"))
			},
			wantErr: true,
		},
		{
			name: "Empty Title",
			s:    s,
			request: &Message{
				Id: 1,
				Title:     "",
				Body:      "update body",
			},
			mock: func() {
				mock.ExpectPrepare("UPDATE messages").ExpectExec().WithArgs("", "update body", 1).WillReturnError(errors.New("Please enter a valid title"))
			},
			wantErr: true,
		},
		{
			name: "Empty Body",
			s:    s,
			request: &Message{
				Id: 1,
				Title:     "update title",
				Body:      "",
			},
			mock: func() {
				mock.ExpectPrepare("UPDATE messages").ExpectExec().WithArgs("update title", "", 1).WillReturnError(errors.New("Please enter a valid body"))
			},
			wantErr: true,
		},
		{
			name: "Update failed",
			s:    s,
			request: &Message{
				Id: 1,
				Title:     "update title",
				Body:      "update body",
			},
			mock: func() {
				mock.ExpectPrepare("UPDATE messages").ExpectExec().WithArgs("update title", "update body", 1).WillReturnResult(sqlmock.NewErrorResult(errors.New("Update failed")))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Update(tt.request)
			if (err != nil) != tt.wantErr {
				fmt.Println("this is the error message: ", err.Message())
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageRepo_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewMessageRepository(db)

	tests := []struct {
		name    string
		s       messageRepoInterface
		msgId   int64
		mock    func()
		want    []Message
		wantErr bool
	}{
		{
			//When everything works as expected
			name:  "OK",
			s:     s,
			mock: func() {
				//We added two rows
				rows := sqlmock.NewRows([]string{"Id", "Title", "Body", "CreatedAt"}).AddRow(1, "first title", "first body", created_at).AddRow(2, "second title", "second body", created_at)
				mock.ExpectPrepare("SELECT (.+) FROM messages").ExpectQuery().WillReturnRows(rows)
			},
			want: []Message{
				{
					Id:        1,
					Title:     "first title",
					Body:      "first body",
					CreatedAt: created_at,
				},
				{
					Id:        2,
					Title:     "second title",
					Body:      "second body",
					CreatedAt: created_at,
				},
			},
		},
		{
			name:  "Invalid SQL Syntax",
			s:     s,
			mock: func() {
				//We added two rows
				_ = sqlmock.NewRows([]string{"Id", "Title", "Body", "CreatedAt"}).AddRow(1, "first title", "first body", created_at).AddRow(2, "second title", "second body", created_at)
				//"SELECTS" is used instead of "SELECT"
				mock.ExpectPrepare("SELECTS (.+) FROM messages").ExpectQuery().WillReturnError(errors.New("Error when trying to prepare all messages"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageRepo_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewMessageRepository(db)

	tests := []struct {
		name    string
		s       messageRepoInterface
		msgId   int64
		mock    func()
		want    *Message
		wantErr bool
	}{
		{
			//When everything works as expected
			name:  "OK",
			s:     s,
			msgId: 1,
			mock: func() {
				mock.ExpectPrepare("DELETE FROM messages").ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:  "Invalid Id/Not Found Id",
			s:     s,
			msgId: 1,
			mock: func() {
				mock.ExpectPrepare("DELETE FROM messages").ExpectExec().WithArgs(100).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
		{
			name:  "Invalid SQL query",
			s:     s,
			msgId: 1,
			mock: func() {
				mock.ExpectPrepare("DELETE FROMSSSS messages").ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := tt.s.Delete(tt.msgId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

//When the right number of arguments are passed
//This test is just to improve coverage
func TestMessageRepo_Initialize(t *testing.T) {
	dbdriver :=  "mysql"
	username := "username"
	password := "password"
	host := "host"
	database := "database"
	port := "port"
	dbConnect := MessageRepo.Initialize(dbdriver, username, password, port, host, database)
	fmt.Println("this is the pool: ", dbConnect)
}
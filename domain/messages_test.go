package domain

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
	"time"
)

func TestMessageRepo_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewMessageRepository(db)
	created_at := time.Now()
	//var msgId int64

	tests := []struct {
		name    string
		s       MessageRepoInterface
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
			//if err != nil {
			//	fmt.Println("this is the generic error: ", err.Message())
			//	fmt.Println("this is the generic error status: ", err.Status())
			//	fmt.Println("this is the generic error error: ", err.Error())
			//}
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
	fmt.Println("this is the current time: ", tm)

	tests := []struct {
		name    string
		s       MessageRepoInterface
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

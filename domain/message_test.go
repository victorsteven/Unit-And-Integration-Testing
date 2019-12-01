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
		//{
		//	name: "lastInsertedId failed",
		//	s: s,
		//	request: &Message{
		//		Title:     "title",
		//		Body:      "body",
		//		CreatedAt: tm,
		//	},
		//	mock: func(){
		//		//mock.ExpectExec("INSERT INTO messages").WithArgs("title", "body", tm).WillReturnResult(sqlmock.NewErrorResult(errors.New("lastInsertId failed")))
		//		mock.ExpectPrepare("INSERT INTO messages").ExpectExec().WithArgs("title", "body", tm).WillReturnError( errors.New("lastInsertId failed"))
		//	},
		//	wantErr: true,
		//},
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
			if err != nil {
				fmt.Println("this is the generic error: ", err.Message())
				fmt.Println("this is the generic error status: ", err.Status())
				fmt.Println("this is the generic error error: ", err.Error())
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
	//tm := time.Now()

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
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				fmt.Println("this is the generic error: ", err.Message())
				fmt.Println("this is the generic error status: ", err.Status())
				fmt.Println("this is the generic error error: ", err.Error())
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestExecExpectations(t *testing.T) {
//	//db, err := sql.Open("mock", "")
//	//if err != nil {
//	//	t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
//	//}
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	result := NewResult(&Message{
//		Id:        1,
//		Title:     "title",
//		Body:      "body",
//		CreatedAt: tm,
//	})
//	//result, _ := &Message{
//	//	Id:        1,
//	//	Title:     "title",
//	//	Body:      "body",
//	//	CreatedAt: tm,
//	//},
//	mock.ExpectExec("^INSERT INTO articles").WithArgs("hello").WillReturnResult(result)
//	res, err := db.Exec("INSERT INTO articles (title) VALUES (?)", "hello")
//	if err != nil {
//		t.Errorf("error '%s' was not expected, while inserting a row", err)
//	}
//	id, err := res.LastInsertId()
//	if err != nil {
//		t.Errorf("error '%s' was not expected, while getting a last insert id", err)
//	}
//	affected, err := res.RowsAffected()
//	if err != nil {
//		t.Errorf("error '%s' was not expected, while getting affected rows", err)
//	}
//	if id != 1 {
//		t.Errorf("expected last insert id to be 1, but got %d instead", id)
//	}
//	if affected != 1 {
//		t.Errorf("expected affected rows to be 1, but got %d instead", affected)
//	}
//	if err = db.Close(); err != nil {
//		t.Errorf("error '%s' was not expected while closing the database", err)
//	}
//}

//func TestIssue4(t *testing.T) {
//	db, err := New()
//	if err != nil {
//		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	ExpectQuery("some sql query which will not be called").
//		WillReturnRows(NewRows([]string{"id"}))
//	err = db.Close()
//	if err == nil {
//		t.Errorf("Was expecting an error, since expected query was not matched")
//	}
//}
package domain

import (
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
	var msgId int64

	tests := []struct{
		name string
		s MessageRepoInterface
		msgId int64
		mock func()
		want *Message
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			msgId: msgId,
			mock: func() {
				rows := sqlmock.NewRows([]string{"Id", "Title", "Body", "CreatedAt"}).AddRow(1, "title", "body", created_at)
				mock.ExpectPrepare("SELECT (.+) FROM messages").ExpectQuery().WillReturnRows(rows)
			},
			want: &Message{
				Id:        1,
				Title:     "title",
				Body:      "body",
				CreatedAt: created_at,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Get(tt.msgId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error new = %v, wantErr %v", err.Message(), tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

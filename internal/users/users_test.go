package users

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	dummyID = 1
	dummyToken = "token"
)

func newMockService(ctrl *gomock.Controller) *Service {
	return NewService(NewMockauthenticator(ctrl), NewMockrepository(ctrl))
}

func dummyUser(user string, pwd string) *User {
	return &User{
		ID:       1,
		Username: user,
		Password: pwd,
	}
}

func TestService_Login(t *testing.T) {
	ass := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := newMockService(ctrl)

	type fields struct {
		user *User
	}

	type want struct {
		session *Session
		err      error
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{name: "ok",
			fields: fields{
				user: dummyUser("johnDoe", "abc123456"),
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).EXPECT().Get(gomock.Any()).Return(nil)
				service.auth.(*Mockauthenticator).EXPECT().Authenticate(gomock.Any()).Return(dummyToken, nil)
			},
			want: want{
				session: newSession(1, "token"),
				err: nil,
			},
		},
		{name: "error: user not found",
			fields: fields{
				user: dummyUser("johnDoe", "abc123456"),
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).EXPECT().Get(gomock.Any()).Return(errors.New("not found"))
			},
			want: want{
				session: nil,
				err: errors.New("not found"),
			},
		},
		{name: "error: can't create session",
			fields: fields{
				user: dummyUser("johnDoe", "abc123456"),
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).EXPECT().Get(gomock.Any()).Return(nil)
				service.auth.(*Mockauthenticator).EXPECT().Authenticate(gomock.Any()).Return("", errors.New("error creating token"))
			},
			want: want{
				session: nil,
				err: errors.New("error creating token"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectations(tt.fields)
			session, err := service.Login(tt.fields.user)
			ass.EqualValues(tt.want.session, session)
			ass.Equal(tt.want.err, err)
		})
	}
}

func TestService_Register(t *testing.T) {
	ass := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := newMockService(ctrl)

	type fields struct {
		user *User
	}

	type want struct {
		id 	 int
		err      error
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{name: "ok",
			fields: fields{
				user: dummyUser("johnDoe", "abc123456"),
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).EXPECT().Save(gomock.Any()).Return(dummyID, nil)
			},
			want: want{
				id: dummyID,
				err: nil,
			},
		},
		{name: "error",
			fields: fields{
				user: dummyUser("johnDoe", "abc123456"),
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).EXPECT().Save(gomock.Any()).Return(0, errors.New("error on registry"))
			},
			want: want{
				id: 0,
				err: errors.New("error on registry"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectations(tt.fields)
			id, err := service.Register(tt.fields.user)
			ass.EqualValues(tt.want.id, id)
			ass.Equal(tt.want.err, err)
		})
	}
}

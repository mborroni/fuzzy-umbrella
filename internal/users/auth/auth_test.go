package auth

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	dummyUser = "johnDoe"
	dummyToken = "token"
)

func newMockService(ctrl *gomock.Controller) *Service {
	return NewService(NewMockrepository(ctrl))
}

func TestService_Authenticate(t *testing.T) {
	ass := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := newMockService(ctrl)

	type fields struct {
		user string
	}

	type want struct {
		isEmpty bool
		err   error
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{name: "ok",
			fields: fields{
				user: dummyUser,
			},
			expectations: func(fields fields) {
				service.cache.(*Mockrepository).EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: want{
				isEmpty: false,
				err:     nil,
			},
		},
		{name: "error on cache save",
			fields: fields{
				user: dummyUser,
			},
			expectations: func(fields fields) {
				service.cache.(*Mockrepository).EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("error on save"))
			},
			want: want{
				isEmpty:   true,
				err:     errors.New("error on save"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectations(tt.fields)
			token, err := service.GenerateToken(tt.fields.user)
			ass.Equal(tt.want.isEmpty, token == "")
			ass.Equal(tt.want.err, err)
		})
	}
}

func TestService_Validate(t *testing.T) {
	ass := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := newMockService(ctrl)

	type fields struct {
		token string
		user string
	}

	type want struct {
		exists bool
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{name: "ok",
			fields: fields{
				token: dummyToken,
				user: dummyUser,
			},
			expectations: func(fields fields) {
				service.cache.(*Mockrepository).EXPECT().Get(gomock.Any()).Return(dummyUser, nil)
			},
			want: want{
				exists: true,
			},
		},
		{name: "not found",
			fields: fields{
				token: dummyToken,
				user: dummyUser,
			},
			expectations: func(fields fields) {
				service.cache.(*Mockrepository).EXPECT().Get(gomock.Any()).Return("", errors.New("not found"))
			},
			want: want{
				exists: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectations(tt.fields)
			exists := service.Validate(tt.fields.token, tt.fields.user)
			ass.Equal(tt.want.exists, exists)
		})
	}
}

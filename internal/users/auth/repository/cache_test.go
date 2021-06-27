package cache

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testTtl = 30
	dummyID = "token"
	dummyValue = "JohnDoe"
)

func newMockRepository(ctrl *gomock.Controller) *Repository {
	return NewRepository(testTtl, NewMockcache(ctrl))
}

func TestRepository_Get(t *testing.T) {
	ass := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := newMockRepository(ctrl)

	type fields struct {
		id string
	}

	type want struct {
		value string
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
				id: dummyID,
			},
			expectations: func(fields fields) {
				repository.cache.(*Mockcache).EXPECT().Get(gomock.Any()).Return([]byte(dummyValue), nil)
			},
			want: want{
				value: dummyValue,
				err:     nil,
			},
		},
		{name: "error",
			fields: fields{
				id: dummyID,
			},
			expectations: func(fields fields) {
				repository.cache.(*Mockcache).EXPECT().Get(gomock.Any()).Return(nil, errors.New("error on get"))
			},
			want: want{
				value: "",
				err:   errors.New("error on get"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectations(tt.fields)
			value, err := repository.Get(tt.fields.id)
			ass.Equal(tt.want.value, value)
			ass.Equal(tt.want.err, err)
		})
	}
}

func TestRepository_Save(t *testing.T) {
	ass := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := newMockRepository(ctrl)

	type fields struct {
		id string
		value string
	}

	type want struct {
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
				id: dummyID,
				value: dummyValue,
			},
			expectations: func(fields fields) {
				repository.cache.(*Mockcache).EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			want: want{
				err:     nil,
			},
		},
		{name: "error",
			fields: fields{
				id: dummyID,
				value: dummyValue,
			},
			expectations: func(fields fields) {
				repository.cache.(*Mockcache).EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error on save"))
			},
			want: want{
				err:   errors.New("error on save"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectations(tt.fields)
			err := repository.Save(tt.fields.id, tt.fields.value)
			ass.Equal(tt.want.err, err)
		})
	}
}
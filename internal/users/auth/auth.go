package auth

import (
	"crypto/md5"
	"encoding/hex"

	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

//go:generate mockgen -source=auth.go -destination=auth_mock.go -package=auth repository

const (
	tokenKeyFormat = "%s_%s"
)

type repository interface {
	Get(string) (string, error)
	Save(string, interface{}) error
}

type Service struct {
	cache repository
}

func NewService(r repository) *Service {
	return &Service{
		cache: r,
	}
}

func (auth *Service) GenerateToken(user string) (string, error) {
	token, err := getHash(user)
	if err != nil {
		return "", err
	}
	if err := auth.cache.Save(token, user); err != nil {
		return "", err
	}
	return token, nil
}

func (auth *Service) Validate(token string, user string) bool {
	storedUser, err := auth.cache.Get(token)
	if err != nil {
		return false
	}
	return user == storedUser
}

func getHash(user string) (string, error) {
	key := fmt.Sprintf(tokenKeyFormat, user, time.Now().String())
	hash, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
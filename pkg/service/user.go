package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/coci/chitchat/pkg/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserService struct {
	UserRepository       repository.UserRepository
	UserAPIKeyRepository repository.UserApiKeyRepository
}

func ParseDelimitedFields(data []byte, expectedParts int) ([]string, error) {
	parts := bytes.SplitN(data, []byte(","), expectedParts)

	if len(parts) != expectedParts {
		return nil, fmt.Errorf("invalid format: expected %d comma-separated fields", expectedParts)
	}

	result := make([]string, expectedParts)
	for i, p := range parts {
		field := string(bytes.TrimSpace(p))
		if field == "" {
			return nil, fmt.Errorf("field %d is empty", i+1)
		}
		result[i] = field
	}

	return result, nil
}

func (u UserService) CreateUser(data []byte) string {

	fields, err := ParseDelimitedFields(data, 2)

	user, err := u.UserRepository.StoreUser(
		fields[0],
		u.hashPassword(fields[1]),
	)
	if err != nil {
		log.Println(err)
	}
	apikey := u.createUserApiKey(
		user.ID,
	)

	return apikey
}

func (u UserService) LoginUser(data []byte) string {
	fields, err := ParseDelimitedFields(data, 2)

	user, err := u.UserRepository.GetUser(
		fields[0],
		u.hashPassword(fields[1]),
	)
	if err != nil {
		log.Println(err)
	}

	apikey := u.createUserApiKey(
		user.ID,
	)

	return apikey
}

func (u UserService) createUserApiKey(userID uint) string {
	buff := make([]byte, 12)
	apikey := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(buff)

	u.UserAPIKeyRepository.StoreUserApiKey(
		userID,
		apikey,
	)

	return apikey
}

func (u UserService) hashPassword(password string) string {
	HashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(HashedPassword)
}

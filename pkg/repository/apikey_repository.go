package repository

import (
	"github.com/redis/go-redis/v9"
)

type UserApiKeyRepository struct {
	client *redis.Client
}

func NewUserApiKeyRepository(client *redis.Client) *UserApiKeyRepository {
	return &UserApiKeyRepository{client: client}
}

func (r *UserApiKeyRepository) StoreUserApiKey(userID uint, apiKey string) {
}

func (r *UserApiKeyRepository) GetUserID(apiKey string) uint {
	return 1
}

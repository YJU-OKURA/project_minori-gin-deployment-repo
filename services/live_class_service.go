package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"github.com/go-redis/redis/v8"
	"time"
)

type LiveClassService interface {
	GetScreenShareInfo(ctx context.Context, uid uint, cid uint) (interface{}, error)
	SaveScreenShareInfo(ctx context.Context, cid uint, info map[string]interface{}) error
}

type liveClassServiceImpl struct {
	classUserRepository repositories.ClassUserRepository
	redisClient         *redis.Client
}

func NewLiveClassService(classUserRepo repositories.ClassUserRepository, redisClient *redis.Client) LiveClassService {
	return &liveClassServiceImpl{
		classUserRepository: classUserRepo,
		redisClient:         redisClient,
	}
}

func (service *liveClassServiceImpl) GetScreenShareInfo(ctx context.Context, uid uint, cid uint) (interface{}, error) {
	isMember, err := service.classUserRepository.IsMember(uid, cid)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("user is not a member of the class")
	}

	// Retrieve screen share data from Redis
	data, err := service.redisClient.Get(ctx, makeRedisKey(cid)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("no active screen sharing session")
		}
		return nil, err
	}

	var info map[string]interface{}
	err = json.Unmarshal([]byte(data), &info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (service *liveClassServiceImpl) SaveScreenShareInfo(ctx context.Context, cid uint, info map[string]interface{}) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	// Save to Redis with an expiration (e.g., 2 hours)
	return service.redisClient.Set(ctx, makeRedisKey(cid), data, 2*time.Hour).Err()
}

func makeRedisKey(cid uint) string {
	return fmt.Sprintf("screen_share:%d", cid)
}

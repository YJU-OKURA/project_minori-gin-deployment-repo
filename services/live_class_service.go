package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"github.com/go-redis/redis/v8"
)

type LiveClassService interface {
	GetScreenShareInfo(ctx context.Context, cid uint) (interface{}, error)
	SaveScreenShareInfo(ctx context.Context, cid uint, info map[string]interface{}) error
	StartStreamingSession(cid uint) (string, error)
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

func (service *liveClassServiceImpl) GetScreenShareInfo(ctx context.Context, cid uint) (interface{}, error) {
	data, err := service.redisClient.Get(ctx, makeRedisKey(cid)).Result()
	if err != nil {
		return nil, err
	}
	var info map[string]interface{}
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, err
	}
	return info, nil
}

func (service *liveClassServiceImpl) SaveScreenShareInfo(ctx context.Context, cid uint, info map[string]interface{}) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	// Save to Redis with expiration
	return service.redisClient.Set(ctx, makeRedisKey(cid), data, 2*time.Hour).Err()
}

func makeRedisKey(cid uint) string {
	return fmt.Sprintf("screen_share:%d", cid)
}

func (service *liveClassServiceImpl) StartStreamingSession(cid uint) (string, error) {
	// API 호출 로직 구현 (예시: HTTP 요청)
	// 예를 들어, 스트리밍 서버로 POST 요청을 보내고 응답에서 URL을 추출
	response, err := http.Post(fmt.Sprintf("https://minoriedu.com/start/%d", cid), "application/json", nil)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	streamURL, ok := result["streamURL"].(string)
	if !ok {
		return "", errors.New("invalid response from streaming service")
	}

	return streamURL, nil
}

func (service *liveClassServiceImpl) MonitorStream(cid uint) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		response, err := http.Get(fmt.Sprintf("https://minoriedu.com/stream/status/%d", cid))
		if err != nil {
			log.Println("Failed to check stream status:", err)
			continue
		}

		var status map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&status); err != nil {
			log.Println("Error decoding status response:", err)
			continue
		}

		if active, ok := status["active"].(bool); ok && !active {
			log.Println("Stream has stopped unexpectedly, attempting to restart...")
			service.StartStreamingSession(cid)
		}
	}
}

func (service *liveClassServiceImpl) StopStreamingSession(cid uint) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://minoriedu.com/stop/%d", cid), nil)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to stop streaming session with status: %s", response.Status)
	}

	return nil
}

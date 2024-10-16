package tests

import (
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/controllers"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockClassUserService struct {
	mock.Mock
}

func (m *MockClassUserService) GetClassMembers(cid uint, roleNames ...string) ([]dto.ClassMemberDTO, error) {
	args := m.Called(cid, roleNames)
	return args.Get(0).([]dto.ClassMemberDTO), args.Error(1)
}

func (m *MockClassUserService) GetClassUserInfo(uid uint, cid uint) (dto.ClassMemberDTO, error) {
	args := m.Called(uid, cid)
	return args.Get(0).(dto.ClassMemberDTO), args.Error(1)
}

func (m *MockClassUserService) GetUserClasses(uid uint, page int, limit int) ([]dto.UserClassInfoDTO, error) {
	args := m.Called(uid, page, limit)
	return args.Get(0).([]dto.UserClassInfoDTO), args.Error(1)
}

func (m *MockClassUserService) GetRole(uid uint, cid uint) (string, error) {
	args := m.Called(uid, cid)
	return args.String(0), args.Error(1)
}

func (m *MockClassUserService) GetFavoriteClasses(uid uint, page int, limit int) ([]dto.UserClassInfoDTO, error) {
	args := m.Called(uid, page, limit)
	return args.Get(0).([]dto.UserClassInfoDTO), args.Error(1)
}

func (m *MockClassUserService) GetUserClassesByRole(uid uint, roleName string, page int, limit int) ([]dto.UserClassInfoDTO, error) {
	args := m.Called(uid, roleName, page, limit)
	return args.Get(0).([]dto.UserClassInfoDTO), args.Error(1)
}

func (m *MockClassUserService) AssignRole(uid uint, cid uint, roleName string) error {
	args := m.Called(uid, cid, roleName)
	return args.Error(0)
}

func (m *MockClassUserService) UpdateUserName(uid uint, cid uint, newName string) error {
	args := m.Called(uid, cid, newName)
	return args.Error(0)
}

func (m *MockClassUserService) ToggleFavorite(uid uint, cid uint) error {
	args := m.Called(uid, cid)
	return args.Error(0)
}

func (m *MockClassUserService) RemoveUserFromClass(uid uint, cid uint) error {
	args := m.Called(uid, cid)
	return args.Error(0)
}

func (m *MockClassUserService) SearchUserClassesByName(uid uint, name string) ([]dto.UserClassInfoDTO, error) {
	args := m.Called(uid, name)
	return args.Get(0).([]dto.UserClassInfoDTO), args.Error(1)
}

func TestGetUserInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockClassUserService)
	controller := controllers.NewClassUserController(mockService)

	router := gin.Default()
	router.GET("/cu/:uid/:cid/info", controller.GetUserClassUserInfo)

	t.Run("Success", func(t *testing.T) {
		uid := uint(1)
		cid := uint(1)
		expectedResponse := dto.ClassMemberDTO{
			Uid:      uid,
			Nickname: "testuser",
			Role:     "USER",
			Image:    "testimage.png",
		}

		mockService.On("GetClassUserInfo", uid, cid).Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/cu/"+strconv.Itoa(int(uid))+"/"+strconv.Itoa(int(cid))+"/info", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid UID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/cu/invalid/1/info", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

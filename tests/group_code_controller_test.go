package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/controllers"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockGroupCodeService はGroupCodeServiceのモックです。
type MockGroupCodeService struct {
	mock.Mock
}

// TestCheckSecretExists はシークレットの存在を確認するテストです。
func setUpTestRouter() (*gin.Engine, *MockGroupCodeService, *controllers.GroupCodeController) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockGroupCodeService)
	controller := controllers.NewGroupCodeController(mockService)
	r := gin.Default()
	r.POST("/gc/checkSecretExists", controller.CheckSecretExists)
	r.POST("/gc/verifyGroupCode", controller.VerifyGroupCode)
	return r, mockService, controller
}

// CheckSecretExists は指定されたグループコードにシークレットがあるかどうかをチェックします。
func (m *MockGroupCodeService) CheckSecretExists(code string) (bool, error) {
	args := m.Called(code)
	return args.Bool(0), args.Error(1)
}

// VerifyGroupCode はグループコードと、該当する場合はそのシークレットを確認します。
func (m *MockGroupCodeService) VerifyGroupCode(code, secret string) (bool, error) {
	args := m.Called(code, secret)
	return args.Bool(0), args.Error(1)
}

// performRequest はリクエストを実行し、結果を返します。
func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// TestCheckSecretExists はシークレットの存在を確認するテストです。
func TestCheckSecretExists(t *testing.T) {
	r, mockService, _ := setUpTestRouter()

	// シークレットが存在するケースのテスト
	t.Run("CheckSecretExists - Secret Exists", func(t *testing.T) {
		mockService.On("CheckSecretExists", "validCode").Return(true, nil)
		body, _ := json.Marshal(dto.GroupCodeCheckRequest{Code: "validCode"})
		req, _ := http.NewRequest(http.MethodPost, "/gc/checkSecretExists", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	// シークレットが存在しないケースのテスト
	t.Run("CheckSecretExists - Secret Not Found", func(t *testing.T) {
		mockService.On("CheckSecretExists", "invalidCode").Return(false, gorm.ErrRecordNotFound)
		body, _ := json.Marshal(dto.GroupCodeCheckRequest{Code: "invalidCode"})
		resp := performRequest(r, http.MethodPost, "/gc/checkSecretExists", body)

		assert.Equal(t, constants.StatusNotFound, resp.Code)
		assert.Contains(t, resp.Body.String(), constants.CodeNotFound)
		mockService.AssertExpectations(t)
	})

	// 無効なリクエストのテスト
	t.Run("CheckSecretExists - Invalid Request", func(t *testing.T) {
		resp := performRequest(r, http.MethodPost, "/gc/checkSecretExists", []byte("{}")) // 空のリクエスト

		assert.Equal(t, constants.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), constants.InvalidRequest)
	})
}

// TestVerifyGroupCode はグループコードとシークレットを検証するテストです。
func TestVerifyGroupCode(t *testing.T) {
	r, mockService, _ := setUpTestRouter()

	// 正常なリクエスト
	t.Run("Success", func(t *testing.T) {
		mockService.On("VerifyGroupCode", "validCode", "validSecret").Return(true, nil)
		body, _ := json.Marshal(dto.GroupCodeRequest{Code: "validCode", Secret: "validSecret"})
		req, _ := http.NewRequest(http.MethodPost, "/gc/verifyGroupCode", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, constants.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), constants.GroupCodeVerified)
		mockService.AssertExpectations(t)
	})

	// シークレットが存在しないケースのテスト
	t.Run("Secret Not Exists", func(t *testing.T) {
		mockService.On("CheckSecretExists", "codeWithoutSecret").Return(false, nil)
		body, _ := json.Marshal(dto.GroupCodeCheckRequest{Code: "codeWithoutSecret"})
		resp := performRequest(r, http.MethodPost, "/gc/checkSecretExists", body)

		assert.Equal(t, constants.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "false") // "secretExists": false
		mockService.AssertExpectations(t)
	})

	// パスワードの不一致
	t.Run("VerifyGroupCode - Secret Mismatch", func(t *testing.T) {
		mockService.On("VerifyGroupCode", "validCode", "invalidSecret").Return(false, nil)
		body, _ := json.Marshal(dto.GroupCodeRequest{Code: "validCode", Secret: "invalidSecret"})
		req, _ := http.NewRequest(http.MethodPost, "/gc/verifyGroupCode", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		mockService.AssertExpectations(t)
	})

	// コードが見つからないケースのテスト
	t.Run("VerifyGroupCode - Code Not Found", func(t *testing.T) {
		mockService.On("VerifyGroupCode", "invalidCode", "anySecret").Return(false, gorm.ErrRecordNotFound)
		body, _ := json.Marshal(dto.GroupCodeRequest{Code: "invalidCode", Secret: "anySecret"})
		req, _ := http.NewRequest(http.MethodPost, "/gc/verifyGroupCode", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, constants.StatusNotFound, resp.Code)
		assert.Contains(t, resp.Body.String(), constants.CodeNotFound)
		mockService.AssertExpectations(t)
	})

	// 無効なリクエストのテスト
	t.Run("VerifyGroupCode - Invalid Request", func(t *testing.T) {
		resp := performRequest(r, http.MethodPost, "/gc/verifyGroupCode", []byte("{}")) // 空のリクエスト

		assert.Equal(t, constants.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), constants.InvalidRequest)
	})

	// サーバーエラーのテスト
	t.Run("Internal Server Error", func(t *testing.T) {
		mockService.On("VerifyGroupCode", "code", "secret").Return(false, errors.New("internal error"))
		body, _ := json.Marshal(dto.GroupCodeRequest{Code: "code", Secret: "secret"})
		resp := performRequest(r, http.MethodPost, "/gc/verifyGroupCode", body)

		assert.Equal(t, constants.StatusInternalServerError, resp.Code)
		assert.Contains(t, resp.Body.String(), constants.InternalServerError)
		mockService.AssertExpectations(t)
	})
}

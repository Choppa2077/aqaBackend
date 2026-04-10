package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/PrimeraAizen/e-comm/internal/domain"
	"github.com/PrimeraAizen/e-comm/internal/service"
	"github.com/PrimeraAizen/e-comm/pkg/logger"
)

// MockAuthService implements service.AuthService for integration testing
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req *domain.User) (*domain.Token, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Token), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.Token, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Token), args.Error(1)
}

func (m *MockAuthService) ValidateToken(tokenString string) (*domain.TokenClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Token), args.Error(1)
}

// setupRouter creates a test gin router with auth routes wired to a mock service
func setupRouter(authSvc service.AuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	appLogger, _ := logger.New(&logger.Config{Level: "error", Format: "text", Output: "stdout"})

	svc := &service.Service{AuthService: authSvc}
	handler := NewHandler(svc, appLogger)

	r := gin.New()
	api := r.Group("/api")
	handler.Init(api)
	return r
}

// --- Integration tests ---

func TestRegisterHandler_Success(t *testing.T) {
	mockAuth := new(MockAuthService)
	router := setupRouter(mockAuth)

	token := &domain.Token{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    900,
	}
	mockAuth.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(token, nil)

	body := `{"email":"newuser@example.com","password":"password123","password_confirm":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "test-access-token", resp["access_token"])
	mockAuth.AssertExpectations(t)
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	mockAuth := new(MockAuthService)
	router := setupRouter(mockAuth)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockAuth.AssertNotCalled(t, "Register")
}

func TestLoginHandler_Success(t *testing.T) {
	mockAuth := new(MockAuthService)
	router := setupRouter(mockAuth)

	token := &domain.Token{
		AccessToken:  "access-xyz",
		RefreshToken: "refresh-xyz",
		TokenType:    "Bearer",
		ExpiresIn:    900,
	}
	mockAuth.On("Login", mock.Anything, mock.AnythingOfType("*domain.LoginRequest")).Return(token, nil)

	body := `{"email":"admin@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "access-xyz", resp["access_token"])
	assert.Equal(t, "refresh-xyz", resp["refresh_token"])
	mockAuth.AssertExpectations(t)
}

func TestLoginHandler_WrongCredentials(t *testing.T) {
	mockAuth := new(MockAuthService)
	router := setupRouter(mockAuth)

	mockAuth.On("Login", mock.Anything, mock.AnythingOfType("*domain.LoginRequest")).Return(nil, domain.ErrInvalidCredentials)

	body := `{"email":"user@example.com","password":"wrongpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid email or password", resp["error"])
	mockAuth.AssertExpectations(t)
}

func TestLoginHandler_MissingFields(t *testing.T) {
	mockAuth := new(MockAuthService)
	router := setupRouter(mockAuth)

	// Missing password field — validator should reject
	body := `{"email":"user@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockAuth.AssertNotCalled(t, "Login")
}

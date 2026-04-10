package service

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/PrimeraAizen/e-comm/config"
	"github.com/PrimeraAizen/e-comm/internal/domain"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// testConfig returns a config suitable for unit tests
func testConfig() *config.Config {
	return &config.Config{
		JWT: config.JWT{
			Secret:               "test-secret-key-for-unit-tests",
			AccessTokenDuration:  "15m",
			RefreshTokenDuration: "168h",
		},
	}
}

func newTestAuthService(userRepo *MockUserRepository) AuthService {
	svc, err := NewAuthService(userRepo, testConfig())
	if err != nil {
		panic("failed to create auth service: " + err.Error())
	}
	return svc
}

// --- Register tests ---

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	user := &domain.User{Email: "newuser@example.com", PasswordHash: "password123"}

	mockRepo.On("GetByEmail", ctx, "newuser@example.com").Return(nil, domain.ErrNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	token, err := svc.Register(ctx, user)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.AccessToken)
	assert.NotEmpty(t, token.RefreshToken)
	assert.Equal(t, "Bearer", token.TokenType)
	mockRepo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	existingUser := &domain.User{ID: 1, Email: "existing@example.com", Status: "active"}
	mockRepo.On("GetByEmail", ctx, "existing@example.com").Return(existingUser, nil)

	user := &domain.User{Email: "existing@example.com", PasswordHash: "password123"}
	token, err := svc.Register(ctx, user)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrAlreadyExists, err)
	assert.Nil(t, token)
	mockRepo.AssertExpectations(t)
}

// --- Login tests ---

func TestLogin_ValidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	existingUser := &domain.User{
		ID:           1,
		Email:        "user@example.com",
		PasswordHash: string(hash),
		Status:       "active",
		CreatedAt:    time.Now(),
	}

	mockRepo.On("GetByEmail", ctx, "user@example.com").Return(existingUser, nil)
	mockRepo.On("UpdateLastLogin", ctx, 1).Return(nil)

	req := &domain.LoginRequest{Email: "user@example.com", Password: "password123"}
	token, err := svc.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.AccessToken)
	assert.NotEmpty(t, token.RefreshToken)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.MinCost)
	existingUser := &domain.User{
		ID:           1,
		Email:        "user@example.com",
		PasswordHash: string(hash),
		Status:       "active",
	}

	mockRepo.On("GetByEmail", ctx, "user@example.com").Return(existingUser, nil)

	req := &domain.LoginRequest{Email: "user@example.com", Password: "wrongpassword"}
	token, err := svc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidCredentials, err)
	assert.Nil(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByEmail", ctx, "notexist@example.com").Return(nil, domain.ErrNotFound)

	req := &domain.LoginRequest{Email: "notexist@example.com", Password: "password123"}
	token, err := svc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidCredentials, err)
	assert.Nil(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InactiveUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	inactiveUser := &domain.User{
		ID:           2,
		Email:        "inactive@example.com",
		PasswordHash: string(hash),
		Status:       "suspended",
	}

	mockRepo.On("GetByEmail", ctx, "inactive@example.com").Return(inactiveUser, nil)

	req := &domain.LoginRequest{Email: "inactive@example.com", Password: "password123"}
	token, err := svc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserInactive, err)
	assert.Nil(t, token)
	mockRepo.AssertExpectations(t)
}

// --- ValidateToken tests ---

func TestValidateToken_Valid(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	// First register to get a real token
	mockRepo.On("GetByEmail", ctx, "token@example.com").Return(nil, domain.ErrNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	user := &domain.User{ID: 1, Email: "token@example.com"}
	token, err := svc.Register(ctx, user)
	assert.NoError(t, err)

	// Now validate the token
	claims, err := svc.ValidateToken(token.AccessToken)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "token@example.com", claims.Email)
}

func TestValidateToken_Invalid(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)

	claims, err := svc.ValidateToken("this.is.not.a.valid.token")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)

	// Token signed with a different secret
	wrongSecretToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjk5OTk5OTk5OTksInVzZXJfaWQiOiIxIn0.wrongsignature"

	claims, err := svc.ValidateToken(wrongSecretToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

// --- Edge case & new tests for Midterm ---

func TestLogin_EmptyEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByEmail", ctx, "").Return(nil, domain.ErrNotFound)

	req := &domain.LoginRequest{Email: "", Password: "password123"}
	token, err := svc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidCredentials, err)
	assert.Nil(t, token)
	mockRepo.AssertExpectations(t)
}

func TestRegister_EmptyPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	// Empty password hash — bcrypt should still create a token, but
	// the service itself does not validate password strength (that is done in DTO layer).
	// We test that if Create fails, Register propagates the error.
	mockRepo.On("GetByEmail", ctx, "edge@example.com").Return(nil, domain.ErrNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(fmt.Errorf("password hash required"))

	user := &domain.User{Email: "edge@example.com", PasswordHash: ""}
	token, err := svc.Register(ctx, user)

	assert.Error(t, err)
	assert.Nil(t, token)
	mockRepo.AssertExpectations(t)
}

func TestValidateToken_Expired(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)

	// Generate a token with -1 second duration (already expired)
	cfg := testConfig()
	cfg.JWT.AccessTokenDuration = "-1s"
	expiredSvc, err := NewAuthService(mockRepo, cfg)
	assert.NoError(t, err)

	ctx := context.Background()
	mockRepo.On("GetByEmail", ctx, "expired@example.com").Return(nil, domain.ErrNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	user := &domain.User{Email: "expired@example.com"}
	token, err := expiredSvc.Register(ctx, user)
	assert.NoError(t, err)

	// Now validate the already-expired access token using the normal service
	claims, err := svc.ValidateToken(token.AccessToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestRefreshToken_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	// Register to get tokens
	mockRepo.On("GetByEmail", ctx, "refresh@example.com").Return(nil, domain.ErrNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	user := &domain.User{ID: 10, Email: "refresh@example.com"}
	token, err := svc.Register(ctx, user)
	assert.NoError(t, err)

	// Refresh: service validates refresh token, then calls GetByID
	activeUser := &domain.User{ID: 10, Email: "refresh@example.com", Status: "active"}
	mockRepo.On("GetByID", ctx, 10).Return(activeUser, nil)

	newToken, err := svc.RefreshToken(ctx, token.RefreshToken)

	assert.NoError(t, err)
	assert.NotNil(t, newToken)
	assert.NotEmpty(t, newToken.AccessToken)
	mockRepo.AssertExpectations(t)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	newToken, err := svc.RefreshToken(ctx, "this.is.invalid.refresh.token")

	assert.Error(t, err)
	assert.Nil(t, newToken)
}

func TestRegister_Concurrent(t *testing.T) {
	// Concurrency test: 10 goroutines try to register the same email simultaneously.
	// Goroutine 0 sees slot free and succeeds; goroutines 1-9 see existing user and get ErrAlreadyExists.
	const goroutines = 10
	results := make([]error, goroutines)
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			mockRepo := new(MockUserRepository)
			svc := newTestAuthService(mockRepo)
			ctx := context.Background()

			if idx == 0 {
				mockRepo.On("GetByEmail", ctx, "concurrent@example.com").Return(nil, domain.ErrNotFound)
				mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
			} else {
				existing := &domain.User{ID: idx, Email: "concurrent@example.com", Status: "active"}
				mockRepo.On("GetByEmail", ctx, "concurrent@example.com").Return(existing, nil)
			}

			user := &domain.User{Email: "concurrent@example.com", PasswordHash: "hash"}
			_, err := svc.Register(ctx, user)
			results[idx] = err
		}(i)
	}

	wg.Wait()

	successCount := 0
	duplicateCount := 0
	for _, err := range results {
		if err == nil {
			successCount++
		} else if err == domain.ErrAlreadyExists {
			duplicateCount++
		}
	}

	assert.Equal(t, 1, successCount)
	assert.Equal(t, goroutines-1, duplicateCount)
}

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/PrimeraAizen/e-comm/internal/domain"
)

// MockInteractionRepository is a mock implementation of InteractionRepository
type MockInteractionRepository struct {
	mock.Mock
}

func (m *MockInteractionRepository) RecordView(ctx context.Context, userID, productID int) error {
	args := m.Called(ctx, userID, productID)
	return args.Error(0)
}

func (m *MockInteractionRepository) GetUserViews(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ProductInteraction), args.Error(1)
}

func (m *MockInteractionRepository) HasViewed(ctx context.Context, userID, productID int) (bool, error) {
	args := m.Called(ctx, userID, productID)
	return args.Bool(0), args.Error(1)
}

func (m *MockInteractionRepository) RecordLike(ctx context.Context, userID, productID int) error {
	args := m.Called(ctx, userID, productID)
	return args.Error(0)
}

func (m *MockInteractionRepository) RemoveLike(ctx context.Context, userID, productID int) error {
	args := m.Called(ctx, userID, productID)
	return args.Error(0)
}

func (m *MockInteractionRepository) GetUserLikes(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ProductInteraction), args.Error(1)
}

func (m *MockInteractionRepository) HasLiked(ctx context.Context, userID, productID int) (bool, error) {
	args := m.Called(ctx, userID, productID)
	return args.Bool(0), args.Error(1)
}

func (m *MockInteractionRepository) RecordPurchase(ctx context.Context, userID, productID int, quantity int, price float64) error {
	args := m.Called(ctx, userID, productID, quantity, price)
	return args.Error(0)
}

func (m *MockInteractionRepository) GetUserPurchases(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ProductInteraction), args.Error(1)
}

func (m *MockInteractionRepository) HasPurchased(ctx context.Context, userID, productID int) (bool, error) {
	args := m.Called(ctx, userID, productID)
	return args.Bool(0), args.Error(1)
}

func (m *MockInteractionRepository) GetUserInteractionSummary(ctx context.Context, userID int) (*domain.UserInteractionSummary, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserInteractionSummary), args.Error(1)
}

func (m *MockInteractionRepository) GetAllUserViews(ctx context.Context) ([]domain.UserProductView, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserProductView), args.Error(1)
}

func (m *MockInteractionRepository) GetAllUserLikes(ctx context.Context) ([]domain.UserProductLike, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserProductLike), args.Error(1)
}

func (m *MockInteractionRepository) GetAllUserPurchases(ctx context.Context) ([]domain.UserProductPurchase, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserProductPurchase), args.Error(1)
}

// helper to create interaction service with mocks
func newTestInteractionService(interRepo *MockInteractionRepository, prodRepo *MockProductRepository) InteractionService {
	return NewInteractionService(interRepo, prodRepo)
}

// --- PurchaseProduct tests ---

func TestPurchaseProduct_Success(t *testing.T) {
	interRepo := new(MockInteractionRepository)
	prodRepo := new(MockProductRepository)
	svc := newTestInteractionService(interRepo, prodRepo)
	ctx := context.Background()

	product := &domain.Product{ID: 1, Name: "iPhone", Price: 999.99, Stock: 10}
	updatedProduct := &domain.Product{ID: 1, Name: "iPhone", Price: 999.99, Stock: 8}

	prodRepo.On("GetByID", ctx, 1).Return(product, nil)
	interRepo.On("RecordPurchase", ctx, 5, 1, 2, 999.99).Return(nil)
	prodRepo.On("Update", ctx, updatedProduct).Return(nil)

	err := svc.PurchaseProduct(ctx, 5, 1, 2)

	assert.NoError(t, err)
	prodRepo.AssertExpectations(t)
	interRepo.AssertExpectations(t)
}

func TestPurchaseProduct_InsufficientStock(t *testing.T) {
	interRepo := new(MockInteractionRepository)
	prodRepo := new(MockProductRepository)
	svc := newTestInteractionService(interRepo, prodRepo)
	ctx := context.Background()

	product := &domain.Product{ID: 1, Name: "iPhone", Price: 999.99, Stock: 1}
	prodRepo.On("GetByID", ctx, 1).Return(product, nil)

	err := svc.PurchaseProduct(ctx, 5, 1, 5) // requesting 5, only 1 in stock

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient stock")
	interRepo.AssertNotCalled(t, "RecordPurchase")
}

func TestPurchaseProduct_ZeroQuantity(t *testing.T) {
	interRepo := new(MockInteractionRepository)
	prodRepo := new(MockProductRepository)
	svc := newTestInteractionService(interRepo, prodRepo)
	ctx := context.Background()

	err := svc.PurchaseProduct(ctx, 5, 1, 0) // quantity 0 should fail

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quantity must be greater than 0")
	prodRepo.AssertNotCalled(t, "GetByID")
}

func TestPurchaseProduct_ProductNotFound(t *testing.T) {
	interRepo := new(MockInteractionRepository)
	prodRepo := new(MockProductRepository)
	svc := newTestInteractionService(interRepo, prodRepo)
	ctx := context.Background()

	prodRepo.On("GetByID", ctx, 999).Return(nil, domain.ErrNotFound)

	err := svc.PurchaseProduct(ctx, 5, 999, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")
	interRepo.AssertNotCalled(t, "RecordPurchase")
}

// --- LikeProduct tests ---

func TestLikeProduct_Success(t *testing.T) {
	interRepo := new(MockInteractionRepository)
	prodRepo := new(MockProductRepository)
	svc := newTestInteractionService(interRepo, prodRepo)
	ctx := context.Background()

	product := &domain.Product{ID: 1, Name: "iPhone", Price: 999.99}
	prodRepo.On("GetByID", ctx, 1).Return(product, nil)
	interRepo.On("RecordLike", ctx, 5, 1).Return(nil)

	err := svc.LikeProduct(ctx, 5, 1)

	assert.NoError(t, err)
	prodRepo.AssertExpectations(t)
	interRepo.AssertExpectations(t)
}

func TestLikeProduct_ProductNotFound(t *testing.T) {
	interRepo := new(MockInteractionRepository)
	prodRepo := new(MockProductRepository)
	svc := newTestInteractionService(interRepo, prodRepo)
	ctx := context.Background()

	prodRepo.On("GetByID", ctx, 999).Return(nil, domain.ErrNotFound)

	err := svc.LikeProduct(ctx, 5, 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")
	interRepo.AssertNotCalled(t, "RecordLike")
}

// --- GetUserPurchaseHistory tests ---

func TestGetUserPurchaseHistory_Success(t *testing.T) {
	interRepo := new(MockInteractionRepository)
	prodRepo := new(MockProductRepository)
	svc := newTestInteractionService(interRepo, prodRepo)
	ctx := context.Background()

	expected := []domain.ProductInteraction{
		{ProductID: 1, ProductName: "iPhone", Price: 999.99, InteractedAt: time.Now()},
		{ProductID: 2, ProductName: "MacBook", Price: 1999.99, InteractedAt: time.Now()},
	}
	interRepo.On("GetUserPurchases", ctx, 5, 50).Return(expected, nil)

	history, err := svc.GetUserPurchaseHistory(ctx, 5, 0) // 0 → defaults to 50

	assert.NoError(t, err)
	assert.Len(t, history, 2)
	interRepo.AssertExpectations(t)
}

func TestGetUserPurchaseHistory_Empty(t *testing.T) {
	interRepo := new(MockInteractionRepository)
	prodRepo := new(MockProductRepository)
	svc := newTestInteractionService(interRepo, prodRepo)
	ctx := context.Background()

	interRepo.On("GetUserPurchases", ctx, 99, 50).Return([]domain.ProductInteraction{}, nil)

	history, err := svc.GetUserPurchaseHistory(ctx, 99, 0)

	assert.NoError(t, err)
	assert.Empty(t, history)
	interRepo.AssertExpectations(t)
}

// --- RecordProductView tests ---

func TestRecordProductView_Success(t *testing.T) {
	interRepo := new(MockInteractionRepository)
	prodRepo := new(MockProductRepository)
	svc := newTestInteractionService(interRepo, prodRepo)
	ctx := context.Background()

	product := &domain.Product{ID: 1, Name: "iPhone"}
	prodRepo.On("GetByID", ctx, 1).Return(product, nil)
	interRepo.On("RecordView", ctx, 5, 1).Return(nil)

	err := svc.RecordProductView(ctx, 5, 1)

	assert.NoError(t, err)
	prodRepo.AssertExpectations(t)
	interRepo.AssertExpectations(t)
}

// --- IsProductLiked tests ---

func TestIsProductLiked_True(t *testing.T) {
	interRepo := new(MockInteractionRepository)
	prodRepo := new(MockProductRepository)
	svc := newTestInteractionService(interRepo, prodRepo)
	ctx := context.Background()

	interRepo.On("HasLiked", ctx, 5, 1).Return(true, nil)

	liked, err := svc.IsProductLiked(ctx, 5, 1)

	assert.NoError(t, err)
	assert.True(t, liked)
	interRepo.AssertExpectations(t)
}

func TestIsProductLiked_False(t *testing.T) {
	interRepo := new(MockInteractionRepository)
	prodRepo := new(MockProductRepository)
	svc := newTestInteractionService(interRepo, prodRepo)
	ctx := context.Background()

	interRepo.On("HasLiked", ctx, 5, 2).Return(false, nil)

	liked, err := svc.IsProductLiked(ctx, 5, 2)

	assert.NoError(t, err)
	assert.False(t, liked)
	interRepo.AssertExpectations(t)
}

// --- New tests for Midterm ---

func TestUnlikeProduct_Success(t *testing.T) {
	interRepo := new(MockInteractionRepository)
	prodRepo := new(MockProductRepository)
	svc := newTestInteractionService(interRepo, prodRepo)
	ctx := context.Background()

	interRepo.On("RemoveLike", ctx, 5, 1).Return(nil)

	err := svc.UnlikeProduct(ctx, 5, 1)

	assert.NoError(t, err)
	interRepo.AssertExpectations(t)
}

func TestHasPurchasedProduct_True(t *testing.T) {
	interRepo := new(MockInteractionRepository)
	prodRepo := new(MockProductRepository)
	svc := newTestInteractionService(interRepo, prodRepo)
	ctx := context.Background()

	interRepo.On("HasPurchased", ctx, 5, 1).Return(true, nil)

	purchased, err := svc.HasPurchasedProduct(ctx, 5, 1)

	assert.NoError(t, err)
	assert.True(t, purchased)
	interRepo.AssertExpectations(t)
}

func TestHasPurchasedProduct_False(t *testing.T) {
	interRepo := new(MockInteractionRepository)
	prodRepo := new(MockProductRepository)
	svc := newTestInteractionService(interRepo, prodRepo)
	ctx := context.Background()

	interRepo.On("HasPurchased", ctx, 5, 99).Return(false, nil)

	purchased, err := svc.HasPurchasedProduct(ctx, 5, 99)

	assert.NoError(t, err)
	assert.False(t, purchased)
	interRepo.AssertExpectations(t)
}

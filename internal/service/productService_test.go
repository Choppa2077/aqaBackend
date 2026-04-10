package service

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/PrimeraAizen/e-comm/internal/domain"
)

// MockProductRepository is a mock implementation of ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) GetByIDWithCategory(ctx context.Context, id int) (*domain.ProductWithCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProductWithCategory), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) ListWithCategories(ctx context.Context, filter domain.ProductFilter) ([]*domain.ProductWithCategory, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.ProductWithCategory), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) Search(ctx context.Context, query string, limit, offset int) ([]*domain.Product, int64, error) {
	args := m.Called(ctx, query, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockProductRepository) GetCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockProductRepository) GetCategoryByName(ctx context.Context, name string) (*domain.Category, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockProductRepository) ListCategories(ctx context.Context) ([]*domain.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *MockProductRepository) UpdateCategory(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockProductRepository) DeleteCategory(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) GetProductStatistics(ctx context.Context, productID int) (*domain.ProductStatistics, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProductStatistics), args.Error(1)
}

func (m *MockProductRepository) RefreshProductStatistics(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// --- CreateProduct tests ---

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	catID := 1
	product := &domain.Product{
		Name:       "Test Product",
		Price:      99.99,
		Stock:      10,
		CategoryID: &catID,
	}

	mockRepo.On("GetCategoryByID", ctx, 1).Return(&domain.Category{ID: 1, Name: "Electronics"}, nil)
	mockRepo.On("Create", ctx, product).Return(nil)

	err := svc.CreateProduct(ctx, product)

	assert.NoError(t, err)
	assert.True(t, product.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_InvalidName(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	product := &domain.Product{
		Name:  "", // empty name — should fail validation
		Price: 99.99,
	}

	err := svc.CreateProduct(ctx, product)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateProduct_InvalidPrice(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	product := &domain.Product{
		Name:  "Valid Name",
		Price: -10.0, // negative price — should fail
	}

	err := svc.CreateProduct(ctx, product)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateProduct_CategoryNotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	catID := 999
	product := &domain.Product{
		Name:       "Test Product",
		Price:      99.99,
		CategoryID: &catID,
	}

	mockRepo.On("GetCategoryByID", ctx, 999).Return(nil, domain.ErrNotFound)

	err := svc.CreateProduct(ctx, product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "category not found")
	mockRepo.AssertExpectations(t)
}

// --- GetProduct tests ---

func TestGetProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	expected := &domain.Product{ID: 1, Name: "iPhone", Price: 999.99}
	mockRepo.On("GetByID", ctx, 1).Return(expected, nil)

	product, err := svc.GetProduct(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, product)
	mockRepo.AssertExpectations(t)
}

func TestGetProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, 999).Return(nil, domain.ErrNotFound)

	product, err := svc.GetProduct(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
	assert.Nil(t, product)
	mockRepo.AssertExpectations(t)
}

// --- UpdateProduct tests ---

func TestUpdateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	existing := &domain.Product{ID: 1, Name: "Old Name", Price: 50.0, Stock: 10}
	updated := &domain.Product{ID: 1, Name: "New Name", Price: 75.0, Stock: 10}

	mockRepo.On("GetByID", ctx, 1).Return(existing, nil)
	mockRepo.On("Update", ctx, updated).Return(nil)

	err := svc.UpdateProduct(ctx, updated)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	product := &domain.Product{ID: 999, Name: "Ghost Product", Price: 50.0}
	mockRepo.On("GetByID", ctx, 999).Return(nil, domain.ErrNotFound)

	err := svc.UpdateProduct(ctx, product)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// --- DeleteProduct tests ---

func TestDeleteProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	existing := &domain.Product{ID: 1, Name: "Product to delete"}
	mockRepo.On("GetByID", ctx, 1).Return(existing, nil)
	mockRepo.On("Delete", ctx, 1).Return(nil)

	err := svc.DeleteProduct(ctx, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, 999).Return(nil, domain.ErrNotFound)

	err := svc.DeleteProduct(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
	mockRepo.AssertExpectations(t)
}

// --- ListProducts tests ---

func TestListProducts_WithFilter(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	expected := []*domain.Product{
		{ID: 1, Name: "Product A", Price: 10.0},
		{ID: 2, Name: "Product B", Price: 20.0},
	}

	// Service sets IsActive=true and Limit=20 by default when Limit=0
	mockRepo.On("List", ctx, mock.AnythingOfType("domain.ProductFilter")).Return(expected, int64(2), nil)

	products, total, err := svc.ListProducts(ctx, domain.ProductFilter{Limit: 10, Offset: 0})

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, products, 2)
	mockRepo.AssertExpectations(t)
}

// --- Edge cases & new tests for Midterm ---

func TestUpdateStock_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	existing := &domain.Product{ID: 1, Name: "Laptop", Price: 999.99, Stock: 10}
	updated := &domain.Product{ID: 1, Name: "Laptop", Price: 999.99, Stock: 15}

	mockRepo.On("GetByID", ctx, 1).Return(existing, nil)
	mockRepo.On("Update", ctx, updated).Return(nil)

	err := svc.UpdateStock(ctx, 1, 5) // add 5 units

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateStock_ProductNotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, 999).Return(nil, domain.ErrNotFound)

	err := svc.UpdateStock(ctx, 999, 5)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestCheckStock_Sufficient(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	product := &domain.Product{ID: 1, Name: "Phone", Price: 499.99, Stock: 20}
	mockRepo.On("GetByID", ctx, 1).Return(product, nil)

	ok, err := svc.CheckStock(ctx, 1, 10)

	assert.NoError(t, err)
	assert.True(t, ok)
	mockRepo.AssertExpectations(t)
}

func TestCheckStock_Insufficient(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	product := &domain.Product{ID: 1, Name: "Phone", Price: 499.99, Stock: 3}
	mockRepo.On("GetByID", ctx, 1).Return(product, nil)

	ok, err := svc.CheckStock(ctx, 1, 10)

	assert.NoError(t, err)
	assert.False(t, ok)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_VeryLongName(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	// Name with 300 characters — service does not enforce max length,
	// so this should pass validation and reach repo.Create.
	longName := ""
	for i := 0; i < 300; i++ {
		longName += "a"
	}
	catID := 1
	product := &domain.Product{
		Name:       longName,
		Price:      9.99,
		CategoryID: &catID,
	}

	mockRepo.On("GetCategoryByID", ctx, 1).Return(&domain.Category{ID: 1, Name: "Electronics"}, nil)
	mockRepo.On("Create", ctx, product).Return(nil)

	err := svc.CreateProduct(ctx, product)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPurchaseProduct_Concurrent(t *testing.T) {
	// Concurrency edge case: 5 goroutines attempt to purchase the last item in stock.
	// Each goroutine uses its own isolated mock (simulating separate DB transactions).
	// Goroutine 0 sees stock=1 and succeeds; goroutines 1-4 see stock=0 and fail.
	const buyers = 5
	results := make([]error, buyers)
	var wg sync.WaitGroup

	for i := 0; i < buyers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			interRepo := new(MockInteractionRepository)
			prodRepo := new(MockProductRepository)
			svc := NewInteractionService(interRepo, prodRepo)
			ctx := context.Background()

			if idx == 0 {
				// First buyer gets the last unit
				product := &domain.Product{ID: 1, Price: 99.99, Stock: 1}
				updated := &domain.Product{ID: 1, Price: 99.99, Stock: 0}
				prodRepo.On("GetByID", ctx, 1).Return(product, nil)
				interRepo.On("RecordPurchase", ctx, 0, 1, 1, 99.99).Return(nil)
				prodRepo.On("Update", ctx, updated).Return(nil)
			} else {
				// Remaining buyers see empty stock
				product := &domain.Product{ID: 1, Price: 99.99, Stock: 0}
				prodRepo.On("GetByID", ctx, 1).Return(product, nil)
			}

			results[idx] = svc.PurchaseProduct(ctx, idx, 1, 1)
		}(i)
	}

	wg.Wait()

	successCount := 0
	for _, err := range results {
		if err == nil {
			successCount++
		}
	}

	assert.Equal(t, 1, successCount, "only one buyer should succeed when stock=1")
}

package service

import (
	"context"
	"errors"

	"github.com/online-shop/internal/models"
	"github.com/online-shop/internal/repository"
)

var ErrInsufficientStock = errors.New("insufficient stock")

type OrderService struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
}

func NewOrderService(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo, productRepo: productRepo}
}

// func (s *OrderService) Create(ctx context.Context, userID int64, req models.CreateOrderRequest) (*models.Order, error) {
// 	var totalPrice float64
// 	var items []models.OrderItem

// 	for _, ri := range req.Items {
// 		product, err := s.productRepo.GetByID(ctx, ri.ProductID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if product.Stock < ri.Quantity {
// 			return nil, ErrInsufficientStock
// 		}
// 		items = append(items, models.OrderItem{
// 			ProductID: ri.ProductID,
// 			Quantity:  ri.Quantity,
// 			Price:     product.Price,
// 		})
// 		totalPrice += product.Price * float64(ri.Quantity)
// 	}

// 	order := &models.Order{
// 		UserID:     userID,
// 		Status:     models.OrderStatusPending,
// 		TotalPrice: totalPrice,
// 		Address:    req.Address,
// 		Items:      items,
// 	}

// 	if err := s.orderRepo.Create(ctx, order); err != nil {
// 		return nil, err
// 	}

// 	return order, nil
// }

func (s *OrderService) GetByID(ctx context.Context, id int64) (*models.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}

func (s *OrderService) ListByUserID(ctx context.Context, userID int64) ([]models.Order, error) {
	return s.orderRepo.ListByUserID(ctx, userID)
}

func (s *OrderService) UpdateStatus(ctx context.Context, id int64, status models.OrderStatus) error {
	return s.orderRepo.UpdateStatus(ctx, id, status)
}

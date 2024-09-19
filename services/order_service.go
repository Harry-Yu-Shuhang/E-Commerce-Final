package services

import (
	"imooc-product/datamodels"
	"imooc-product/repositories"
)

type IOrderService interface {
	GetOederByID(int64) (*datamodels.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(*datamodels.Order) error
	InsertOrder(*datamodels.Order) (int64, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
}

func NewOrderService(repository repositories.IOrderRepository) IOrderService {
	return &OrderService{OrderRepository: repository}
}

type OrderService struct {
	OrderRepository repositories.IOrderRepository
}

func (o *OrderService) GetOederByID(orderID int64) (order *datamodels.Order, err error) {
	order, err = o.OrderRepository.SelectByKey(orderID)
	return
}

func (o *OrderService) DeleteOrderByID(orderID int64) (isOk bool) {
	isOk = o.OrderRepository.Delete(orderID)
	return
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) (err error) {
	err = o.OrderRepository.Update(order)
	return
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (orderID int64, err error) {
	orderID, err = o.OrderRepository.Insert(order)
	return
}

func (o *OrderService) GetAllOrder() (orders []*datamodels.Order, err error) {
	orders, err = o.OrderRepository.SelectAll()
	return
}

func (o *OrderService) GetAllOrderInfo() (orderInfo map[int]map[string]string, err error) {
	orderInfo, err = o.OrderRepository.SelectAllWithInfo()
	return
}

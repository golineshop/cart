package service

import (
	"github.com/golineshop/cart/domain/model"
	"github.com/golineshop/cart/domain/repository"
)

type ICartService interface {
	AddCart(cart *model.Cart) (cartID int64, err error)
	DeleteCart(cartID int64) (err error)
	UpdateCart(cart *model.Cart) (err error)
	FindCartByID(cartID int64) (cart *model.Cart, err error)
	FindAllCart(userID int64) (cartList []model.Cart, err error)
	CleanCart(userID int64) (err error)
	DecrNum(cartID int64, num int64) (err error)
	IncrNum(cartID int64, num int64) (err error)
}

func NewCartService(cartRepository repository.ICartRepository) ICartService {
	return &CartService{cartRepository: cartRepository}
}

type CartService struct {
	cartRepository repository.ICartRepository
}

func (c *CartService) AddCart(cart *model.Cart) (cartID int64, err error) {
	return c.cartRepository.CreateCart(cart)
}

func (c *CartService) DeleteCart(cartID int64) (err error) {
	return c.cartRepository.DeleteCartByID(cartID)
}

func (c *CartService) UpdateCart(cart *model.Cart) (err error) {
	return c.UpdateCart(cart)
}

func (c *CartService) FindCartByID(cartID int64) (cart *model.Cart, err error) {
	return c.FindCartByID(cartID)
}

func (c *CartService) FindAllCart(userID int64) (cartList []model.Cart, err error) {
	return c.FindAllCart(userID)
}

func (c *CartService) CleanCart(userID int64) (err error) {
	return c.CleanCart(userID)
}

func (c *CartService) DecrNum(cartID int64, num int64) (err error) {
	return c.DecrNum(cartID, num)
}

func (c *CartService) IncrNum(cartID int64, num int64) (err error) {
	return c.IncrNum(cartID, num)
}

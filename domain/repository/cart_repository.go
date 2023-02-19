package repository

import (
	"errors"
	"github.com/golineshop/cart/domain/model"
	"github.com/jinzhu/gorm"
)

type ICartRepository interface {
	InitTable() (err error)
	FindCartByID(cartID int64) (cart *model.Cart, err error)
	CreateCart(cart *model.Cart) (ID int64, err error)
	DeleteCartByID(cartID int64) (err error)
	UpdateCart(cart *model.Cart) (err error)
	FindAll(userID int64) (cartList []model.Cart, err error)
	CleanCart(userID int64) (err error)
	IncrNum(cartID int64, num int64) (err error)
	DecrNum(cartID int64, num int64) (err error)
}

func NewCartRepository(db *gorm.DB) ICartRepository {
	return &CartRepository{mysqlDB: db}
}

type CartRepository struct {
	mysqlDB *gorm.DB
}

func (c *CartRepository) InitTable() (err error) {
	return c.mysqlDB.CreateTable(&model.Cart{}).Error
}

func (c *CartRepository) FindCartByID(cartID int64) (cart *model.Cart, err error) {
	cart = &model.Cart{}
	return cart, c.mysqlDB.First(cartID).Error
}

func (c *CartRepository) CreateCart(cart *model.Cart) (ID int64, err error) {
	db := c.mysqlDB.FirstOrCreate(cart, model.Cart{ProductID: cart.ProductID, SizeID: cart.SizeID, UserID: cart.UserID})
	if db.Error != nil {
		return 0, db.Error
	}
	if db.RowsAffected == 0 {
		return 0, errors.New("购物车插入失败")
	}
	return cart.ID, nil
}

func (c *CartRepository) DeleteCartByID(cartID int64) (err error) {
	return c.mysqlDB.Where("id = ?", cartID).Delete(&model.Cart{}).Error
}

func (c *CartRepository) UpdateCart(cart *model.Cart) (err error) {
	return c.mysqlDB.Model(cart).Update(cart).Error
}

func (c *CartRepository) FindAll(userID int64) (cartList []model.Cart, err error) {
	return cartList, c.mysqlDB.Where("user_id = ?", userID).Find(&cartList).Error
}

func (c *CartRepository) CleanCart(userID int64) (err error) {
	return c.mysqlDB.Where("user_id = ?", userID).Delete(&model.Cart{}).Error
}

func (c *CartRepository) IncrNum(cartID int64, num int64) (err error) {
	cart := &model.Cart{ID: cartID}
	return c.mysqlDB.Model(cart).UpdateColumn("num", gorm.Expr("num + ?", num)).Error
}

func (c *CartRepository) DecrNum(cartID int64, num int64) (err error) {
	cart := &model.Cart{ID: cartID}
	db := c.mysqlDB.Model(cart).Where("num >= ?", num).UpdateColumn("num", gorm.Expr("num - ?", num))
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return errors.New("减少失败")
	}
	return nil
}

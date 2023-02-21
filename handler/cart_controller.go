package handler

import (
	"context"
	"github.com/golineshop/cart/domain/model"
	"github.com/golineshop/cart/domain/service"
	proto "github.com/golineshop/cart/proto"
	"github.com/golineshop/common"
)

type CartController struct {
	CartService service.ICartService
}

func (c *CartController) AddCart(ctx context.Context, request *proto.CartInfo, response *proto.ResponseAdd) (err error) {
	cart := &model.Cart{}
	if err = common.SwapTo(request, cart); err != nil {
		return err
	}
	response.CartId, err = c.CartService.AddCart(cart)
	return err
}

func (c *CartController) CleanCart(ctx context.Context, request *proto.Clean, response *proto.Response) (err error) {
	if err = c.CartService.CleanCart(request.UserId); err != nil {
		return err
	}
	response.Meg = "购物车清空成功"
	return nil
}

func (c *CartController) Incr(ctx context.Context, request *proto.Item, response *proto.Response) (err error) {
	if err = c.CartService.IncrNum(request.Id, request.ChangeNum); err != nil {
		return err
	}
	response.Meg = "购物车添加成功"
	return nil
}

func (c *CartController) Decr(ctx context.Context, request *proto.Item, response *proto.Response) (err error) {
	if err = c.CartService.DecrNum(request.Id, request.ChangeNum); err != nil {
		return err
	}
	response.Meg = "购物程减少成功"
	return nil
}

func (c *CartController) DeleteItemByID(ctx context.Context, request *proto.CartID, response *proto.Response) (err error) {
	if err = c.CartService.DeleteCart(request.Id); err != nil {
		return err
	}
	response.Meg = "购物车删除成功"
	return nil
}

func (c *CartController) GetAll(ctx context.Context, request *proto.CartFindAll, response *proto.CartAll) (err error) {
	cartAll, err := c.CartService.FindAllCart(request.UserId)
	if err != nil {
		return err
	}

	for _, v := range cartAll {
		cart := &proto.CartInfo{}
		if err := common.SwapTo(v, cart); err != nil {
			return err
		}
		response.CartInfo = append(response.CartInfo, cart)
	}
	return nil
}

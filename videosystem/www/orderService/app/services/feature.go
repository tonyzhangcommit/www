package services

import (
	"errors"
	"userservice/app/request"
	"userservice/models"
)

/*
	业务逻辑实现
*/

// 秒杀活动下单
func TakeFlashOrder(form *request.TakeFlashOrder) (order models.Order, err error) {
	err = errors.New("下单失败")
	return
}

// 查询订单
func GetOrder(form *request.TakeFlashOrder) (order models.Order, err error) {
	return
}

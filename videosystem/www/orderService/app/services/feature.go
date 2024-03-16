package services

import (
	"errors"
	"userservice/app/request"
	"userservice/models"
)

/*
	业务逻辑实现
*/

/*
	业务流程：
	将所有请求存放至指定的消息队列中
	暂定消息队列有如下几个：
	1. 订单队列
	2. 死信队列
	3. 结果队列
*/

func TakeFlashOrder(form *request.TakeFlashOrder) (order models.Order, err error) {
	err = errors.New("下单失败")
	return
}

// 查询订单
func GetOrder(form *request.TakeFlashOrder) (order models.Order, err error) {
	return
}

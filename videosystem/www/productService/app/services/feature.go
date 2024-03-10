package services

import (
	"userservice/app/request"
	"userservice/app/response"
)

/*
	业务逻辑实现,为了区分不同的功能模块，这里新定义一个空结构体
*/

type flashEvent struct{}

var FE = new(flashEvent)

// 获取活动以及商品信息
func (f *flashEvent) GetFEventProduct(form *request.GetFlashEventProduct) (fpinfo response.FlashEventProduct, err error) {
	return
}

// 获取活动基本信息
func (f *flashEvent) GetFEventInfo(form *request.GetFlashEvent) (finfo response.FlashEvent, err error) {
	return
}

package management

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"userservice/app/request"
	"userservice/app/response"
	"userservice/app/services"
	"userservice/global"

	"github.com/gin-gonic/gin"
)

/*
	视图函数
*/

// 秒杀活动下单
func TakeFlashOrder(c *gin.Context) {
	var form request.TakeFlashOrder
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	res, err := services.TakeFlashOrder(&form)
	if err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, res)
	}
}

func Test(c *gin.Context) {
	// getmessage("flasheventreqqueue")
	global.SendLogs("info", "test logs")
	getmessage("orderinfoqueue")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// 获取队列信息
func getmessage(qname string) {
	type QueueInfo struct {
		Name      string `json:"name"`
		Messages  int    `json:"messages"`
		Consumers int    `json:"consumers"`
	}
	url := fmt.Sprintf("http://localhost:15672/api/queues/%%2fvideosystemdev/%s", qname)
	// 创建一个HTTP客户端并设置请求
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 替换为你的RabbitMQ管理界面的用户名和密码
	req.SetBasicAuth("adminroot", "adminroot")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 读取并解析响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var queueInfo QueueInfo
	if err := json.Unmarshal(body, &queueInfo); err != nil {
		log.Fatal(err)
	}
	fmt.Println("队列详情", queueInfo)
}

// 普通下单
func TakeOrder(c *gin.Context) {

}

// 修改订单
func EditOrder(c *gin.Context) {

}

// 取消订单
func CancelOrder(c *gin.Context) {

}

// 删除订单
func DeleteOrder(c *gin.Context) {

}

// 查看订单
func GetOrder(c *gin.Context) {

}

// 申请退款
func RefundOrder(c *gin.Context) {

}

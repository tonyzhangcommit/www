package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
	"userservice/app/request"
	"userservice/global"
	"userservice/models"
	"userservice/utils"

	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
)

/*
	业务逻辑实现
*/

var WebSocketclients = make(map[string]*websocket.Conn)

// websocket 发送消息函数
func WebsocketSendMessage(userID, message string) {
	ws, ok := WebSocketclients[userID]
	if !ok {
		// 用户连接不存在
		return
	}

	// 向用户发送消息
	if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		// 错误处理
	}
}

/*
	业务流程
*/

type FlashReqOrder struct {
	EventID   uint   `json:"eventid"`
	ProductID uint   `json:"peoductid"`
	UserID    uint   `json:"userid"`
	Count     int    `json:"count"`
	Uid       string `json:"uid"`
}

type FlashResOrder struct {
	FlashReqOrder
	ProductName    string  `json:"productname"`
	OriginPrice    float64 `json:"originprice"`
	FlashSalePrice float64 `json:"flashsaleprice"`
}

// 秒杀订单处理逻辑
func TakeFlashOrder(form *request.TakeFlashOrder) (res string, err error) {
	// 为每个请求生成唯一性id
	uid := utils.GenerateUniqueID()
	// 创建消息结构体
	order := FlashReqOrder{
		EventID:   form.EventID,
		ProductID: form.ProductID,
		UserID:    form.UserID,
		Count:     form.Count,
		Uid:       uid,
	}

	body, err := json.Marshal(order)
	if err != nil {
		global.SendLogs("error", "秒杀请求-序列化请求体错误", err)
		err = errors.New("请求错误")
		return
	}

	// 发布消息
	err = global.RabbitMQ.Channel.PublishWithContext(
		context.Background(),     // 上下文
		global.RabbitMQ.Exchange, // 交换机名称
		"flashevent.req.row",     // 路由键
		false,                    // mandatory
		false,                    // immediate
		amqp091.Publishing{
			ContentType:     "application/json",             // 内容类型
			ContentEncoding: "utf-8",                        // 内容编码
			DeliveryMode:    amqp091.Persistent,             // 持久化消息
			Priority:        0,                              // 优先级
			CorrelationId:   strconv.Itoa(int(form.UserID)), // 关联ID
			Expiration:      "3600000",                      // 消息过期时间，单位为毫秒
			Timestamp:       time.Now(),                     // 时间戳
			Type:            "order_flashevent_message",     // 消息类型
			UserId:          strconv.Itoa(int(form.UserID)), // 用户ID
			AppId:           "order_flashevent",             // 应用ID
			Body:            body,                           // 消息正文
		},
	)
	if err != nil {
		global.SendLogs("error", "秒杀请求-请求入队错误", err)
		err = errors.New("请求内部错误")
		return
	}
	res = "请求成功"
	return
}

// 秒杀订单请求消费者(参数为消费者的数量
// 需要判断消息是否重复消费
// 订单预生成，将价格计算逻辑，订单入库sql
func FlashEventCustomer(consumercount int) {
	for cid := 1; cid < consumercount+1; cid++ {
		go func(consumerID int) {
			// 初始化监听队列
			msgs, err := global.RabbitMQ.Channel.Consume(
				global.App.Config.RabbitMQ.FlashEventReqQueue,
				strconv.Itoa(cid),
				true,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				// 服务挂掉了
				global.SendLogs("error", "秒杀活动宕机-------------->", err)
				return
			}
			for msg := range msgs {
				var order FlashReqOrder
				err := json.Unmarshal(msg.Body, &order)
				if err != nil {
					global.SendLogs("error", fmt.Sprintf("秒杀订单消费失败：Consumer %d: Error decoding JSON: %v", consumerID, err))
					// 这里通过websocket通知用户订单抢购失败
					go WebsocketSendMessage(msg.CorrelationId, "处理您的请求时出现问题，请检查订单详情或联系客服。")
					// 丢弃消息
					msg.Ack(false)
					continue
				}
				// 首先判断此消息是否已经消费
				ordercachekey := "flashevent:order:req:" + order.Uid
				key := utils.JoinStrings("flashID:", strconv.Itoa(int(order.EventID)), "pid", strconv.Itoa(int(order.ProductID)), ":remainingquantity") // 商品库存键
				if _, err := global.App.Redis.Get(context.Background(), ordercachekey).Result(); err != nil {
					// 使用 Redis 的 DECRBY 命令原子减少商品数量
					result, err := global.App.Redis.DecrBy(context.Background(), key, int64(order.Count)).Result()
					if err != nil {
						// 这里通过websocket通知用户订单抢购失败
						global.SendLogs("error", fmt.Sprintf("秒杀订单抢购失败：Consumer %d", consumerID), err)
						go WebsocketSendMessage(msg.CorrelationId, "当前抢购的人太多了，请稍后再试")
						continue
					}
					if result < 0 {
						// 库存不足，回滚
						_, err := global.App.Redis.IncrBy(context.Background(), key, int64(order.Count)).Result()
						if err != nil {
							global.SendLogs("error", fmt.Sprintf("秒杀订单库存不足：Consumer %d", consumerID), err)
							go WebsocketSendMessage(msg.CorrelationId, "手慢一步，商品已经抢光了哦~")
							continue
						}
					} else {
						// 此处可进行后续处理，如通知用户、生成订单等
						go WebsocketSendMessage(msg.CorrelationId, "抢购成功，正在为您生成订单......")
						// 将抢购结果发送至结果队列
						pname := ""
						originprice := ""
						flashprice := ""
						pnamekey := utils.JoinStrings("flashID:", strconv.Itoa(int(order.EventID)), "pid", strconv.Itoa(int(order.ProductID)), ":", "pname")
						if pname, errredis := global.App.Redis.Get(context.Background(), pnamekey).Result(); errredis != nil {
							pname = ""
						}

						originpricekey := utils.JoinStrings("flashID:", strconv.Itoa(int(order.EventID)), "pid", strconv.Itoa(int(order.ProductID)), ":", "originprice")
						if originprice, errredis := global.App.Redis.Get(context.Background(), originpricekey).Result(); errredis != nil {
							originprice = "0.0"
						}
						flashpricekey := utils.JoinStrings("flashID:", strconv.Itoa(int(order.EventID)), "pid", strconv.Itoa(int(order.ProductID)), ":", "flashprice")
						if flashprice, errredis := global.App.Redis.Get(context.Background(), flashpricekey).Result(); errredis != nil {
							flashprice = "0.0"
						}

						resorder := FlashResOrder{
							FlashReqOrder:  order,
							ProductName:    pname,
							OriginPrice:    originprice,
							FlashSalePrice: flashprice,
						}
						body, err := json.Marshal(order)
						if err != nil {
							go global.SendLogs("error", "秒杀订单抢购成功序列化失败", err)
							go WebsocketSendMessage(msg.CorrelationId, "生成订单失败，请联系客服")
							continue
						}
						// 发布消息
						err = global.RabbitMQ.Channel.PublishWithContext(
							context.Background(),     // 上下文
							global.RabbitMQ.Exchange, // 交换机名称
							"flashevent.res.row",     // 路由键
							false,                    // mandatory
							false,                    // immediate
							amqp091.Publishing{
								ContentType:     "application/json",              // 内容类型
								ContentEncoding: "utf-8",                         // 内容编码
								DeliveryMode:    amqp091.Persistent,              // 持久化消息
								Priority:        0,                               // 优先级
								CorrelationId:   strconv.Itoa(int(order.UserID)), // 关联ID
								Expiration:      "3600000",                       // 消息过期时间，单位为毫秒
								Timestamp:       time.Now(),                      // 时间戳
								Type:            "order_flashevent_message",      // 消息类型
								UserId:          strconv.Itoa(int(order.UserID)), // 用户ID
								AppId:           "order_flashevent",              // 应用ID
								Body:            body,                            // 消息正文
							},
						)
						if err != nil {
							go global.SendLogs("error", "秒杀订单抢购成功后入队生成订单队列失败", err)
							go WebsocketSendMessage(msg.CorrelationId, "生成订单失败，请联系客服")
							continue
						}
					}
					// 将消息已处理信息保存至redis
					if err := global.App.Redis.Set(context.Background(), ordercachekey, "1", 3*time.Hour).Err(); err != nil {
						// 这里增加本地缓存备份，待完善------->
						global.SendLogs("error", fmt.Sprintf("秒杀订单处理后添加缓存标记失败：Consumer %d", consumerID), err)
					}
					msg.Ack(false)
				} else {
					// 已经被消费过了
					msg.Ack(false)
					continue
				}
			}
		}(cid)
	}
	forever := make(chan bool)
	<-forever
}

// 抢购结果消费者（主要为入库操作）
func FlashEventsnapupresCustomer() {
	// 初始化监听队列
	msgs, err := global.RabbitMQ.Channel.Consume(
		global.App.Config.RabbitMQ.FlashEventReqQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		// 服务挂掉了
		global.SendLogs("error", "秒杀活动生成订单队列宕机-------------->", err)
		go WebsocketSendMessage("all_user", "生成订单失败，请联系客服")
		return
	}
	for msg := range msgs {
		var order FlashReqOrder
		err := json.Unmarshal(msg.Body, &order)
		if err != nil {
			global.SendLogs("error", "秒杀活动生成订单--序列化结构体失败", err)
			// 这里通过websocket通知用户订单抢购失败
			go WebsocketSendMessage(msg.CorrelationId, "生成订单失败，详情请联系客服。")
			// 丢弃消息
			msg.Ack(false)
			continue
		}
		orderrescachekey := "flashevent:order:res:" + order.Uid
		if _, err := global.App.Redis.Get(context.Background(), orderrescachekey).Result(); err != nil {

			// 开启事务
			tx := global.App.DB.Begin()
			// 插入订单
			order := models.Order{
				UserID:        1,
				OrderType:     "normal",
				Status:        "pending",
				PaymentStatus: "unpaid",
				PayWay:        "card",
				TotalAmount:   100.0,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			if err := tx.Create(&order).Error; err != nil {
				tx.Rollback() // 回滚事务
				log.Fatalf("failed to insert order: %v", err)
			}
			// 生成订单item
			// 插入订单项
			orderItem := models.OrderItem{
				OrderID:      order.ID,
				ProductID:    1,
				ProductName:  "Product A",
				Quantity:     2,
				UnitPrice:    50.0,
				Subtotal:     100.0,
				OriginAmount: 120.0,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			if err := tx.Create(&orderItem).Error; err != nil {
				tx.Rollback() // 回滚事务
				log.Fatalf("failed to insert order item: %v", err)
			}

			// 提交事务
			if err := tx.Commit().Error; err != nil {
				log.Fatalf("failed to commit transaction: %v", err)
			}

		} else {
			// 消息已经处理，丢弃
			msg.Ack(false)
			continue
		}
		time.Sleep(100 * time.Millisecond) // 限制速率
	}
	forever := make(chan bool)
	<-forever
}

// 查询订单
func GetOrder(form *request.TakeFlashOrder) (order models.Order, err error) {
	return
}

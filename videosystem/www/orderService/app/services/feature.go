package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
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
var ClientsLock = sync.RWMutex{}

// websocket 发送消息函数
func WebsocketSendMessage(userID, message string) {
	if userID == "allUser" {
		// 向全部用户发送消息
		for _, ws := range WebSocketclients {
			// 向用户发送消息
			if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				// 错误处理
				global.SendLogs("error", fmt.Sprintf("websocket 向用户 %s 消息失败", "allUser"), err)
			}
		}

	} else {
		ws, ok := WebSocketclients[userID]
		if !ok {
			// 用户连接不存在
			return
		}
		// 向用户发送消息
		if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			// 错误处理
			global.SendLogs("error", fmt.Sprintf("websocket 向用户 %s 消息失败", userID), err)
		}
	}

}

/*
	业务流程
*/

// 秒杀活动请求消息格式
type FlashReqOrder struct {
	EventID   uint   `json:"eventid"`
	ProductID uint   `json:"peoductid"`
	UserID    uint   `json:"userid"`
	Count     int    `json:"count"`
	Uid       string `json:"uid"`
}

// 生成订单消息格式
type FlashResOrder struct {
	FlashReqOrder
	ProductName    string  `json:"productname"`
	OriginPrice    float64 `json:"originprice"`
	FlashSalePrice float64 `json:"flashsaleprice"`
}

// 待付款消息格式
type PendingPayment struct {
	OrderID uint   `json:"orderid"`
	Uid     string `json:"uid"`
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
			ch, err := global.RabbitMQ.Conn.Channel()
			if err != nil {
				global.SendLogs("error", "秒杀活动创建请求消费者失败", err)
				return
			}
			msgs, err := ch.Consume(
				global.App.Config.RabbitMQ.FlashEventReqQueue,
				strconv.Itoa(consumerID),
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				// 服务挂掉了
				global.SendLogs("error", "秒杀活动-获取订单消息失败", err)
				return
			}
			for msg := range msgs {
				var order FlashReqOrder
				err := json.Unmarshal(msg.Body, &order)
				if err != nil {
					global.SendLogs("error", fmt.Sprintf("秒杀订单消费失败：Consumer %d: Error decoding JSON: %v", 1, err))
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
						go WebsocketSendMessage(msg.CorrelationId, "当前抢购的人太多啦~，请稍后再试")
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

						// 将抢购结果发送至结果队列
						pname := ""
						originprice := 0.0
						flashprice := 0.0
						// 从缓存中获取该秒杀活动中的商品名称，原始价格，秒杀价格
						pnamekey := utils.JoinStrings("flashID:", strconv.Itoa(int(order.EventID)), "pid", strconv.Itoa(int(order.ProductID)), ":", "pname")
						if pnameRedis, errredis := global.App.Redis.Get(context.Background(), pnamekey).Result(); errredis != nil {
							global.SendLogs("error", fmt.Sprintf("redis 获取秒杀活动商品名称错误：Key %s", pnamekey), err)
						} else {
							pname = pnameRedis
						}
						originpricekey := utils.JoinStrings("flashID:", strconv.Itoa(int(order.EventID)), "pid", strconv.Itoa(int(order.ProductID)), ":", "originprice")
						if originpriceRedis, errredis := global.App.Redis.Get(context.Background(), originpricekey).Result(); errredis != nil {
							global.SendLogs("error", fmt.Sprintf("redis 获取秒杀活动商品原始价格错误：Key %s", originpricekey), err)
						} else {
							// 转化为float类型
							floatOPrice, err := strconv.ParseFloat(originpriceRedis, 64)
							if err != nil {
								global.SendLogs("error", fmt.Sprintf("redis 秒杀活动商品原始价格转float错误：Value %s", originpriceRedis), err)
							} else {
								originprice = floatOPrice
							}
						}
						flashpricekey := utils.JoinStrings("flashID:", strconv.Itoa(int(order.EventID)), "pid", strconv.Itoa(int(order.ProductID)), ":", "flashprice")
						if flashpriceRedis, errredis := global.App.Redis.Get(context.Background(), flashpricekey).Result(); errredis != nil {
							global.SendLogs("error", fmt.Sprintf("redis 获取秒杀活动商品秒杀价格错误：Key %s", originpricekey), err)
						} else {
							// 转化为float类型
							floatflashpriceRedis, err := strconv.ParseFloat(flashpriceRedis, 64)
							if err != nil {
								global.SendLogs("error", fmt.Sprintf("redis 秒杀活动商品秒杀价格转float错误：Value %s", flashpriceRedis), err)
							} else {
								flashprice = floatflashpriceRedis
							}
						}
						resorder := FlashResOrder{
							FlashReqOrder:  order,
							ProductName:    pname,
							OriginPrice:    originprice,
							FlashSalePrice: flashprice,
						}
						fmt.Println("抢购后消息格式信息", resorder)
						// 构建订单抢购后的消息体，用于生成订单
						body, err := json.Marshal(resorder)
						if err != nil {
							go global.SendLogs("error", "秒杀订单抢购成功序列化失败", err)
							go WebsocketSendMessage(msg.CorrelationId, "生成订单失败，请联系客服")
							continue
						}
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
								AppId:           "order_flashevent",              // 应用ID
								Body:            body,                            // 消息正文
							},
						)
						if err != nil {
							go global.SendLogs("error", "秒杀订单抢购成功后入队生成订单队列失败", err)
							go WebsocketSendMessage(msg.CorrelationId, "生成订单失败，请联系客服")
							continue
						}
						// 抢购成功
						go WebsocketSendMessage(msg.CorrelationId, "抢购成功，正在为您生成订单......")
						// 将消息已处理信息保存至redis
						if err := global.App.Redis.Set(context.Background(), ordercachekey, "1", 3*time.Hour).Err(); err != nil {
							// 这里增加本地缓存备份，待完善------->
							global.SendLogs("error", fmt.Sprintf("秒杀订单处理后添加缓存标记失败：Consumer %d", consumerID), err)
						}
					}
				} else {
					// 此消息已经被消费过了，直接废弃
					msg.Ack(false)
				}
				msg.Ack(false)
			}
		}(cid)
	}
}

// 抢购结果消费者（主要为入库操作）
func FlashEventsnapupresCustomer() {
	go func() {
		// 初始化监听队列
		ch, err := global.RabbitMQ.Conn.Channel()
		if err != nil {
			global.SendLogs("error", "秒杀活动创建生成订单消费者失败", err)
			go WebsocketSendMessage("allUser", "生成订单失败001，详情请联系客服。")
			return
		}
		msgs, err := ch.Consume(
			global.App.Config.RabbitMQ.FlashEventResQueue,
			"", // 这里只定义一个消费者
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			// 服务挂掉了
			global.SendLogs("error", "秒杀活动生成订单队列宕机-------------->", err)
			go WebsocketSendMessage("allUser", "生成订单失败002，详情请联系客服。")
			return
		}
		for msg := range msgs {
			var orderres FlashResOrder
			err := json.Unmarshal(msg.Body, &orderres)
			if err != nil {
				global.SendLogs("error", fmt.Sprintf("秒杀活动生成订单--序列化结构体失败，userID:%s", msg.CorrelationId), err)
				// 这里通过websocket通知用户订单抢购失败
				go WebsocketSendMessage(msg.CorrelationId, "生成订单失败003，详情请联系客服。")
				// 丢弃消息
				msg.Ack(false)
				continue
			}
			// 防止重复消费
			orderrescachekey := "flashevent:order:res:" + orderres.Uid
			if _, err := global.App.Redis.Get(context.Background(), orderrescachekey).Result(); err != nil {
				flashprice := orderres.FlashSalePrice
				originprice := orderres.OriginPrice
				productname := orderres.ProductName
				type flashPruductPrice struct {
					OriginalPrice  float64 `gorm:"column:originalprice"`
					FlashSalePrice float64 `gorm:"column:flashsaleprice"`
				}

				type ProductName struct {
					Name string `gorm:"column:name;type:varchar(100);not null"`
				}
				// 表示初步缓存失败：（抢购过程获取失败）
				if flashprice < 0.1 {
					// 表示上一步中构建消息时获取商品秒杀价格失败，这里需要重新获取商品信息
					sqlSelectPrice := "select originalprice,flashsaleprice from flasheventproduct where eventid= '%d' and productid = '%d'"
					var flashpruductpriceres flashPruductPrice
					if err := global.App.DB.Raw(sqlSelectPrice).Scan(&flashpruductpriceres).Error; err != nil {
						// 查询出错，生成订单失败
						global.SendLogs("error", "秒杀活动---生成订单---查询秒杀商品价格失败")
						go WebsocketSendMessage(msg.CorrelationId, "生成订单失败，详情请联系客服。")
						continue
					}
					flashprice = flashpruductpriceres.FlashSalePrice
					originprice = flashpruductpriceres.OriginalPrice
				}
				if productname == "" {
					sqlSelectPName := "select name from membership where id = '%d'"
					var productnameres ProductName
					if err := global.App.DB.Raw(sqlSelectPName).Scan(&productnameres).Error; err != nil {
						// 查询出错，生成订单失败
						global.SendLogs("error", "秒杀活动---生成订单---查询秒杀商品名称失败")
						go WebsocketSendMessage(msg.CorrelationId, "生成订单失败，详情请联系客服。")
						continue
					}
				}
				// 计算订单价格
				totalAmount := float64(orderres.Count) * flashprice
				originAmount := float64(orderres.Count) * originprice

				// 开启事务
				tx := global.App.DB.Begin()
				// 插入订单
				fmt.Println(orderres)
				order := models.Order{
					UserID:        orderres.UserID,
					OrderType:     "秒杀订单",
					Status:        "已创建",
					PaymentStatus: "unpaid",
					PayWay:        "",
					TotalAmount:   totalAmount, // 需要计算
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}
				if err := tx.Create(&order).Error; err != nil {
					tx.Rollback() // 回滚事务
					go WebsocketSendMessage(msg.CorrelationId, "生成订单失败，详情请联系客服。")
					go global.SendLogs("error", fmt.Sprintf("mysql 事务创建订单失败，userid:%d,eventID:%d,productID:%d", orderres.UserID, orderres.EventID, orderres.ProductID))
					// 丢弃消息
					msg.Ack(false)
					continue
				}
				// 生成订单item
				// 插入订单项
				orderItem := models.OrderItem{
					OrderID:      order.ID,
					ProductID:    orderres.ProductID,
					ProductName:  orderres.ProductName,
					Quantity:     orderres.Count,
					UnitPrice:    flashprice,   // 单个价格
					Subtotal:     totalAmount,  // 实际金额
					OriginAmount: originAmount, // 原始金额
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}

				if err := tx.Create(&orderItem).Error; err != nil {
					tx.Rollback() // 回滚事务
					go WebsocketSendMessage(msg.CorrelationId, "生成订单失败，详情请联系客服。")
					go global.SendLogs("error", fmt.Sprintf("mysql 事务创建订单详情失败，userid:%d,eventID:%d,productID:%d", orderres.UserID, orderres.EventID, orderres.ProductID))
					// 丢弃消息
					msg.Ack(false)
					continue
				}

				// 提交事务
				if err := tx.Commit().Error; err != nil {
					go WebsocketSendMessage(msg.CorrelationId, "生成订单失败，详情请联系客服。")
					go global.SendLogs("error", fmt.Sprintf("mysql 事务提交失败，userid:%d,eventID:%d,productID:%d", orderres.UserID, orderres.EventID, orderres.ProductID))
					// 丢弃消息
					msg.Ack(false)
					continue
				}

				// 将处理结果发送到待支付队列
				// 构建消息
				pendingpayment := PendingPayment{
					OrderID: order.ID,
					Uid:     utils.GenerateUniqueID(),
				}
				body, err := json.Marshal(pendingpayment)
				if err != nil {
					global.SendLogs("error", "订单处理结果消息构建失败", err)
					go WebsocketSendMessage(msg.CorrelationId, "获取支付异常，请稍后再试或手动提交")
					continue
				}
				// 发送消息
				err = global.RabbitMQ.Channel.PublishWithContext(
					context.Background(),         // 上下文
					global.RabbitMQ.Exchange,     // 交换机名称
					"flashevent.pay.order.queue", // 路由键
					false,                        // mandatory
					false,                        // immediate
					amqp091.Publishing{
						ContentType:     "application/json",                 // 内容类型
						ContentEncoding: "utf-8",                            // 内容编码
						DeliveryMode:    amqp091.Persistent,                 // 持久化消息
						Priority:        0,                                  // 优先级
						CorrelationId:   strconv.Itoa(int(orderres.UserID)), // 关联ID
						Expiration:      "3600000",                          // 消息过期时间，单位为毫秒
						Timestamp:       time.Now(),                         // 时间戳
						Type:            "order_flashevent_message",         // 消息类型
						AppId:           "order_flashevent",                 // 应用ID
						Body:            body,                               // 消息正文
					},
				)
				if err != nil {
					global.SendLogs("error", "订单处理结果消息发布失败", err)
					go WebsocketSendMessage(msg.CorrelationId, "获取支付异常，请稍后再试或联系客服")
					continue
				}
				go WebsocketSendMessage(msg.CorrelationId, "订单处理成功，正在获取支付信息.....")
				// 消息已经处理，丢弃
				msg.Ack(false)
			} else {
				// 消息已经处理，丢弃
				msg.Ack(false)
				continue
			}
			time.Sleep(50 * time.Millisecond) // 限制速率
		}
	}()
}

// 查询订单
func GetOrder(form *request.TakeFlashOrder) (order models.Order, err error) {
	return
}

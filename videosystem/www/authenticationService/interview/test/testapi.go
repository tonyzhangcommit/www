package test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// 注册10000个测试用户
type Resister struct {
	Name        string `form:"name" json:"name"`
	Password    string `form:"password" json:"password"`
	Phonenumber string `form:"phonenum" json:"phonenum"`
	VarifiCode  string `form:"varificode" json:"varificode"`
	AgentCode   string `form:"agentcode" json:"agentcode"`
}

// 用户登录请求结构体
type LoginPVC struct {
	Phonenumber      string `form:"phonenum" json:"phonenum" `
	VerificationCode string `form:"varificode" json:"varificode"`
}

// 秒杀下单请求结构体
type TakeFlashOrder struct {
	EventID   uint `json:"eventid"`
	ProductID uint `json:"peoductid"`
	UserID    uint `json:"userid"`
	Count     int  `json:"count"`
}

func RegistUsers(count int) {
	baseurl := "http://39.105.9.252:9999/api/consumer/register"

	for i := 0; i < count; i++ {
		tetailpnm := ""
		if i < 10 {
			tetailpnm = "0000" + strconv.Itoa(i)
		} else if i < 100 {
			tetailpnm = "000" + strconv.Itoa(i)
		} else if i < 1000 {
			tetailpnm = "00" + strconv.Itoa(i)
		} else if i < 10000 {
			tetailpnm = "0" + strconv.Itoa(i)
		} else {
			tetailpnm = strconv.Itoa(i)
		}

		form := Resister{
			Name:        "CustomUser" + strconv.Itoa(i),
			Password:    "Admin888.",
			Phonenumber: "185101" + tetailpnm,
			VarifiCode:  "000000",
		}
		body, err := json.Marshal(form)
		if err != nil {
			log.Fatal("序列化请求体失败")
		}
		// 构建请求

		req, err := http.NewRequest("POST", baseurl, bytes.NewBuffer(body))
		if err != nil {
			log.Fatal("构建请求失败")
		}
		// 设置请求头
		for key, value := range map[string]string{"Content-Type": "application/json"} {
			req.Header.Set(key, value)
		}
		// 发出请求
		client := http.Client{
			Timeout: 3 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			return
		}

		defer resp.Body.Close()
		// 读取响应体
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		var response Response
		// 反序列化结构体
		err = json.Unmarshal(responseBody, &response)
		if err != nil {
			log.Fatal("反序列化响应失败")
		}
		fmt.Println(response.ErrorCode, response.Data)
	}
}

// 获取测试用户登录jwttoken
func GetJwtToken(count int) {
	baseurl := "http://39.105.9.252:9999/api/consumer/login"
	var jwtmap = make(map[string]string, 0)
	for i := 0; i < count; i++ {
		// 构造手机号并且构造对应的userid
		// 注意：构造用户id时， 需要看测试用户id起始数量
		tetailpnm := ""
		userid := 104
		if i < 10 {
			tetailpnm = "0000" + strconv.Itoa(i)
		} else if i < 100 {
			tetailpnm = "000" + strconv.Itoa(i)
		} else if i < 1000 {
			tetailpnm = "00" + strconv.Itoa(i)
		} else if i < 10000 {
			tetailpnm = "0" + strconv.Itoa(i)
		} else {
			tetailpnm = strconv.Itoa(i)
		}
		userid = 14104 + i
		// 开始登录请求，获取jwttoken
		form := LoginPVC{
			Phonenumber:      "185101" + tetailpnm,
			VerificationCode: "000000",
		}

		body, err := json.Marshal(form)
		if err != nil {
			log.Fatal("序列化请求体失败")
		}

		// 构建请求
		req, err := http.NewRequest("POST", baseurl, bytes.NewBuffer(body))
		if err != nil {
			log.Fatal("构建请求失败")
		}
		// 设置请求头
		for key, value := range map[string]string{"Content-Type": "application/json"} {
			req.Header.Set(key, value)
		}

		// 发出请求
		client := http.Client{
			Timeout: 3 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			return
		}

		defer resp.Body.Close()
		// 读取响应体
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		var response LoginResponse
		// 反序列化结构体
		err = json.Unmarshal(responseBody, &response)
		if err != nil {
			log.Fatal("反序列化响应失败")
		}
		fmt.Println(response.Message, response.Data.AccessToken, userid)
		jwtmap[strconv.Itoa(userid)] = response.Data.AccessToken
	}
	file, err := os.OpenFile("C:\\Users\\Administrator\\Desktop\\debugvideosystem\\tools\\test.txt", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		log.Fatal("文件打开失败")
	}
	defer file.Close()
	for userid, jwttoken := range jwtmap {
		str := userid + "\t" + jwttoken + "\n"
		_, err = file.WriteString(str)
		if err != nil {
			return
		}
	}
}

// 秒杀活动下单测试
// 参数 count 为同时抢购人数
// 需要发送消息并监听websocket响应
func TestFlashOrder(count int) {
	baseurl := "http://39.105.9.252:9999/api/flashevent/takeorder"
	var alljwtmap = make(map[string]string, 0)
	file, err := os.Open("C:\\Users\\Administrator\\Desktop\\debugvideosystem\\tools\\test.txt")
	if err != nil {
		log.Fatalf("open file failed: %s \n", err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		countstrip := strings.Split(line, "\t")
		alljwtmap[countstrip[0]] = strings.Replace(countstrip[1], "\n", "", -1)
	}
	subjwtmap := make(map[string]string, 0)
	index := 0
	for k, v := range alljwtmap {
		if index < count {
			subjwtmap[k] = v
		} else {
			break
		}
		index++
	}
	fmt.Println(subjwtmap)

	var funclist []func()
	// 创建请求列表
	for userid, jwttoekn := range subjwtmap {
		intuserid, _ := strconv.Atoi(userid)
		form := TakeFlashOrder{
			EventID:   12,
			ProductID: 1,
			UserID:    uint(intuserid),
			Count:     1,
		}
		body, err := json.Marshal(form)
		if err != nil {
			log.Fatal("序列化请求体失败")
		}

		req, err := http.NewRequest("POST", baseurl, bytes.NewBuffer(body))
		if err != nil {
			log.Fatal("构建请求失败")
		}
		// 设置请求头
		for key, value := range map[string]string{"Content-Type": "application/json"} {
			req.Header.Set(key, value)
		}
		// 发出请求
		client := http.Client{
			Timeout: 3 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		// 读取响应体
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		var response TaskFlashOrder
		// 反序列化结构体
		err = json.Unmarshal(responseBody, &response)
		if err != nil {
			log.Fatal("反序列化响应失败")
		}
		// 获取初步请求结果
		firstReq := response.Data.Data
		_ = firstReq
		_ = jwttoekn
		_ = form
		// 创建websocket连接，并动态获取内容并保存
	}
	_ = funclist

	// 并发访问，设置每个客户端以纳秒为间隔区分
}

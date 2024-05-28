### **AuthService API V1 document**

#### **UserService** 
1. **获取验证码**
    - 请求地址：http://39.105.9.252:9999/api/consumer/getverifcode
    - 请求方式：POST
    - 数据格式：application/json
    - 参数列表
    ```json
        {
            "phonenum":"18510100000",
        }
    ```
    - 响应数据
    ```json
    {
        "errorCode": 0,
        "data": {
            "errorCode": 0,
            "data": "验证码已发送",
            "msg": "success"
        },
        "msg": "success"
    }
    ```
2. **客户端登录**
    - 请求地址：http://39.105.9.252:9999/api/consumer/login
    - 请求方式：POST
    - 数据格式：application/json
    - 参数列表
    ```json
    {
        "phonenum":"18510100000",
        "varificode":"292441"
    }
    ```

    - 响应数据
    ```json
    {
    "errorCode": 0,
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY0NjMwNDUsImp0aSI6IjI4MTA0IiwiaXNzIjoiYXBwIiwicm9sZXMiOlsibW9udGhseVZpcCJdfQ.Dwip1zsYw23svLWGA0zfYuN2ZyRDgyaXi8uPYoJ_uQA",
        "expires_in": 3000,
        "token_type": "bearer"
    },
    "msg": "success"
    }
    ```
3. **客户端退出登录**
    - 请求地址：http://39.105.9.252:9999/api/consumer/logout
    - 请求方式：POST
    - 数据格式：application/json
    - 参数列表： null
    - 认证头：Bearer Toekn:jwt toekn
    - 响应数据
    ```json
    {
        "errorCode": 0,
        "data": null,
        "msg": "success"
    }
    ```
4. **客户端注册**
    - 请求地址：http://39.105.9.252:9999/api/consumer/register
    - 请求方式：POST
    - 数据格式：application/json
    - 参数列表：
    ```json
    {
        "name":"customuser0001",
        "password":"Testpwd123.",
        "phonenum":"18212340000",
        "varificode":"672376",
        "agentcode":""  // 可选参数
    }
    ```
    - 认证头：null
    - 响应数据
    ```json
    {
    "errorCode": 0,
    "data": {
        "errorCode": 0,
        "data": {
            "AgentCode": "",
            "Children": null,
            "CreatedAt": "2024-05-24T10:12:35.621+08:00",
            "DeletedAt": null,
            "ID": 128102,
            "ParentID": 1,
            "UpdatedAt": "2024-05-24T10:12:35.621+08:00",
            "isbanned": false,
            "parentagentcode": "NKvVu1",
            "phonenumber": "18212340000",
            "roles": [
                {
                    "CreatedAt": "2024-03-11T19:51:24+08:00",
                    "DeletedAt": null,
                    "ID": 3,
                    "UpdatedAt": "2024-03-11T19:51:24+08:00",
                    "desc": "",
                    "rolename": "regularUser"
                }
            ],
            "username": "customuser0001"
        },
        "msg": "success"
    },
    "msg": "success"
    }
    ```
5. **客户端完善个人信息**
    - 请求地址：http://39.105.9.252:9999/api/consumer/login
    - 请求方式：POST
    - 数据格式：application/json
    - 参数列表： 
     ```json
     {
        "uid":128102,
        "address":"北京昌平",
        "sex":0,
        "identification":"410721000000000000",
        "email":"775200000@qq.com",
        "Preferences":""
     }
     ```
    - 认证头：Bearer Toekn:jwt toekn
    - 响应数据
    ```json
    {
        "errorCode": 0,
        "data": {
            "errorCode": 0,
            "data": "编辑/完善资料成功",
            "msg": "success"
        },
        "msg": "success"
    }
    ```
5. **客户端获取个人基本信息**
    - 请求地址：http://39.105.9.252:9999/api/consumer/login
    - 请求方式：GET
    - 数据格式：:
    - 参数列表： 
     ```
        phonenum=18212340000
        uid=128102  (二选一)
     ```
    - 认证头：Bearer Toekn:jwt toekn
    - 响应数据
    ```json
    {
    "errorCode": 0,
    "data": {
        "errorCode": 0,
        "data": {
            "address": "北京昌平",
            "agentcode": "",
            "email": "775200000@qq.com",
            "expvipdate": "2024-05-24T10:12:36+08:00",
            "idcard": "410721000000000000",
            "isbanned": false,
            "parentagentcode": "NKvVu1",
            "phonenumber": "18212340000",
            "preferences": "",
            "roles": [
                "regularUser"
            ],
            "sex": 0,
            "typevip": "月会员",
            "username": "customuser0001",
            "vip": false
        },
        "msg": "success"
    },
    "msg": "success"
    }
    ```
6. **客户端**
    - 请求地址：http://39.105.9.252:9999/api/consumer/login
    - 请求方式：POST
    - 数据格式：application/json
    - 参数列表： 
     ```json
     {
        
     }
     ```
    - 认证头：Bearer Toekn:jwt toekn
    - 响应数据
    ```json
    {
        "errorCode": 0,
        "data": null,
        "msg": "success"
    }
    ```
#### 待完成
- 查看订单
- 查看已购商品
- 开通会员
- 升级会员
- 商品管理相关（查看，筛选等）
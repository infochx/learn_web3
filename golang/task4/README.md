# 个人博客系统
这是一个使用 Go 语言、Gin 框架和 GORM 库开发的个人博客系统后端，支持用户认证、文章管理和评论功能。
## 功能特性

- 用户注册和登录（JWT 认证）
- 文章的创建、读取、更新和删除（CRUD）
- 评论功能
- 数据库使用 MySQL
## 技术栈

- Go 语言
- Gin Web 框架
- GORM ORM 库
- MySQL 数据库
- JWT 用户认证

## 项目结构

```
.
├── config/
│   └── config.go          # 数据库配置
├── api/
│   ├── auth.go            # 认证相关控制器
│   ├── user.go            # 用户相关控制器
│   ├── post.go            # 文章相关控制器
│   └── comment.go         # 评论相关控制器
├── middleware/
│   └── auth.go            # 认证中间件
├── models/
│   ├── user.go            # 用户模型
│   ├── post.go            # 文章模型
│   └── comment.go         # 评论模型
├── routes/
│   └── routes.go          # 路由定义
├── utils/
│   └── jwt.go             # JWT 工具
├── go.mod                 # Go 模块定义
├── go.sum                 # Go 模块校验和
├── main.go                # 程序入口
└── README.md              # 项目说明文档
```


## 安装和运行

### 环境要求

- Go 1.16 或更高版本
- MySQL 数据库


### 安装依赖
```
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
go get -u gorm.io/driver/mysql
go get -u github.com/gin-gonic/gin
go get -u github.com/dgrijalva/jwt-go
```

### 安装步骤

1. 克隆项目或在当前目录下创建文件

2. 安装依赖包：
   ```bash
   go mod tidy
   ```

3. 设置数据库环境变量：
   - DB_USER: 数据库用户名（默认: root）
   - DB_PASS: 数据库密码（默认: ）
   - DB_HOST: 数据库主机（默认: localhost）
   - DB_PORT: 数据库端口（默认: 3306）
   - DB_NAME: 数据库名称（默认: blog）

4. 运行程序：
   ```bash
   go run main.go
   ```
## API 接口说明

### 用户认证
- `POST /auth/register` - 用户注册
- `POST /auth/login` - 用户登录

### 用户管理
- `GET /users/:id` - 获取用户信息

### 文章管理
- `GET /posts/` - 获取所有文章
- `GET /posts/:id` - 获取单个文章
- `POST /posts/` - 创建文章
- `PUT /posts/:id` - 更新文章
- `DELETE /posts/:id` - 删除文章

### 评论管理
- `POST /comments/` - 创建评论
- `GET /comments/post/:postId` - 获取某篇文章的所有评论

## 数据库设计

### 用户表 (users)
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | uint | 主键 |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |
| deleted_at | time | 删除时间 |
| username | string | 用户名（唯一，非空） |
| password | string | 密码（非空） |
| email | string | 邮箱（唯一，非空） |

### 文章表 (posts)
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | uint | 主键 |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |
| deleted_at | time | 删除时间 |
| title | string | 标题（非空） |
| content | string | 内容（非空） |
| user_id | uint | 用户ID（外键） |

### 评论表 (comments)
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | uint | 主键 |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |
| deleted_at | time | 删除时间 |
| content | string | 内容（非空） |
| user_id | uint | 用户ID（外键） |
| post_id | uint | 文章ID（外键） |

## 测试

1. 用户注册

   ```bash 
    curl -X POST http://localhost:8080/auth/register \
         -H "Content-Type: application/json" \
         -d '{"username":"testuser","password":"123","email":"test@example.com"}'
   ```

   预期响应

   ```json
   {
       "message": "User registered successfully"
   }
   ```

2. 用户登录

   ```bash
    curl -X POST http://localhost:8080/auth/login \
         -H "Content-Type: application/json" \
         -d '{"username":"testuser","password":"123"}'
   ```

   预期响应

   ```json
   {
       "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzU2MDMwMjEwfQ.R1glu6C7Kl1YG8mRx-JpB9WfzG1B03F9frS0kKGD-VU"
   }
   ```

3. 获取用户信息

   ```bash 
   curl -X GET http://localhost:8080/users/1 \
         -H "Content-Type: application/json" \
         -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzU2MDMwMjEwfQ.R1glu6C7Kl1YG8mRx-JpB9WfzG1B03F9frS0kKGD-VU" 
   ```

   预期响应

   ```json
   {
       "ID": 1,
       "CreatedAt": "2025-08-24T17:14:18.859+08:00",
       "UpdatedAt": "2025-08-24T17:14:18.859+08:00",
       "DeletedAt": null,
       "Username": "testuser",
       "Password": "",
       "Email": "test@example.com"
   }
   ```

4. 获取所有文章

   ```bash 
   curl -X GET http://localhost:8080/post/ \
         -H "Content-Type: application/json" \
         -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzU2MDM1MDE3fQ.syjmsADlIbmB-aC1q7Mbi_UMRtf3GKEQ-4MS8TeuMLE" 
   ```

5. 获取单个文章   

   ```bash 
   curl -X GET http://localhost:8080/post/1 \
         -H "Content-Type: application/json" \
         -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzU2MTIwMDAzfQ.ZIwjVbBbTwZaBo278ZWICajHcpX8bW5nFXNmFHcPgJE" 
   ```

      预期响应

   ```json
   {
       "ID": 1,
       "CreatedAt": "2025-08-24T18:03:26.811+08:00",
       "UpdatedAt": "2025-08-24T18:25:32.924+08:00",
       "DeletedAt": null,
       "Title": "gin����3",
       "Content": "gin��������",
       "UserID": 1,
       "User": {
           "ID": 1,
           "CreatedAt": "2025-08-24T17:14:18.859+08:00",
           "UpdatedAt": "2025-08-24T17:14:18.859+08:00",
           "DeletedAt": null,
           "Username": "testuser",
           "Password": "$2a$10$pKlepE7mFhnvrweHZmcnPOEMmzZRytGAgVO0pciFVcTwgKfQfI3IW",
           "Email": "test@example.com"
       }
   } 
   ```

6. 创建文章

   ```bash 
   curl -X POST http://localhost:8080/post/ \
         -H "Content-Type: application/json;charset=utf-8" \
         -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzU2MTIwMDAzfQ.ZIwjVbBbTwZaBo278ZWICajHcpX8bW5nFXNmFHcPgJE" \
         -d '{"title":"gin测试","content":"gin测试内容"}'
   ```

      预期响应

   ```json
   {
       "ID": 1,
       "CreatedAt": "2025-08-24T18:03:26.8112559+08:00",
       "UpdatedAt": "2025-08-24T18:03:26.8112559+08:00",
       "DeletedAt": null,
       "Title": "gin����",
       "Content": "gin��������",
       "UserID": 1,
       "User": {
           "ID": 0,
           "CreatedAt": "0001-01-01T00:00:00Z",
           "UpdatedAt": "0001-01-01T00:00:00Z",
           "DeletedAt": null,
           "Username": "",
           "Password": "",
           "Email": ""
       }
   }
   ```

7. 更新文章

   ```bash 
   curl -X PUT http://localhost:8080/post/ \
         -H "Content-Type: application/json" \
         -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzU2MDMyMTQ0fQ.bMMuEsc4chwlAaXeY-xXfonVRUr_kZTslQGveCI6afc" \
         -d '{"ID":1,"title":"gin测试","content":"gin测试内容"}'
   ```

8. 删除文章

   ```bash
   curl -X DELETE http://localhost:8080/post/4 \
         -H "Content-Type: application/json" \
         -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzU2MDMzMDc1fQ.Zpcl1bb1g2XVIHD5JXzt2DSu8L8BuIg9G93G7YqqFFo" 
   ```

   预期响应

   ```json
   {"message":"deleted successfuly"}
   ```

10. 创建评论

    ```bash 
    curl -X POST http://localhost:8080/comments/ \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzU2MDMzMDc1fQ.Zpcl1bb1g2XVIHD5JXzt2DSu8L8BuIg9G93G7YqqFFo" \
          -d '{"content":"评论内容11","postId":1}'
    ```

11. 查询评论

    ```bash 
    curl -X GET http://localhost:8080/comments/post/1 \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzU2MDMzMDc1fQ.Zpcl1bb1g2XVIHD5JXzt2DSu8L8BuIg9G93G7YqqFFo"
    ```

























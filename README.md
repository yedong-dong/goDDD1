# Go DDD 后端服务

这是一个使用Gin框架搭建的简单后端服务，采用领域驱动设计(DDD)的思想组织代码结构，并封装了数据库等配置类。

## 项目结构

```
├── config/         # 配置相关代码
├── controllers/    # 控制器层，处理HTTP请求
├── models/         # 数据模型层
├── routes/         # 路由配置
├── services/       # 服务层，业务逻辑
├── .env            # 环境变量配置
├── .env.example    # 环境变量示例
├── go.mod          # Go模块文件
└── main.go         # 应用入口
```

## 功能特性

- 基于Gin的RESTful API
- GORM数据库ORM
- 环境变量配置
- 用户CRUD操作
- 分层架构设计

## 安装和运行

1. 克隆项目

```bash
git clone <项目地址>
cd goDDD1
```

2. 安装依赖

```bash
go mod tidy
```

3. 配置环境变量

复制`.env.example`文件为`.env`，并根据实际情况修改配置：

```bash
cp .env.example .env
# 编辑.env文件，设置数据库连接信息等
```

4. 运行应用

```bash
go run main.go
```

服务将在配置的端口上启动（默认为8080）。

## API接口

### 用户管理

- `POST /api/users` - 创建新用户
- `GET /api/users/:id` - 获取用户信息
- `PUT /api/users/:id` - 更新用户信息
- `DELETE /api/users/:id` - 删除用户

## 数据库配置

项目使用MySQL数据库，通过`.env`文件配置连接信息：

```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=goDDD1
```

## 扩展开发

1. 添加新模型：在`models`目录下创建新的模型文件
2. 添加新服务：在`services`目录下创建新的服务文件
3. 添加新控制器：在`controllers`目录下创建新的控制器文件
4. 配置新路由：在`routes/router.go`中添加新的路由
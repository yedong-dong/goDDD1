# Gin Context (*gin.Context) 常用API详解

## 1. 请求参数获取

### 路径参数
```go
// 获取路径参数，如 /user/:id 中的 id
id := ctx.Param("id")

// 获取所有路径参数
params := ctx.Params
```

### 查询参数
```go
// 获取查询参数，如 ?name=john
name := ctx.Query("name")

// 获取查询参数，带默认值
name := ctx.DefaultQuery("name", "default")

// 获取查询参数数组，如 ?tags=go&tags=web
tags := ctx.QueryArray("tags")

// 获取所有查询参数
queryMap := ctx.QueryMap("user")
```

### 表单数据
```go
// 获取表单字段
username := ctx.PostForm("username")

// 获取表单字段，带默认值
username := ctx.DefaultPostForm("username", "guest")

// 获取表单数组
hobbies := ctx.PostFormArray("hobbies")

// 获取表单映射
user := ctx.PostFormMap("user")
```

### 文件上传
```go
// 获取单个文件
file, err := ctx.FormFile("upload")

// 获取多个文件
form, err := ctx.MultipartForm()
files := form.File["upload[]"]

// 保存文件
err := ctx.SaveUploadedFile(file, "./uploads/"+file.Filename)
```

## 2. 请求体处理

### JSON绑定
```go
// 绑定JSON到结构体
var user User
if err := ctx.ShouldBindJSON(&user); err != nil {
    // 处理错误
}

// 必须绑定JSON（失败会自动返回400错误）
var user User
if err := ctx.BindJSON(&user); err != nil {
    return
}
```

### 其他格式绑定
```go
// 绑定XML
ctx.ShouldBindXML(&user)

// 绑定YAML
ctx.ShouldBindYAML(&user)

// 绑定表单
ctx.ShouldBind(&user)

// 绑定查询参数
ctx.ShouldBindQuery(&user)

// 绑定URI参数
ctx.ShouldBindUri(&user)

// 绑定Header
ctx.ShouldBindHeader(&user)
```

## 3. 响应处理

### JSON响应
```go
// 返回JSON
ctx.JSON(http.StatusOK, gin.H{
    "message": "success",
    "data": user,
})

// 返回带缩进的JSON
ctx.IndentedJSON(http.StatusOK, user)

// 返回安全的JSON（防止JSON劫持）
ctx.SecureJSON(http.StatusOK, user)

// 返回JSONP
ctx.JSONP(http.StatusOK, user)
```

### 其他格式响应
```go
// 返回XML
ctx.XML(http.StatusOK, user)

// 返回YAML
ctx.YAML(http.StatusOK, user)

// 返回纯文本
ctx.String(http.StatusOK, "Hello %s", name)

// 返回HTML
ctx.HTML(http.StatusOK, "index.html", gin.H{
    "title": "Main website",
})

// 返回数据流
ctx.Data(http.StatusOK, "application/octet-stream", []byte("some data"))
```

### 文件响应
```go
// 返回文件
ctx.File("./assets/image.png")

// 返回文件附件（下载）
ctx.FileAttachment("./assets/file.zip", "download.zip")

// 从文件系统返回文件
ctx.FileFromFS("assets/image.png", http.Dir("./public"))
```

## 4. 请求头操作

```go
// 获取请求头
contentType := ctx.GetHeader("Content-Type")
userAgent := ctx.Request.UserAgent()

// 设置响应头
ctx.Header("X-Custom-Header", "value")

// 设置Cookie
ctx.SetCookie("name", "value", 3600, "/", "localhost", false, true)

// 获取Cookie
cookie, err := ctx.Cookie("name")
```

## 5. 重定向

```go
// HTTP重定向
ctx.Redirect(http.StatusMovedPermanently, "https://google.com")

// 路由重定向
ctx.Request.URL.Path = "/new-path"
r.HandleContext(ctx)
```

## 6. 中间件和上下文

```go
// 设置键值对
ctx.Set("user", user)

// 获取值
value, exists := ctx.Get("user")
if exists {
    user := value.(User)
}

// 获取特定类型的值
user := ctx.MustGet("user").(User)

// 获取字符串值
username := ctx.GetString("username")

// 获取整数值
age := ctx.GetInt("age")

// 获取布尔值
isAdmin := ctx.GetBool("isAdmin")
```

## 7. 请求信息

```go
// 获取请求方法
method := ctx.Request.Method

// 获取请求URL
url := ctx.Request.URL.String()

// 获取客户端IP
clientIP := ctx.ClientIP()

// 获取请求路径
path := ctx.Request.URL.Path

// 获取完整URL
fullPath := ctx.FullPath()

// 检查请求是否为WebSocket升级
isWebSocket := ctx.IsWebsocket()
```

## 8. 错误处理

```go
// 添加错误
ctx.Error(err)

// 获取所有错误
errors := ctx.Errors

// 中止请求处理
ctx.Abort()

// 中止并返回状态码
ctx.AbortWithStatus(http.StatusUnauthorized)

// 中止并返回JSON错误
ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
    "error": "Invalid request",
})
```

## 9. 流处理

```go
// 服务器发送事件(SSE)
ctx.Stream(func(w io.Writer) bool {
    ctx.SSEvent("message", "Hello World")
    return false
})

// 写入数据流
ctx.Writer.Write([]byte("streaming data"))

// 刷新缓冲区
ctx.Writer.Flush()
```

## 10. 状态码设置

```go
// 设置状态码
ctx.Status(http.StatusCreated)

// 获取状态码
status := ctx.Writer.Status()
```

## 11. 实用工具方法

```go
// 检查请求是否已中止
if ctx.IsAborted() {
    return
}

// 获取请求的Content-Type
contentType := ctx.ContentType()

// 检查是否接受某种MIME类型
if ctx.NegotiateFormat(gin.MIMEJSON, gin.MIMEXML) == gin.MIMEJSON {
    // 客户端接受JSON
}

// 复制上下文（用于goroutine）
cCtx := ctx.Copy()
go func() {
    // 在goroutine中使用cCtx
}()
```

## 12. 常用HTTP状态码

```go
// 成功响应
http.StatusOK                    // 200
http.StatusCreated               // 201
http.StatusAccepted              // 202
http.StatusNoContent             // 204

// 重定向
http.StatusMovedPermanently      // 301
http.StatusFound                 // 302
http.StatusNotModified           // 304

// 客户端错误
http.StatusBadRequest            // 400
http.StatusUnauthorized          // 401
http.StatusForbidden             // 403
http.StatusNotFound              // 404
http.StatusMethodNotAllowed      // 405
http.StatusConflict              // 409
http.StatusUnprocessableEntity   // 422

// 服务器错误
http.StatusInternalServerError   // 500
http.StatusNotImplemented        // 501
http.StatusBadGateway            // 502
http.StatusServiceUnavailable    // 503
```

## 13. 常用示例

### 用户注册API
```go
func Register(ctx *gin.Context) {
    var user User
    if err := ctx.ShouldBindJSON(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid JSON format",
            "details": err.Error(),
        })
        return
    }
    
    // 业务逻辑处理
    if err := userService.CreateUser(&user); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create user",
        })
        return
    }
    
    ctx.JSON(http.StatusCreated, gin.H{
        "message": "User created successfully",
        "user": user,
    })
}
```

### 文件上传API
```go
func UploadFile(ctx *gin.Context) {
    file, err := ctx.FormFile("file")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "No file uploaded",
        })
        return
    }
    
    filename := filepath.Base(file.Filename)
    if err := ctx.SaveUploadedFile(file, "./uploads/"+filename); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to save file",
        })
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{
        "message": "File uploaded successfully",
        "filename": filename,
    })
}
```

### 分页查询API
```go
func GetUsers(ctx *gin.Context) {
    page := ctx.DefaultQuery("page", "1")
    limit := ctx.DefaultQuery("limit", "10")
    
    pageInt, _ := strconv.Atoi(page)
    limitInt, _ := strconv.Atoi(limit)
    
    users, total, err := userService.GetUsers(pageInt, limitInt)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get users",
        })
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{
        "data": users,
        "pagination": gin.H{
            "page": pageInt,
            "limit": limitInt,
            "total": total,
        },
    })
}
```

---

这些是Gin框架中`*gin.Context`最常用的API方法，涵盖了Web开发中的大部分场景。通过这些API，您可以轻松处理HTTP请求和响应，构建功能完整的Web应用程序。
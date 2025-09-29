package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 响应码常量
const (
	CodeSuccess     = "20000" // 成功
	CodeClientError = "40000" // 客户端错误
	CodeServerError = "50000" // 服务器错误
)

// Response 统一响应结构体
type Response struct {
	Code string      `json:"code"`
	Data interface{} `json:"data"`
}

// ResponseUtil 响应工具类
type ResponseUtil struct{}

// NewResponseUtil 创建响应工具类实例
func NewResponseUtil() *ResponseUtil {
	return &ResponseUtil{}
}

// Success 成功响应
// message: 成功消息
// data: 响应数据
func (r *ResponseUtil) Response(ctx *gin.Context, message string, data interface{}) {
	responseData := gin.H{
		"message": message,
	}

	if data != nil {
		responseData["data"] = data
	}

	ctx.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Data: responseData,
	})
}

// ClientError 客户端错误响应
// message: 错误消息
func (r *ResponseUtil) ResponseError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, Response{
		Code: CodeClientError,
		Data: gin.H{
			"message": message,
		},
	})
}

// ServerError 服务器错误响应
// err: 错误对象
func (r *ResponseUtil) ResponseServerError(ctx *gin.Context, err error) {
	message := "服务器内部错误"
	if err != nil {
		message = err.Error()
	}

	ctx.JSON(http.StatusOK, Response{
		Code: CodeServerError,
		Data: gin.H{
			"message": message,
		},
	})
}

// 全局响应工具实例
var ResponseUtilInstance = NewResponseUtil()

// 便捷函数，方便直接调用
func ResSuccess(ctx *gin.Context, message string, data interface{}) {
	ResponseUtilInstance.Response(ctx, message, data)
}

func ResClientError(ctx *gin.Context, message string) {
	ResponseUtilInstance.ResponseError(ctx, message)
}

func ResServerError(ctx *gin.Context, err error) {
	ResponseUtilInstance.ResponseServerError(ctx, err)
}

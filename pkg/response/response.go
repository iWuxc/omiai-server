package response

import (
	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/validator"
	"net/http"
)

type JSONResult struct {
	Code int         `json:"code" ` // 返回码，0表示成功，非0表示异常
	Msg  string      `json:"msg"`   // 返回提示信息
	Data interface{} `json:"data"`  // 返回数据
}

// ValidateError 表单验证失败处理函数
func ValidateError(c *gin.Context, err error, code Code) {
	ErrorResponse(c, code, validator.ValidateErr(err))
}

// ErrorResponse 失败响应
func ErrorResponse(c *gin.Context, code Code, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	})
}

// ErrorResponseWithData 失败响应
func ErrorResponseWithData(c *gin.Context, code Code, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  msg,
		"data": data,
	})
}

// MiddlewareErrorResponse 中间件失败响应
func MiddlewareErrorResponse(c *gin.Context, code Code, msg string) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	})
}

func WechatNotifySuccessResponse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": "SUCCESS",
	})
}
func WechatNotifyFailResponse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": "FAIL",
	})
}

func AlipayNotifySuccessResponse(c *gin.Context) {
	c.String(http.StatusOK, "%s", "success")
}
func AlipayNotifyFailResponse(c *gin.Context) {
	c.String(http.StatusOK, "%s", "fail")
}

// StreamErrorResponse 流式响应错误返回
func StreamErrorResponse(c *gin.Context, msg interface{}) {
	c.SSEvent("error", msg)
}

// StreamSuccessResponse 流式响应成功返回
func StreamSuccessResponse(c *gin.Context, data interface{}) {
	c.SSEvent("message", data)
}

// StreamCloseResponse 流式响应结束返回
func StreamCloseResponse(c *gin.Context, data interface{}) {
	c.SSEvent("close", data)
}

// StreamFilterResponse 内容触发过滤机制
func StreamFilterResponse(c *gin.Context, data interface{}) {
	c.SSEvent("filter", data)
}

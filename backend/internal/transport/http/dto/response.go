package dto

// 响应码定义
const (
	CodeSuccess       = 0    // 成功
	CodeInvalidParams = 1001 // 参数错误
	CodeUnauthorized  = 1002 // 未授权
	CodeForbidden     = 1003 // 禁止访问
	CodeNotFound      = 1004 // 资源不存在
	CodeBusinessError = 2001 // 业务错误
	CodeInternalError = 9999 // 内部错误
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`    // 响应码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data,omitempty"`    // 响应数据
}

// PageResponse 分页响应结构
type PageResponse struct {
	Total int64       `json:"total"` // 总数
	Items interface{} `json:"items"` // 列表项
}
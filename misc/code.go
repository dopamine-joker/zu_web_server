package misc

// 自定义错误码
const (
	CodeSuccess      = 0
	CodeFail         = 1
	CodeUnknownError = -1
	CodeTokenError   = 400
)

//MsgCodeMap 默认错误码对应信息
var MsgCodeMap = map[int]string{
	CodeSuccess:      "success",
	CodeFail:         "fail",
	CodeUnknownError: "unknown error",
	CodeTokenError:   "Token error",
}

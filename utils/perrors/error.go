// 扩展错误码
package perrors

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	emptyStruct interface{} = nil

	SUCCESS_CODE = ExtraCode(0, "success", emptyStruct)
	ERROR_CODE   = ExtraCode(-1, "failed", emptyStruct)
	ERROR_1001   = ExtraCode(1001, "签名失效", emptyStruct)
	ERROR_1002   = ExtraCode(1002, "签名错误", emptyStruct)
	ERROR_2001   = ExtraCode(2001, "数据记录不存在", emptyStruct)
	ERROR_2002   = ExtraCode(2002, "数据重复", emptyStruct)
	ERROR_2003   = ExtraCode(2003, "创建/更新数据失败", emptyStruct)
	ERROR_3000   = ExtraCode(3000, "参数验证不通过", emptyStruct)
	ERROR_3001   = ExtraCode(3001, "操作过于频繁,请稍后再试", emptyStruct)
	ERROR_3002   = ExtraCode(3002, "无操作权限", emptyStruct)
	ERROR_3003   = ExtraCode(3003, "上传失败", emptyStruct)
	ERROR_3004   = ExtraCode(3004, "数据格式不正确", emptyStruct)
	ERROR_3005   = ExtraCode(3005, "提交的数据不符合字典约束范围值", emptyStruct)
	ERROR_3006   = ExtraCode(3006, "提交的数据校验不通过,验证失败", emptyStruct)
	ERROR_3054   = ExtraCode(3054, "系统繁忙,请稍后再试", emptyStruct)
	ERROR_4001   = ExtraCode(4001, "账户未授权", emptyStruct)
	ERROR_4002   = ExtraCode(4002, "设备登录已达上限", emptyStruct)
	ERROR_4003   = ExtraCode(4003, "禁止访问", emptyStruct)
	ERROR_4004   = ExtraCode(4004, "页面未定义", emptyStruct)
	ERROR_4100   = ExtraCode(4100, "未知错误", emptyStruct)
	ERROR_5000   = ExtraCode(5000, "内部服务异常", emptyStruct)
	ERROR_5001   = ExtraCode(5001, "系统维护中", emptyStruct)
	ERROR_9000   = ExtraCode(9000, "账户授权Token值无效", emptyStruct)
	ERROR_9001   = ExtraCode(9001, "账户异常,请联系管理员", emptyStruct)
	ERROR_9002   = ExtraCode(9002, "账号/密码错误,请检查后重试", emptyStruct)
	ERROR_9003   = ExtraCode(9003, "账号已存在", emptyStruct)
	ERROR_9005   = ExtraCode(9005, "密码错误次数过多,请稍后再试", emptyStruct)
	ERROR_9006   = ExtraCode(9006, "账户信息异常", emptyStruct)
	ERROR_9007   = ExtraCode(9007, "身份信息验证不通过", emptyStruct)
	ERROR_9008   = ExtraCode(9008, "登录失败", emptyStruct)
	ERROR_9009   = ExtraCode(9009, "当前手机号未绑定账号或账号不可用", emptyStruct)
)

// code: int 错误码
// message: string 错误信息
// detail: interface{} 返回数据
func ExtraCode(code int, message string, detail interface{}) gcode.Code {
	return gcode.New(code, message, detail)
}

// 用于业务逻辑中主动抛出错误
func Throw(message string) error {
	return gerror.NewCode(ERROR_CODE, message)
}

// 用于业务逻辑中主动抛出格式化错误
func Throwf(format string, args ...interface{}) error {
	return gerror.NewCodef(ERROR_CODE, format, args...)
}

// 用于业务逻辑中主动抛出错误，并指定错误码
func ThrowCode(code gcode.Code, message string) error {
	return gerror.NewCode(code, message)
}

// 用于业务逻辑中主动抛出格式化错误
func ThrowCodef(code gcode.Code, format string, args ...interface{}) error {
	return gerror.NewCodef(code, format, args...)
}

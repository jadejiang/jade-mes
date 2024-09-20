package ecode

var (
	// ErrUnhandledException ...
	ErrUnhandledException = func(err error) ECode { return newError(500, "未知错误", err) }
	// ErrParseJSONException ...
	ErrParseJSONException = func(err error) ECode { return newError(501, "解析JSON数据出错", err) }
	// ErrDBCASLimitExceeded ...
	ErrDBCASLimitExceeded = newError(502, "乐观锁重试次数达上限")

	ErrInvalidStateForWasher = newError(1001, "洗衣机状态错误")

	ErrInvalidDeviceStatus = newError(1002, "设备状态错误")
)

package hook

type FailHandler func(c IHttpContext, err error) error // 返回nil时能够从致命错误中恢复
type ErrorHandler func(c IHookContextRead, err error, isDeadly bool)
type BeforeHandle func(c IHttpContext) (error, error)
type AfterHandle func(c IHttpContext) (error, error)

package hook

type FailHandler func(c IHttpContext, err error) error // 返回nil时能够从致命错误中恢复
type ErrorHandler func(c IHookContextRead, err error, isDeadly bool)
type BeforeHandle func(c IHttpContext) (error, error)
type AfterHandle func(c IHttpContext) (error, error)

//
// type FailHandlerMap struct {
// 	m sync.Map
// }
//
// func NewFailHandlerMap() *FailHandlerMap {
// 	return &FailHandlerMap{sync.Map{}}
// }
//
// func (fhm *FailHandlerMap) Add(fh FailHandler) {
// 	p := reflect.ValueOf(fh).Pointer()
// 	fhm.m.Store(p, fh)
// }
//
// func (fhm *FailHandlerMap) Del(fh FailHandler) {
// 	p := reflect.ValueOf(fh).Pointer()
// 	fhm.m.Delete(p)
// }
//
// func (fhm *FailHandlerMap) InvokeAll(c *HttpContext, err error) {
// 	fhm.m.Range(func(_, value interface{}) bool {
// 		fh := value.(FailHandler)
// 		fh(c, err)
// 		return true
// 	})
// }
//
// type ErrorHandlerMap struct {
// 	m sync.Map
// }
//
// func NewErrorHandlerMap() *ErrorHandlerMap {
// 	return &ErrorHandlerMap{sync.Map{}}
// }
//
// func (fhm *ErrorHandlerMap) Add(fh ErrorHandler) {
// 	p := reflect.ValueOf(fh).Pointer()
// 	fhm.m.Store(p, fh)
// }
//
// func (fhm *ErrorHandlerMap) Del(fh ErrorHandler) {
// 	p := reflect.ValueOf(fh).Pointer()
// 	fhm.m.Delete(p)
// }
//
// func (fhm *ErrorHandlerMap) InvokeAll(c *HttpContext, err error, isDeadly bool) {
// 	fhm.m.Range(func(_, value interface{}) bool {
// 		fh := value.(ErrorHandler)
// 		fh(c, err, isDeadly)
// 		return true
// 	})
// }
//
// type BeforeHandleMap struct {
// 	m sync.Map
// }
//
// func NewBeforeHandleMap() *BeforeHandleMap {
// 	return &BeforeHandleMap{sync.Map{}}
// }
//
// func (bhm *BeforeHandleMap) Add(bh BeforeHandle) {
// 	p := reflect.ValueOf(bh).Pointer()
// 	bhm.m.Store(p, bh)
// }
//
// func (bhm *BeforeHandleMap) Del(bh BeforeHandle) {
// 	p := reflect.ValueOf(bh).Pointer()
// 	bhm.m.Delete(p)
// }
//
// type AfterHandleMap struct {
// 	m sync.Map
// }
//
// func NewAfterHandleMap() *AfterHandleMap {
// 	return &AfterHandleMap{sync.Map{}}
// }
//
// func (ahm *AfterHandleMap) Add(ah AfterHandle) {
// 	p := reflect.ValueOf(ah).Pointer()
// 	ahm.m.Store(p, ah)
// }
//
// func (ahm *AfterHandleMap) Del(ah AfterHandle) {
// 	p := reflect.ValueOf(ah).Pointer()
// 	ahm.m.Delete(p)
// }

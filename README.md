# go-proxy
proxy lib for golang
## How to use
'go get it': go get https://github.com/crzbird/go-proxy

define a proxy struct for target struct,example:
target struct:
```
func (bean *Bean) Add(a int, b string, c interface{}) *Result {
	log.Println("I am invoking!!", a, b, c)
	//do bussiness...
	res := &Result{
		Code: a,
		Msg:  b,
		Data: c,
	}
	return res
}

type Bean struct {
}
```
proxy struct:
```
type BeanProxy struct {
	*Bean
	Add func(a int, b string, c interface{}) *Result
}
```
now you can build a proxy:
```
beanProxy := &BeanProxy{Bean: &Bean{}}
	proxyInfo := &ProxyInfo{
		MethodName: "Add",
		Before: func(a int, b string, c interface{}) interface{} {
			log.Println("before....")
			return nil
		},
		After: func(a int, b string, c interface{}) interface{} {
			log.Println("after....")
			return nil
		},
	}
	Create(beanProxy, []*ProxyInfo{proxyInfo})
	result := beanProxy.Add(1, "a", "c")
	log.Println(result)
```
about ProxyInfo:
```
type ProxyInfo struct {
	MethodName string //target method
	Before     interface{} //it must be a func type ,invoke before bussiness method
	After      interface{} //it must be a func type ,invoke after bussiness method
}
```
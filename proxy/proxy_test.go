package proxy

import (
	"fmt"
	"log"
	"testing"
)

func TestProxy(t *testing.T) {
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
	err := Create(beanProxy, []*ProxyInfo{proxyInfo})
	fmt.Println(err)
	result := beanProxy.Add(1, "a", "c")
	log.Println(result)
}

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

type Result struct {
	Code int
	Msg  string
	Data interface{}
}

type BeanProxy struct {
	*Bean
	Add func(a int, b string, c interface{}) *Result
}

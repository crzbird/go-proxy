package proxy

import (
	"errors"
	"reflect"
)

type ProxyInfo struct {
	MethodName string
	Before     interface{}
	After      interface{}
}

func Create(proxyBean interface{}, proxyInfos []*ProxyInfo) {

	bpv := reflect.ValueOf(proxyBean)
	bpe := bpv.Elem()
	bpt := bpe.Type()
	fieldNum := bpt.NumField()
	beans := []reflect.Value{}
	for i := 0; i < fieldNum; i++ {
		fieldV := bpe.Field(i)
		if fieldV.Kind() == reflect.Ptr && fieldV.Elem().Kind() == reflect.Struct {
			beans = append(beans, fieldV)
		}

	}
	for _, bean := range beans {
		for _, proxyInfo := range proxyInfos {
			beanMethod := bean.MethodByName(proxyInfo.MethodName)
			if !beanMethod.IsValid() {
				continue
			}
			targetField := bpe.FieldByName(proxyInfo.MethodName)
			if targetField.Kind() == reflect.Func && targetField.IsValid() && targetField.CanSet() {
				proxyFunc := reflect.MakeFunc(targetField.Type(), func(args []reflect.Value) (results []reflect.Value) {
					reflect.ValueOf(proxyInfo.Before).Call(args)
					res, err := InvokeMethod(bean, proxyInfo.MethodName, args)
					if err != nil {
						return []reflect.Value{}
					}
					reflect.ValueOf(proxyInfo.After).Call(args)
					return res
				})
				targetField.Set(proxyFunc)
			}
		}
	}
}

func InvokeMethod(bean reflect.Value, methodName string, args []reflect.Value) ([]reflect.Value, error) {
	beanType := reflect.TypeOf(bean)
	if beanType.Kind() != reflect.Struct {
		return []reflect.Value{}, errors.New("proxy bean must be a struct")
	}
	targetMethod := bean.MethodByName(methodName)
	if !targetMethod.IsValid() || targetMethod.IsNil() || targetMethod.IsZero() {
		return []reflect.Value{}, errors.New("proxy bean must contains proxy method")
	}
	return targetMethod.Call(args), nil
}

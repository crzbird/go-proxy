package proxy

import (
	"errors"
	"fmt"
	"reflect"
)

type ProxyInfo struct {
	MethodName string      //target method
	Before     interface{} //it must be a func type ,invoke before bussiness method
	After      interface{} //it must be a func type ,invoke after bussiness method
}

func Create(proxyBean interface{}, proxyInfos []*ProxyInfo) error {
	var err error
	defer func() {
		if pErr := recover(); pErr != nil {
			err = errors.New(fmt.Sprintf("%v", pErr))
		}
	}()
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
			beforeMethodV := reflect.ValueOf(proxyInfo.Before)
			afterMethodV := reflect.ValueOf(proxyInfo.After)
			if !beforeMethodV.IsValid() || beforeMethodV.Kind() != reflect.Func || !afterMethodV.IsValid() || afterMethodV.Kind() != reflect.Func {
				return errors.New("illegal proxy method")
			}
			targetField := bpe.FieldByName(proxyInfo.MethodName)
			if targetField.Kind() == reflect.Func && targetField.IsValid() && targetField.CanSet() {
				proxyFunc := reflect.MakeFunc(targetField.Type(), func(args []reflect.Value) (results []reflect.Value) {
					beforeMethodV.Call(args)
					res, err := InvokeMethod(bean, proxyInfo.MethodName, args)
					if err != nil {
						return []reflect.Value{}
					}
					afterMethodV.Call(args)
					return res
				})
				targetField.Set(proxyFunc)
			}
		}
	}
	return err
}

func InvokeMethod(bean reflect.Value, methodName string, args []reflect.Value) ([]reflect.Value, error) {
	var err error
	defer func() {
		if pErr := recover(); pErr != nil {
			err = errors.New(fmt.Sprintf("%v", pErr))
		}
	}()
	beanType := reflect.TypeOf(bean)
	if beanType.Kind() != reflect.Struct {
		return []reflect.Value{}, errors.New("proxy bean must be a struct")
	}
	targetMethod := bean.MethodByName(methodName)
	if !targetMethod.IsValid() || targetMethod.IsNil() || targetMethod.IsZero() {
		return []reflect.Value{}, errors.New("proxy bean must contains proxy method")
	}
	return targetMethod.Call(args), err
}

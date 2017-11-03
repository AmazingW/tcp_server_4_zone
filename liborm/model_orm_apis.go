package liborm

import (
	"fmt"
	"reflect"

	"github.com/wuqifei/server_lib/libmodel"
)

func (orm *Orm) registerModel(tablename string, model interface{}, tags []string) {
	reflectVal := reflect.ValueOf(model)
	reflectType := reflect.Indirect(reflectVal).Type()

	if reflectVal.Kind() != reflect.Ptr && reflectVal.Kind() != reflect.Struct {
		panic(fmt.Errorf("cannot register non ptr or struct model:[%s],[%v]", libmodel.GetObjFullName(reflectType), reflectVal.Kind()))
	}

	// 如果不是结构体，直接error
	if reflectType.Kind() != reflect.Struct {
		panic(fmt.Errorf("cannot register non struct model or ptr model:[%s] type:[%s]", libmodel.GetObjFullName(reflectType), reflectType.String()))
	}

	if len(tablename) == 0 {
		objName := reflectType.Name()
		tablename = libmodel.ObjName2SqlName(objName)
	}

	// 判断重复注册
	if orm.modelCache.Get(tablename) != nil {
		panic(fmt.Errorf("register repeat please check again[%s]", tablename))
	}

	m := newModelTableInfo(tablename, reflectVal, tags)

	orm.modelCache.Set(m.Name, m)

}

func (orm *Orm) getModelInfoAndIndtype(model interface{}) (*ModelTableInfo, reflect.Value) {
	val := reflect.ValueOf(model)
	ind := reflect.Indirect(val)

	if ind.Type().Kind() == reflect.Ptr {
		panic("orm : model not support ** model")
	}

	name := ind.Type().Name()

	m := orm.modelCache.Get(name)
	if m == nil {
		panic(fmt.Errorf("model has called a invalid name :[%s]", name))
	}

	return m.(*ModelTableInfo), ind
}

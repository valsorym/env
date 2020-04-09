package env

import (
	"reflect"
)

// instance is an auxiliary structure for performing reflection.
type instance struct {
	Ptr   reflect.Value
	Type  reflect.Type
	Kind  reflect.Kind
	Value reflect.Value

	IsPtr    bool
	IsStruct bool
	IsValid  bool
}

// Init defines the main reflect's parameters.
func (inst *instance) Init(obj interface{}) {
	inst.Type = reflect.TypeOf(obj)
	inst.Value = reflect.ValueOf(obj)
	inst.Kind = inst.Type.Kind()

	if inst.Kind == reflect.Ptr {
		inst.IsPtr = true
		inst.Ptr = inst.Value
		inst.Type = inst.Type.Elem()
		inst.Kind = inst.Type.Kind()
		inst.Value = inst.Value.Elem()
	} else {
		inst.Ptr = reflect.New(inst.Type)
		tmp := inst.Ptr.Elem()
		tmp.Set(inst.Value)
	}

	inst.IsStruct = inst.Kind == reflect.Struct
	inst.IsValid = inst.Value.IsValid()
}

// Implements returns true if instance implements interface.
//
// Usage:
//	if inst.Implements((*CustomInterface)(nil)) { ... }
func (inst *instance) Implements(ifc interface{}) bool {
	return inst.Ptr.Type().Implements(reflect.TypeOf(ifc).Elem())
}

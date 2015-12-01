package context

import (
	"fmt"
	"reflect"
)

type Invoker struct {
	values map[reflect.Type]reflect.Value
}

// CreateInvoker return standard reflector to invoke a method
func CreateInvoker() *Invoker {
	return &Invoker{
		values: make(map[reflect.Type]reflect.Value),
	}
}

//// MARK: Struct's constructors
//func (inj *Invoker) Apply(val interface{}) error {
//	v := reflect.ValueOf(val)

//	for v.Kind() == reflect.Ptr {
//		v = v.Elem()
//	}

//	if v.Kind() != reflect.Struct {
//		return nil // Should not panic here ?
//	}

//	t := v.Type()

//	for i := 0; i < v.NumField(); i++ {
//		f := v.Field(i)
//		structField := t.Field(i)
//		if f.CanSet() && (structField.Tag == "inject" || structField.Tag.Get("inject") != "") {
//			ft := f.Type()
//			v := inj.Get(ft)
//			if !v.IsValid() {
//				return fmt.Errorf("Value not found for type %v", ft)
//			}

//			f.Set(v)
//		}

//	}

//	return nil
//}

func (i *Invoker) Invoke(function interface{}) ([]reflect.Value, error) {
	/* Condition validation: Input must be a function type */
	reflectFunction := reflect.TypeOf(function)
	if reflectFunction.Kind() != reflect.Func {
		return nil, fmt.Errorf("Input is not a function type.")
	}

	input := make([]reflect.Value, reflectFunction.NumIn())
	for idx := 0; idx < reflectFunction.NumIn(); idx++ {
		argument := reflectFunction.In(idx)
		value := i.Get(argument)
		input[idx] = value
	}
	return reflect.ValueOf(function).Call(input), nil
}

// InterfaceOf dereferences a pointer to an Interface type.
// It panics if value is not an pointer to an interface.
func InterfaceOf(value interface{}) reflect.Type {
	t := reflect.TypeOf(value)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Interface {
		panic("Called inject.InterfaceOf with a value that is not a pointer to an interface. (*MyInterface)(nil)")
	}

	return t
}

// Maps the concrete value of val to its dynamic type using reflect.TypeOf,
// It returns the TypeMapper registered in.
func (i *Invoker) Map(val interface{}) *Invoker {
	i.values[reflect.TypeOf(val)] = reflect.ValueOf(val)
	return i
}

func (i *Invoker) MapTo(val interface{}, ifacePtr interface{}) *Invoker {
	i.values[InterfaceOf(ifacePtr)] = reflect.ValueOf(val)
	return i
}

// Maps the given reflect.Type to the given reflect.Value and returns
// the Typemapper the mapping has been registered in.
func (i *Invoker) Set(typ reflect.Type, val reflect.Value) *Invoker {
	i.values[typ] = val
	return i
}

func (i *Invoker) Get(t reflect.Type) reflect.Value {
	val := i.values[t]

	// reflect.Array
	// reflect.Chan		// Channel
	// reflect.Interface
	// reflect.Map
	// reflect.Ptr
	// reflect.Slice
	// reflect.String
	// reflect.Struct

	if val.IsValid() {
		return val

	} else {
		// switch t.Kind() {

		// case reflect.Complex64, reflect.Complex128, reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// 	val = reflect.ValueOf(0)
		// 	break

		// case reflect.Bool:
		// 	val = reflect.ValueOf(false)
		// 	break

		// case reflect.String:
		// 	val = reflect.ValueOf("")
		// 	break
		// }

		if t.Kind() == reflect.Bool {
			val = reflect.ValueOf(false)
		} else if t.Kind() == reflect.Complex64 || t.Kind() == reflect.Complex128 || t.Kind() == reflect.Float32 || t.Kind() == reflect.Float64 || t.Kind() == reflect.Int || t.Kind() == reflect.Int8 || t.Kind() == reflect.Int16 || t.Kind() == reflect.Int32 || t.Kind() == reflect.Int64 || t.Kind() == reflect.Uint || t.Kind() == reflect.Uint8 || t.Kind() == reflect.Uint16 || t.Kind() == reflect.Uint32 || t.Kind() == reflect.Uint64 {
			val = reflect.ValueOf(0)
		} else if t.Kind() == reflect.String {
			val = reflect.ValueOf("")
		} else if t.Kind() == reflect.Array {
			val = reflect.MakeSlice(t, 0, 0)
		} else if t.Kind() == reflect.Map {
			val = reflect.MakeMap(t)
		} else if t.Kind() == reflect.Interface {
			for k, v := range i.values {
				if k.Implements(t) {
					val = v
					break
				}
			}
		} else {
			val = reflect.ValueOf(nil)
		}
	}
	return val
}

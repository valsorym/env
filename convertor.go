package env

// /// import (
// /// 	"fmt"
// /// 	"net/url"
// /// 	"reflect"
// /// )
// ///
// /// // Error's values.
// /// var (
// /// 	InitializedError = fmt.Errorf("object must be initialized")
// /// 	ObjectError      = fmt.Errorf("object must be a structure")
// /// 	//TypeError        = fmt.Errorf("incorrect type")
// /// )
// ///
// /// // Marshaller describes an interface for implementing
// /// // a custom method for marshaling.
// /// type Marshaller interface {
// /// 	MarshalENV() ([]string, error)
// /// }
// ///
// /// /*
// /// // Unmarshaller describes an interface for implementing
// /// // a custom method for unmarshaling.
// /// type Unmarshaller interface {
// /// 	UnmarshalENV() error
// /// }
// /// */
// ///
// /// // marshalENVOld saves scope into environment data.
// /// //
// /// // Supported types: int, int8, int16, int32, int64, uint, uint8, uint16,
// /// // uint32, uint64, bool, float32, float64, string, and slice from thous types.
// /// func marshalENVOld(scope interface{}, prefix string) ([]string, error) {
// /// 	var (
// /// 		rt reflect.Type  // type
// /// 		rv reflect.Value // value
// /// 		rp reflect.Value // pointer
// ///
// /// 		err    error
// /// 		result []string
// /// 	)
// ///
// /// 	// Define: type, value and pointer.
// /// 	rt = reflect.TypeOf(scope)
// /// 	rv = reflect.ValueOf(scope)
// /// 	if rt.Kind() == reflect.Ptr {
// /// 		rp, rt, rv = rv, rt.Elem(), rv.Elem()
// /// 	} else {
// /// 		rp = reflect.New(reflect.TypeOf(scope))
// /// 		tmp := rp.Elem()
// /// 		tmp.Set(rv)
// /// 	}
// ///
// /// 	// Scope validation.
// /// 	switch {
// /// 	case rt.Kind() != reflect.Struct:
// /// 		return result, ObjectError
// /// 	case !rv.IsValid():
// /// 		return result, InitializedError
// /// 	}
// ///
// /// 	// Implements Marshaler interface.
// /// 	if rp.Type().Implements(reflect.TypeOf((*Marshaller)(nil)).Elem()) {
// /// 		// Try to run custom MarshalENV function.
// /// 		if m := rp.MethodByName("MarshalENV"); m.IsValid() {
// /// 			tmp := m.Call([]reflect.Value{})
// /// 			value := tmp[0].Interface()
// /// 			err := tmp[1].Interface()
// /// 			if err != nil {
// /// 				return []string{}, fmt.Errorf("marshal: %v", err)
// /// 			}
// /// 			return value.([]string), nil
// /// 		}
// /// 	}
// ///
// /// 	// Walk through the fields.
// /// 	result = make([]string, 0, rv.NumField()-1)
// /// 	for i := 0; i < rv.NumField(); i++ {
// /// 		var key, value, sep string
// /// 		field := rv.Type().Field(i)
// /// 		instance := rv.FieldByName(field.Name)
// ///
// /// 		key, sep = parseTag(field.Tag.Get("env"), field.Name, " ")
// /// 		kind := field.Type.Kind()
// /// 		switch kind {
// /// 		case reflect.Int, reflect.Int8, reflect.Int16,
// /// 			reflect.Int32, reflect.Int64:
// /// 			value = fmt.Sprintf("%d", instance.Int())
// /// 		case reflect.Uint, reflect.Uint8, reflect.Uint16,
// /// 			reflect.Uint32, reflect.Uint64:
// /// 			value = fmt.Sprintf("%d", instance.Uint())
// /// 		case reflect.Float32, reflect.Float64:
// /// 			value = fmt.Sprintf("%f", instance.Float())
// /// 		case reflect.Bool:
// /// 			value = fmt.Sprintf("%t", instance.Bool())
// /// 		case reflect.String:
// /// 			value = fmt.Sprintf("%s", instance.String())
// /// 		case reflect.Array:
// /// 			value, err = getSequence(&instance, sep)
// /// 			if err != nil {
// /// 				return result, err
// /// 			}
// /// 		case reflect.Slice:
// /// 			value, err = getSequence(&instance, sep)
// /// 			if err != nil {
// /// 				return result, err
// /// 			}
// /// 		case reflect.Ptr:
// /// 			instance = instance.Elem()
// /// 			if instance.Kind() != reflect.Struct {
// /// 				return result, TypeError
// /// 			}
// /// 			fallthrough
// /// 		case reflect.Struct:
// /// 			if u, ok := instance.Interface().(url.URL); ok {
// /// 				// Type of url.URL.
// /// 				value = u.String()
// /// 			} else {
// /// 				// Recursive parsing
// /// 				p := fmt.Sprintf("%s%s_", prefix, key)
// /// 				value, err := marshalENVOld(instance.Interface(), p)
// /// 				if err != nil {
// /// 					return result, err
// /// 				}
// ///
// /// 				// Expand the result.
// /// 				for _, v := range value {
// /// 					result = append(result, v)
// /// 				}
// /// 				continue // object doesn't save
// /// 			}
// /// 		default:
// /// 			return result, TypeError
// /// 		} // switch
// ///
// /// 		// Set into environment and add to result list.
// /// 		key = fmt.Sprintf("%s%s", prefix, key)
// /// 		Set(key, value)
// /// 		result = append(result, fmt.Sprintf("%s=%s", key, value))
// /// 	} // for
// ///
// /// 	return result, nil
// /// }

/*
// unmarshalENVOld gets variables from the environment and sets them by
// pointer into scope. Returns an error if something went wrong.
//
// Supported types: int, int8, int16, int32, int64, uint, uint8, uint16,
// uint32, uint64, bool, float32, float64, string, and slice from thous types.
func unmarshalENVOld(scope interface{}, prefix string) error {
	var (
		rt reflect.Type  // type
		rv reflect.Value // value
		rp reflect.Value // pointer
	)

	// Define: type, value and pointer.
	rt = reflect.TypeOf(scope)
	rv = reflect.ValueOf(scope)

	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot use scope (type %s) as type *%s "+
			"in argument to decode", rt, rt)
	} else if rv.IsNil() {
		return fmt.Errorf("null element")
	}

	rp, rt, rv = rv, rt.Elem(), rv.Elem()

	return loadData(rt, rv, rp, prefix)
}

// loadData sets data inside an object.
func loadData(rt reflect.Type, rv, rp reflect.Value, prefix string) error {
	// Scope validation.
	switch {
	case rt.Kind() != reflect.Struct:
		return ObjectError
	case !rv.IsValid():
		return InitializedError
	}

	// Implements Unmarshaler interface.
	if rp.Type().Implements(reflect.TypeOf((*Unmarshaller)(nil)).Elem()) {
		// If there is the custom method, MarshlaENV - run it.
		if m := rp.MethodByName("UnmarshalENV"); m.IsValid() {
			tmp := m.Call([]reflect.Value{})
			err := tmp[0].Interface()
			if err != nil {
				return fmt.Errorf("unmarshal: %v", err)
			}
			return nil
		}
	}

	// Walk through all the fields of the transferred object.
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		key, sep := parseTag(field.Tag.Get("env"), field.Name, " ")
		key = fmt.Sprintf("%s%s", prefix, key)

		fmt.Println(":", key, "::", field.Name)

		// Change value.
		instance := rv.FieldByName(field.Name)
		kind := field.Type.Kind()
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			r, err := strToIntKind(Get(key), kind)
			if err != nil {
				return err
			}
			instance.SetInt(r)
		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64:
			r, err := strToUintKind(Get(key), kind)
			if err != nil {
				return err
			}
			instance.SetUint(r)
		case reflect.Float32, reflect.Float64:
			r, err := strToFloatKind(Get(key), kind)
			if err != nil {
				return err
			}
			instance.SetFloat(r)
		case reflect.Bool:
			r, err := strToBool(Get(key))
			if err != nil {
				return err
			}
			instance.SetBool(r)
		case reflect.String:
			instance.SetString(Get(key))
		case reflect.Array:
			max := instance.Type().Len()
			seq := strings.Split(Get(key), sep)
			if len(seq) > max {
				return fmt.Errorf("%d items overwhelms the [%d]array",
					len(seq), max)
			}

			err := setSequence(&instance, strings.Split(Get(key), sep))
			if err != nil {
				return err
			}
		case reflect.Slice:
			seq := strings.Split(Get(key), sep)
			tmp := reflect.MakeSlice(instance.Type(), len(seq), len(seq))
			err := setSequence(&tmp, strings.Split(Get(key), sep))
			if err != nil {
				return err
			}

			instance.Set(reflect.AppendSlice(instance, tmp))
		case reflect.Ptr:
			// Pointer support (structures only).
			item := instance.Type().Elem()
			if item.Kind() != reflect.Struct {
				return TypeError
			}

			// Determine the type of structure.
			elem := reflect.Indirect(reflect.New(item)).Interface()
			if _, ok := elem.(url.URL); ok {
				// Parse url.URL.
				u, err := url.Parse(Get(key))
				if err != nil {
					return err
				}
				instance.Set(reflect.ValueOf(u))
			} else {

				//p := reflect.New(reflect.TypeOf(instance))
				//p.Elem().Set(reflect.ValueOf(instance))

				//ins := p.Interface()
				//err := unmarshalENVOld(&ins, prefix)
				//fmt.Println("E:", err)
				//fmt.Printf("z: %T\n", ins)

				// tmp := instance.Elem()
				// tmp.Set(instance.Elem())
				// err := loadData(field.Type, tmp, instance, fmt.Sprintf("%s_", key))
				// fmt.Println("---L", tmp, err)
				//val := reflect.ValueOf(*instance)
				//vp := reflect.New(val.Type())
				//vp.Elem().Set(val)

				//loadData(field.Type, val, vp, fmt.Sprintf("%s_", key))
				//fmt.Println("x:", val, vp)
				//err := unmarshalENVOld(val, prefix)
				//fmt.Printf("%T, %v\n", elem, err)
				// err := loadData(field.Type, reflect.ValueOf(elem), instance, fmt.Sprintf("%s_", key))
				// fmt.Println("Error:", err)
				// tmp := instance.Elem()
				// //tmp.Set(nil)

				// err := unmarshalENVOld(tmp, prefix)
				// fmt.Printf("T~~>%T\n", tmp)
				// fmt.Printf("U~~>%T\n", item)
				// fmt.Printf("U++>%T\n", elem)
				// fmt.Printf("%v %v\n", elem, err)

				// // // Define: type, value and pointer.
				// // tmp := instance.Elem()
				// // //tmp.Set(instance.Ptr)
				// // err := loadData(field.Type, tmp, rp, fmt.Sprintf("%s_", key))
				// // fmt.Println("xxx:=>", tmp, err)
				// // // // // Another type of structure.
				// // // // //rp := instance
				// // // // ////rv := instance.Elem()
				// // // // rp = reflect.New(field.Type)
				// // // // tmp := rp.Elem()
				// // // // tmp.Set(instance)
				// // // // rv := reflect.ValueOf(field.Type)
				// // // // loadData(field.Type, rv, rp, fmt.Sprintf("%s_", key))
				// // // rp = reflect.New(field.Type)
				// // // tmp := rp.Elem()
				// // // tmp.Set(instance)
				// // // unmarshalENVOld(tmp, prefix)
				// // // fmt.Printf("xxx: %T", rp)
			}
		case reflect.Struct:
			if _, ok := instance.Interface().(url.URL); ok {
				// Parse url.URL.
				u, err := url.Parse(Get(key))
				if err != nil {
					return err
				}
				instance.Set(reflect.ValueOf(*u))
			} else {
				// Another type of structure.
				rp = reflect.New(field.Type)
				tmp := rp.Elem()
				tmp.Set(instance)
				loadData(field.Type, instance, rp, fmt.Sprintf("%s_", key))
			}
		default:
			return TypeError
		} // switch
	} // for

	return nil
}
*/

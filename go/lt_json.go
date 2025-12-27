package lt

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strconv"
)

type JSON map[string]any

func (j JSON) String() string {
	b, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

func Load[T any](j JSON, key ...string) (value T, ok bool) {
	var k string
	switch len(key) {
	case 0:
	case 1:
		k = key[0]
	default:
		return value, false
	}
	return value, load(j, k, reflect.ValueOf(&value))
}

func load(j JSON, key string, value reflect.Value) bool {
	value = reflect.Indirect(value)

	switch value.Kind() {
	// Bool
	case reflect.Bool:
		d, ok := j[key].(bool)
		value.SetBool(d)
		return ok

	// Int
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d, ok := j[key].(float64)
		value.SetInt(int64(d))
		return ok

	// Uint
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		d, ok := j[key].(float64)
		value.SetUint(uint64(d))
		return ok

	// Float
	case reflect.Float64, reflect.Float32:
		d, ok := j[key].(float64)
		value.SetFloat(d)
		return ok

	// String
	case reflect.String:
		switch d := j[key].(type) {
		case string:
			value.SetString(d)
			return true
		case float64:
			value.SetString(strconv.FormatFloat(d, 'g', -1, 64))
			return true
		default: // Unsupported
			return false
		}

	// [n]T (array)
	case reflect.Array:
		a, ok := j[key].([]any)
		if !ok || len(a) != value.Len() {
			return false
		}
		switch value.Type().Elem().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			for i, d := range a {
				d, ok := d.(float64)
				if !ok {
					return false
				}
				value.Index(i).SetInt(int64(d))
			}
			return true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			for i, d := range a {
				d, ok := d.(float64)
				if !ok {
					return false
				}
				value.Index(i).SetUint(uint64(d))
			}
			return true
		case reflect.Float64, reflect.Float32:
			for i, d := range a {
				d, ok := d.(float64)
				if !ok {
					return false
				}
				value.Index(i).SetFloat(float64(d))
			}
			return true
		default: // Unsupported
			return false
		}

	// []T (slice)
	case reflect.Slice:
		s, ok := j[key]
		if !ok {
			return false
		}
		switch s := s.(type) {
		case nil:
			value.Set(reflect.Zero(value.Type()))
			return true
		case string: // []byte (base64)
			if value.Type().Elem().Kind() != reflect.Uint8 {
				return false
			}
			b, err := base64.StdEncoding.DecodeString(s)
			value.SetBytes(b)
			return err == nil
		case []any:
			value.Set(reflect.MakeSlice(value.Type(), len(s), len(s)))
			switch value.Type().Elem().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				for i, d := range s {
					d, ok := d.(float64)
					if !ok {
						return false
					}
					value.Index(i).SetInt(int64(d))
				}
				return true
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				for i, d := range s {
					d, ok := d.(float64)
					if !ok {
						return false
					}
					value.Index(i).SetUint(uint64(d))
				}
				return true
			case reflect.Float64, reflect.Float32:
				for i, d := range s {
					d, ok := d.(float64)
					if !ok {
						return false
					}
					value.Index(i).SetFloat(d)
				}
				return true
			case reflect.String:
				for i, d := range s {
					d, ok := d.(float64)
					if !ok {
						return false
					}
					value.Index(i).SetString(strconv.FormatFloat(d, 'g', -1, 64))
				}
				return true
			default: // Unsupported
				return false
			}
		default: // Unsupported
			return false
		}

	// Map
	case reflect.Map:
		if key != "" {
			var ok bool
			if j, ok = j[key].(map[string]any); !ok {
				return false
			}
		}
		switch value.Type() {
		case reflect.TypeOf(JSON{}), reflect.TypeOf(map[string]any{}): // JSON
			value.Set(reflect.ValueOf(j))
			return true
		default: // map[int]struct
			if value.Type().Key().Kind() != reflect.Int && value.Type().Elem().Kind() != reflect.Struct {
				return false
			}
			value.Set(reflect.MakeMap(value.Type()))
			for k := range j {
				v := reflect.Indirect(reflect.New(value.Type().Elem()))
				if ok := load(j, k, v); !ok {
					return false
				}
				i, err := strconv.Atoi(k)
				if err != nil {
					return false
				}
				value.SetMapIndex(reflect.ValueOf(i), v)
			}
			return true
		}

	// Struct
	case reflect.Struct:
		if key != "" {
			var ok bool
			if j, ok = j[key].(map[string]any); !ok {
				return false
			}
		}
		for i := 0; i < value.NumField(); i++ {
			if k, ok := value.Type().Field(i).Tag.Lookup("json"); ok {
				if ok := load(j, k, value.Field(i)); !ok {
					return false
				}
			}
		}
		return true

	// Unsupported
	default:
		return false
	}
}

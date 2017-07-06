package tteok

import (
	"reflect"
	"runtime/debug"
)

func (l *log) addError(v interface{}) {
	err, ok := v.(error)
	if !ok {
		return
	}

	l.Error = err.Error()
	l.Stack = string(debug.Stack())
}

func (l *log) addMessage(v interface{}) {
	s, ok := v.(string)
	if !ok {
		return
	}

	if l.Messages != nil {
		l.Messages = append(l.Messages, s)
	} else if l.Message != "" {
		l.Messages = []string{l.Message, s}
		l.Message = ""
	} else {
		l.Message = s
	}
}

func (l *log) addMapToData(v interface{}) {
	val := reflect.ValueOf(v)
	for _, keyVal := range val.MapKeys() {
		data := val.MapIndex(keyVal)
		key, ok := keyVal.Interface().(string)
		if !ok {
			return
		}

		l.addPropToData(key, data.Interface())
	}
}

func (l *log) dereferencePointer(v interface{}) {
	ptr := reflect.ValueOf(v).Elem()
	i := ptr.Interface()
	l.enrich(i)
}

func (l *log) spreadSlice(v interface{}) {
	sl := reflect.ValueOf(v)

	for i := 0; i < sl.Len(); i++ {
		l.enrich(sl.Index(i).Interface())
	}
}

func (l *log) addStructToData(v interface{}) {
	val := reflect.ValueOf(v)
	fieldsCount := val.NumField()

	for i := 0; i < fieldsCount; i++ {
		field := val.Field(i)
		key := val.Type().FieldByIndex([]int{i}).Name

		switch field.Kind() {
		case reflect.Struct:
			l.addPropToData(key, field.Interface())
		case reflect.Ptr:
			l.addPropToData(key, field.Elem().Interface())
		default:
			l.addPropToData(key, field.Interface())
		}
	}
}

func (l *log) addPropToData(key string, value interface{}) {
	if l.Data == nil {
		l.Data = make(map[string]interface{})
	}

	l.Data[key] = value
}

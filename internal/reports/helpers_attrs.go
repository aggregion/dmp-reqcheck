package reports

import (
	"reflect"
	"strings"
)

func getReportAttrFromField(v reflect.Value) interface{} {
	kind := v.Type().Kind()

	switch kind {
	case reflect.String:
		return v.String()
	case reflect.Int:
	case reflect.Int32:
	case reflect.Int64:
		return v.Int()
	}

	return nil
}

func getReportInterfaceAttr(obj interface{}, attrName string) interface{} {
	v := reflect.ValueOf(obj).Elem()
	t := reflect.TypeOf(obj).Elem()
	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		attrsPrefix := t.Field(fieldIndex).Tag.Get("attrMap")
		if len(attrsPrefix) > 0 && strings.HasPrefix(attrName, attrsPrefix) {
			mp := v.Field(fieldIndex)
			key := reflect.ValueOf(strings.TrimPrefix(attrName, attrsPrefix))
			mv := mp.MapIndex(key)
			if !mv.IsValid() {
				return nil
			}

			return getReportAttrFromField(mp.MapIndex(key).Elem())
		}
		attrs := t.Field(fieldIndex).Tag.Get("attr")
		if attrs == attrName {
			return getReportAttrFromField(v.Field((fieldIndex)))
		}
	}

	return nil
}

func getReportStrAttr(obj interface{}, attrName string) string {
	attrVal := getReportInterfaceAttr(obj, attrName)
	if val, ok := attrVal.(string); ok {
		return val
	}

	return ""
}

func getReportIntAttr(obj interface{}, attrName string) int64 {
	attrVal := getReportInterfaceAttr(obj, attrName)
	if val, ok := attrVal.(int64); ok {
		return val
	}

	return 0
}

// GetReportAttrs .
func GetReportAttrs(obj interface{}) (result map[string]interface{}) {
	v := reflect.ValueOf(obj).Elem()
	t := reflect.TypeOf(obj).Elem()
	result = make(map[string]interface{})
	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		attrsPrefix := t.Field(fieldIndex).Tag.Get("attrMap")
		if len(attrsPrefix) > 0 {
			mp := v.Field(fieldIndex)
			for _, key := range mp.MapKeys() {
				result[attrsPrefix+key.String()] = getReportAttrFromField(mp.MapIndex(key).Elem())
			}
			continue
		}

		attr := t.Field(fieldIndex).Tag.Get("attr")
		if len(attr) > 0 {
			result[attr] = getReportAttrFromField(v.Field(fieldIndex))
		}
	}

	return
}

package tools

import (
	"ProjectX/library/storage/core"
	"encoding/json"
	"reflect"
	"strings"
)

const (
	VersionMark  = "Version"
	TagHashMark  = "hash"
	TagRangeMark = "range"
)

func GetSliceStructName(value interface{}) string {
	tp := reflect.TypeOf(value)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	if tp.Kind() != reflect.Slice {
		return ""
	}
	tp = tp.Elem()
	if tp.Kind() != reflect.Struct {
		return ""
	}
	return tp.Name()
}

func GetStructName(value interface{}) string {
	tp := reflect.TypeOf(value)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	if tp.Kind() != reflect.Struct {
		return ""
	}
	return tp.Name()
}

func GetStructOnlyName(value interface{}) string {
	if reflect.TypeOf(value).Kind() != reflect.Struct {
		return ""
	}
	return reflect.TypeOf(value).Name()
}

func GetHashAndRangeKey(value interface{}, useTag bool) (hashKey string, rangeKey string) {
	// 类型检查
	tp := reflect.TypeOf(value)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	k := tp.Kind()
	if k == reflect.Slice {
		tp = tp.Elem()
		k = tp.Kind()
	}
	if k != reflect.Struct {
		return
	}

	for i := 0; i < tp.NumField(); i++ {
		fieldType := tp.Field(i)
		tag := fieldType.Tag.Get("dynamo")
		tagArr := strings.Split(tag, ",")
		name := fieldType.Name
		if useTag {
			if len(tagArr) > 0 && tagArr[0] != "" {
				name = tagArr[0]
			}
		}
		for j := 1; j < len(tagArr); j++ {
			if tagArr[j] == TagHashMark {
				hashKey = name
				continue
			}
			if tagArr[j] == TagRangeMark {
				rangeKey = name
			}
		}
	}
	return
}

func TrySetStructVersion(value interface{}) (uint64, error) {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	} else {
		return 0, core.ErrUnsupportedValueType
	}
	if val.Kind() != reflect.Struct {
		return 0, core.ErrUnsupportedValueType
	}

	field := val.FieldByName(VersionMark)
	if field.Kind() != reflect.Invalid && field.CanInterface() {
		fieldInterface := field.Interface()
		if fieldInterface != nil {
			if version, ok := fieldInterface.(uint64); ok {
				field.SetUint(version + 1)
				return version, nil
			}
		}
	}
	return 0, core.ErrUnsupportedValueType
}

func GetFieldValueByName(value interface{}, name string) interface{} {
	tp := reflect.ValueOf(value)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	if tp.Kind() != reflect.Struct {
		return nil
	}

	field := tp.FieldByName(name)
	if field.Kind() != reflect.Invalid && field.CanInterface() {
		return field.Interface()
	}
	return nil
}

func GetFieldValueByRealName(value interface{}, name string) interface{} {
	tp := reflect.ValueOf(value)
	kp := reflect.TypeOf(value)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
		kp = kp.Elem()
	}
	if tp.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < kp.NumField(); i++ {
		fieldType := kp.Field(i)
		tag := fieldType.Tag.Get("dynamo")
		tagArr := strings.Split(tag, ",")
		realName := fieldType.Name
		if len(tagArr) > 0 && tagArr[0] != "" {
			realName = tagArr[0]
		}
		if realName == name {
			field := tp.Field(i)
			if field.Kind() != reflect.Invalid && field.CanInterface() {
				return field.Interface()
			}
		}
	}

	return nil
}

func DeepCopy(source interface{}, target interface{}) error {
	b, err := json.Marshal(source)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, target)
	return err
}

func GetSliceFromInterfacePtr(value interface{}) []interface{} {
	tp := reflect.TypeOf(value).Elem().Kind()
	if tp == reflect.Slice {
		s := reflect.ValueOf(value).Elem()
		temp := make([]interface{}, 0, 0)
		for i := 0; i < s.Len(); i++ {
			temp = append(temp, s.Index(i).Interface())
		}
		return temp
	}
	return nil
}

func GetInterfacePtr(value interface{}) interface{} {
	vp := reflect.ValueOf(value)
	if vp.Kind() == reflect.Ptr {
		return value
	} else {
		vNew := reflect.New(vp.Type())
		vNew.Elem().Set(vp)
		return vNew.Interface()
	}
}

// GenerateKeyMapFromStructSlice
// sl为StorageModel数组 []Storage
func GenerateKeyMapFromStructSlice(sl interface{}) map[interface{}]uint8 {
	var res = make(map[interface{}]uint8, 0)
	realVal := reflect.ValueOf(sl)
	valLen := realVal.Len()
	hashName := GetHashName(realVal.Interface())
	for i := 0; i < valLen; i++ {
		key := realVal.Index(i).FieldByName(hashName).Interface()
		res[key] = 0
	}
	return res
}

// GetHashName 获取主键在结构体中的名字
func GetHashName(model interface{}) string {
	vStruct := reflect.TypeOf(model).Elem()
	filedNum := vStruct.NumField()
	for i := 0; i < filedNum; i++ {
		filedTag := vStruct.Field(i).Tag.Get("dynamo")
		tagSlice := strings.Split(filedTag, ",")
		for _, v := range tagSlice {
			if v == "hash" {
				return vStruct.Field(i).Name
			}
		}
	}
	return ""
}

// GetStructVersionFromOriginData 获取版本
func GetStructVersionFromOriginData(value interface{}) (uint64, error) {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Struct {
		return 0, core.ErrUnsupportedValueType
	}

	field := val.FieldByName(VersionMark)
	if field.Kind() != reflect.Invalid && field.CanInterface() {
		fieldInterface := field.Interface()
		if fieldInterface != nil {
			if version, ok := fieldInterface.(uint64); ok {
				return version, nil
			}
		}
	}
	return 0, core.ErrUnsupportedValueType
}

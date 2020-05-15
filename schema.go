/*
 *  Copyright (c) 2020. aberic - All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gnomon

import (
	"go/ast"
	"reflect"
)

// Field 上一结构体所属参数信息
type Field struct {
	Name   string        // 参数名
	Value  reflect.Value // 参数值
	Tag    string        // 参数标签
	Kind   reflect.Kind  // 参数类型
	Fields []*Field      // 自身作为根结构体信息
}

// Schema 根结构体信息
type Schema struct {
	Model    interface{}
	Name     string // 结构体名称
	Fields   []*Field
	fieldMap map[string]*Field
}

// GetField GetField
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// GetValue GetValue
func (schema *Schema) GetValue(name string) reflect.Value {
	return schema.fieldMap[name].Value
}

// SchemaParse SchemaParse
func SchemaParse(key string, dest interface{}) *Schema {
	destValue := reflect.ValueOf(dest)
	destModelType := reflect.Indirect(destValue).Type()
	schema := &Schema{
		Model:    dest,
		Name:     destModelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < destModelType.NumField(); i++ {
		sf := destModelType.Field(i)
		if !sf.Anonymous && ast.IsExported(sf.Name) {
			destValue = ptrParse(destValue)
			field := fieldParse(key, sf, destValue.Field(i).Kind(), destValue.Field(i))
			schema.Fields = append(schema.Fields, field)
			schema.fieldMap[sf.Name] = field
		}
	}
	return schema
}

// fieldParse
func fieldParse(key string, sf reflect.StructField, kind reflect.Kind, value reflect.Value) *Field {
	field := &Field{
		Name:  sf.Name,
		Value: value,
		Kind:  kind,
	}
	value = ptrParse(value)
	field.Fields = fieldsParse(key, value)
	if v, ok := sf.Tag.Lookup(key); ok {
		field.Tag = v
	} else if v, ok := sf.Tag.Lookup(StringBuild(";", key)); ok {
		field.Tag = v
	}
	return field
}

func fieldsParse(key string, value reflect.Value) []*Field {
	switch value.Kind() {
	default:
		return nil
	case reflect.Struct:
		var fields []*Field
		fieldModelType := reflect.Indirect(value).Type()
		for i := 0; i < fieldModelType.NumField(); i++ {
			sf := fieldModelType.Field(i)
			if !sf.Anonymous && ast.IsExported(sf.Name) {
				childValue := value.Field(i)
				fields = append(fields, fieldParse(key, sf, childValue.Kind(), childValue))
			}
		}
		return fields
	case reflect.Slice:
		var fields []*Field
		for i := 0; i < value.Len(); i++ {
			oneOfArrValue := ptrParse(value.Index(i))
			field := &Field{
				Value:  oneOfArrValue,
				Kind:   oneOfArrValue.Kind(),
				Fields: fieldsParse(key, oneOfArrValue),
			}
			fields = append(fields, field)
		}
		return fields
		//case reflect.Map:
		//	fmt.Println("map.len", value.Len())
		//	var fields []*Field
		//	iter := value.MapRange()
		//	for iter.Next() {
		//		oneOfArrValue := ptrParse(iter.Value())
		//		field := &Field{
		//			Value:  oneOfArrValue,
		//			Kind:   oneOfArrValue.Kind(),
		//			Fields: fieldsParse(key, oneOfArrValue),
		//		}
		//		fields = append(fields, field)
		//		fmt.Println("iter key = ", iter.Key())
		//		fmt.Println("iter value = ", iter.Value())
		//	}
		//	return fields
	}
}

// ptrParse 解析指针类型为标准类型
func ptrParse(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		return ptrParse(value)
	}
	return value
}

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

import "testing"

type User struct {
	Name string `orm:"column:id" json:"name"`
	Age  int    `json:"age"`
}

type Family struct {
	Father    User              `json:"father"`
	Mother    *User             `json:"mother"`
	Children  []User            `json:"children"`
	Furniture map[string]string `json:"furniture"`
}

func TestParse1(t *testing.T) {
	schema1 := SchemaParse("orm", &User{Name: "tom", Age: 11})
	if schema1.Name != "User" || len(schema1.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	t.Log("schema1 Name", schema1.Name)
	t.Log("schema1 Model", schema1.Model)
	for _, f := range schema1.Fields {
		t.Log("schema1 field name", f.Name)
		t.Log("schema1 field tag", f.Tag)
		t.Log("schema1 field value", f.Value.Interface())
		t.Log("schema1 field kind", f.Kind)
	}

	schema2 := SchemaParse("json", &User{})
	t.Log("schema2 tag", schema2.GetField(schema1.Fields[0].Name).Tag)
	t.Log("schema2 tag", schema2.GetField(schema1.Fields[1].Name).Tag)
}

func TestParse2(t *testing.T) {
	schema := SchemaParse("json", &Family{
		Father: User{Name: "dad", Age: 33},
		Mother: &User{Name: "mom", Age: 30},
		Children: []User{
			{Name: "tom", Age: 11},
			{Name: "jack", Age: 13},
		},
		Furniture: map[string]string{
			"a": "A",
			"b": "B",
			"c": "C",
		},
	})
	rangeFields(schema.Fields, t)
}

func rangeFields(fields []*Field, t *testing.T) {
	for _, f := range fields {
		t.Log("field name", f.Name)
		t.Log("Field tag", f.Tag)
		t.Log("Field value", f.Value.Interface())
		t.Log("Field kind", f.Kind)
		rangeFields(f.Fields, t)
	}
}

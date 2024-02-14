package schema

import (
	"deferORM/dialect"
	"go/ast"
	"reflect"
)

// Field represents a column of database
type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}
	//遍歷結構體字段
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		//忽略struct中匿名和私有字段
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			//https://vimsky.com/examples/usage/reflect-tag-lookup-function-in-golang-with-examples.html
			//用于查找与标签字符串中的键相关联的值
			if v, ok := p.Tag.Lookup("deferorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

//reflect
//https://www.lanseyujie.com/post/golang-reflection.html

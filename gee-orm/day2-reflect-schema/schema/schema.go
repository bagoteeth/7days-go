package schema

import (
	"7days-go/gee-orm/day2-reflect-schema/dialect"
	"go/ast"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	//对象
	Model interface{}
	//表名
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func (r *Schema) GetField(name string) *Field {
	return r.fieldMap[name]
}

//dest为结构体对象指针
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	sc := &Schema{
		Model:      dest,
		Name:       modelType.Name(),
		Fields:     nil,
		FieldNames: nil,
		fieldMap:   make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		//解析结构体的每一个字段
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			//如果有，解析tag
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			sc.Fields = append(sc.Fields, field)
			sc.FieldNames = append(sc.FieldNames, p.Name)
			sc.fieldMap[p.Name] = field
		}
	}
	return sc
}

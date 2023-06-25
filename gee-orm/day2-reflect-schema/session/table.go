package session

import (
	"7days-go/gee-orm/day2-reflect-schema/log"
	"7days-go/gee-orm/day2-reflect-schema/schema"
	"fmt"
	"reflect"
	"strings"
)

//解析schema
func (r *Session) Model(value interface{}) *Session {
	if r.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(r.refTable.Model) {
		r.refTable = schema.Parse(value, r.dialect)
	}
	return r
}

//获取schema
func (r *Session) RefTable() *schema.Schema {
	if r.refTable == nil {
		log.Error("Model is not set")
	}
	return r.refTable
}

func (r *Session) CreateTable() error {
	table := r.refTable
	var columns []string
	for _, v := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", v.Name, v.Type, v.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := r.Raw(fmt.Sprintf("create table %s (%s);", table.Name, desc)).Exec()
	return err
}

func (r *Session) DropTable() error {
	_, err := r.Raw(fmt.Sprintf("drop table if exists %s;", r.refTable.Name)).Exec()
	return err
}

func (r *Session) HasTable() bool {
	sql, values := r.dialect.TableExistSQL(r.refTable.Name)
	row := r.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == r.RefTable().Name
}

package schema

import (
	"7days-go/gee-orm/day2-reflect-schema/dialect"
	"testing"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

var TestDial, _ = dialect.GetDialect("sqlite3")

func TestParse2(t *testing.T) {
	sc := Parse(&User{}, TestDial)
	if sc.Name != "User" || len(sc.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if sc.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}

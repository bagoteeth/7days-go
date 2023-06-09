package main

import (
	"github.com/golang/protobuf/proto"
	"log"
)

//测试proto数据序列化/反序列化的一致性
func main() {
	test := &Student{
		Name:   "bagoteeth",
		Male:   true,
		Scores: []int32{55, 66, 77},
	}

	data, err := proto.Marshal(test)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	newTest := &Student{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}

	if test.GetName() != newTest.GetName() {
		log.Fatal("data mismatch %q != %q", test.GetName(), newTest.GetName())
	}
}

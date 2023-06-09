package main

import (
	"log"
	"net/rpc"
)

type Result struct {
	Num, Ans int
}

func main() {
	c, _ := rpc.DialHTTP("tcp", "localhost:1234")
	var r Result
	//if err := c.Call("Cal.Square", 12, &r); err != nil{
	//	log.Fatal("Failed to call Cal.Square. ", err)
	//}

	asyncCall := c.Go("Cal.Square", 12, &r, nil)
	log.Printf("%d^2 = %d", r.Num, r.Ans)

	<-asyncCall.Done
	log.Printf("%d^2 = %d", r.Num, r.Ans)
}

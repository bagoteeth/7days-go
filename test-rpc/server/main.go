package main

import (
	"log"
	"net/http"
	"net/rpc"
)

type Result struct {
	Num, Ans int
}

type Cal int

func (r *Cal) Square(num int, res *Result) error {
	res.Num = num
	res.Ans = num * num
	return nil
}

func main() {
	rpc.Register(new(Cal))
	rpc.HandleHTTP()

	log.Printf("Serving RPC server on port %d\n", 1234)
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("Error serving: ", err)
	}
}

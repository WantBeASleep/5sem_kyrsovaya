package main

import (
	"fmt"
	gn "lib/generatelib"
	mt "lib/matrixes"
	rq "lib/requests"
)

const (
	m int = 2
	n int = 3
)

func getFormedTask() rq.ClientReq {
	expr := "a+b"
	a := gn.GenerateRandMatrix(m, n, 100)
	b := gn.GenerateRandMatrix(m, n, 100)
	
	h := map[string]mt.Matrix{
		"a": a,
		"b": b,
	}

	req := rq.ClientReq{
		Expr: expr,
		Matrixes: h,
	}
	return req
}

func SendRequest(port string) {
	Req := getFormedTask()
	var Ans mt.Matrix
	rq.SendRequest(port, "msolveproblem", Req, &Ans)
	fmt.Println(Ans)
}
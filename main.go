package main

import (
	"fmt"
	"go-spider/module/request"
)

func main() {
	res, _ := request.Get("www.baidu.com")
	fmt.Println(res.Response.Body)
}

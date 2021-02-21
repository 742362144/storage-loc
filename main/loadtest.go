package main

import (
	"fmt"
	"storage-loc/util"
	"time"
)

func main()  {
	start := time.Now().UnixNano()
	for i := 0; i < 10000; i++ {
		util.Load("/root/go/src/storage-loc/main/add.so", "Add")
	}
	end := time.Now().UnixNano()
	fmt.Println((end-start)/ 10000)
}

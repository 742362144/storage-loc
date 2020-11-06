package main

import (
	"fmt"
	"github.com/742362144/storage-loc/util"
	"github.com/dterei/gotsc"
)

const N = 100

func main() {
	tsc := gotsc.TSCOverhead()
	fmt.Println("TSC Overhead:", tsc)

	start := gotsc.BenchStart()
	value := util.GENERATEVALUE(128)
	util.MD5V(value)
	end := gotsc.BenchEnd()
	avg := (end - start - tsc) / N

	fmt.Println("Cycles:", avg)
	for i:=1; i<9; i++ {
		fmt.Println(i * 25826)
	}
}
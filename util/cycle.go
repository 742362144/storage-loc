package util

import (
	"fmt"
	"github.com/dterei/gotsc"
	"time"
)

const N = 10000

func main() {
	value := GENERATEVALUE(8)

	tsc := gotsc.TSCOverhead()
	fmt.Println("TSC Overhead:", tsc)

	start := gotsc.BenchStart()
	cycle_start := time.Now().UnixNano()
	for i:=1; i<N; i++ {
		MD5V(value)
	}
	//util.CYCLE1000()
	cycle_end := time.Now().UnixNano()
	end := gotsc.BenchEnd()

	avg := (end - start - tsc) / N
	fmt.Println(cycle_start)
	fmt.Println(cycle_end)
	cycle_avg := (cycle_end - cycle_start) / N

	fmt.Println("avg time:", cycle_avg)
	fmt.Println("avg cycles:", avg)
	for i:=1; i<9; i++ {
		fmt.Println(i * 922)
	}
}
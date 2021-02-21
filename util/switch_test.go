package util

import (
	"sync"
	"testing"
)

func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	channel := make(chan struct{})

	sender := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			channel <- struct{}{}
		}
	}

	receiver := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			<-channel
		}
	}

	wg.Add(2)
	go sender()
	go receiver()

	b.StartTimer()
	close(begin)
	wg.Wait()
}
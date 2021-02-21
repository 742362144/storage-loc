package sche

import (
	"fmt"
	"time"
	"context"
)

func F(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("TIME OUT")
			cancel()
			return ctx.Err()
		default:
			time.Sleep(1 * time.Second)
			fmt.Println("No: ", i)
		}
	}
	fmt.Println("ALL DONE")
	return nil
}

/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"github.com/742362144/storage-loc/pb"
	"log"
	"math/big"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKVServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//for i:=0; i<100000; i++ {
	//	r, err := c.Op(ctx, &pb.Request{OpType: 1, Key: "key", Value: ""})
	//	if err != nil {
	//		log.Fatalf("could not greet: %v", err)
	//	}
	//	if i % 1000 == 0 {
	//		log.Printf("Greeting: %s", r.GetValue())
	//	}
	//}

	base := time.Now().UnixNano() / 1e9

	//添加kv
	val := generateValue(8)
	for i := 0; i < 100000; i++ {
		r, err := c.Op(ctx, &pb.Request{OpType: 2, Key: strconv.Itoa(i), Value: val})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		if i % 1000 == 0 {
			println(time.Now().UnixNano() / 1e9 - base)
			log.Printf("Greeting: %s", r.GetValue())
		}
	}
}

func generateValue(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0;i < len ;i++  {
		randomInt,_ := rand.Int(rand.Reader,bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}
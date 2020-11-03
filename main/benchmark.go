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
	"crypto/md5"
	"encoding/hex"
	"flag"
	"github.com/742362144/storage-loc/pb"
	"log"
	"math/big"
	"strconv"
	"sync"
	"time"
	"crypto/rand"
	mrand "math/rand"
	"google.golang.org/grpc"
)

const (
	GET = 1
	PUT = 2
	RUN = 3
)

func main() {
	HOST := *flag.String("ip", "133.133.135.22", "host ip")
	PORT := *flag.String("port", "50051", "port")
	MODEL := *flag.Int("model", 2, "operation model")
	PARALLEL := *flag.Int("parallel", 32, "PARALLEL nums")
	SIZE := *flag.Int("size", 16, "request value SIZE")
	NUM := *flag.Int("num", 10000, "request value SIZE")

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()


	var wg sync.WaitGroup
	timeChan := make(chan int, 100)
	wg.Add(PARALLEL)
	for j:=0; j< PARALLEL; j++ {
		if MODEL == PUT {
			go func(index int, timeChan chan int) {
				put(HOST, PORT, SIZE, NUM, index, timeChan)
				wg.Done()
			}(j, timeChan)
		} else if MODEL == GET {
			go func(timeChan chan int) {
				get(HOST, PORT, NUM, PARALLEL, timeChan)
				wg.Done()
			}(timeChan)
		}
	}
	wg.Wait()
	total := 0
	for i := 0; i < PARALLEL; i++ {
		total += <-timeChan
	}
	throught := float64(NUM * PARALLEL * PARALLEL * 1000)/ float64(total)
	log.Printf("throught: %f", throught)
}

func put(host, port string, size, num, index int, timeChan chan int) {
	base := time.Now().UnixNano() / 1e6
	ctx := context.Background()
	//添加kv
	pre := "client" + strconv.Itoa(index) + "_"

	for i := 0; i < num; i++ {
		client, conn := getConn(host, port)
		val := generateValue(size)
		_, err := client.Op(ctx, &pb.Request{OpType: PUT, Key: pre+strconv.Itoa(i), Value: val})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		conn.Close()
	}
	cost := time.Now().UnixNano() / 1e6 - base
	//log.Printf("%s finish %d\n", pre, time.Now().UnixNano() / 1e6 - base )
	timeChan <- int(cost)
}

func get(host, port string, num, parallel int, timeChan chan int) {

	base := time.Now().UnixNano() / 1e6

	//添加kv
	pre := "client" + strconv.Itoa(mrand.Intn(parallel)) + "_"
	mrand.Seed(time.Now().UnixNano())
	ctx := context.Background()

	for i := 0; i < num; i++ {
		for j:=0; j<8; j++ {
			client, conn := getConn(host, port)
			res, err := client.Op(ctx, &pb.Request{OpType: GET, Key: pre+strconv.Itoa(mrand.Intn(num)), Value: ""})
			for i:=0; i<100; i++ {
				md5V(res.GetValue())
			}
			conn.Close()
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
		}
	}
	cost := time.Now().UnixNano() / 1e6 - base
	log.Printf("%s finish %d\n", pre, time.Now().UnixNano() / 1e6 - base )
	timeChan <- int(cost)

}

func md5V(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
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

func getConn(host, port string) (pb.KVServiceClient, *grpc.ClientConn) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(host + ":" + port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewKVServiceClient(conn)
	return client, conn
}
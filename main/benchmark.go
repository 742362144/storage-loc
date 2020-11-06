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
	"context"
	"flag"
	"fmt"
	"github.com/742362144/storage-loc/pb"
	"github.com/742362144/storage-loc/util"
	"google.golang.org/grpc"
	"log"
	mrand "math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	GET = 1
	PUT = 2
	RUN = 3
)

func main() {
	HOST := *flag.String("host", "133.133.133.127", "host ip")
	PORT := *flag.String("port", "50051", "port")
	MODEL := *flag.Int("model", 1, "operation model")
	PARALLEL := *flag.Int("parallel", 16, "PARALLEL nums")
	SIZE := *flag.Int("size", 128, "request value SIZE")
	NUM := *flag.Int("num", 10000, "request value SIZE")
	DEEPTH := *flag.Int("deepth", 6, "request value SIZE")

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
				get(HOST, PORT, DEEPTH,NUM, PARALLEL, timeChan)
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
	fmt.Printf("throught: %f\n", throught)

	//_, filedir,_,_ := runtime.Caller(1)
	//fmt.Println(filedir)
	//f, err := os.OpenFile("C:/Users/win/go/src/github.com/742362144/storage-loc/main/result", os.O_RDWR,0600)
	//defer f.Close()
	//if err !=nil {
	//	fmt.Println(err.Error())
	//} else {
	//	jsonStr,_ :=ioutil.ReadAll(f)
	//	var mapResult map[string]string
	//	json.Unmarshal(jsonStr, &mapResult)
	//	if mapResult == nil {
	//		mapResult = make(map[string]string)
	//	}
	//	mapResult["client" + strconv.Itoa(DEEPTH)] = strconv.Itoa(int(throught))
	//	jsonStr, err = json.Marshal(mapResult)
	//
	//	_,err=f.Write(jsonStr)
	//	fmt.Println(string(jsonStr))
	//}

}

func put(host, port string,size, num, index int, timeChan chan int) {
	base := time.Now().UnixNano() / 1e6
	ctx := context.Background()
	conn, err := grpc.Dial(host + ":" + port, grpc.WithInsecure(), grpc.WithBlock())
	client := pb.NewKVServiceClient(conn)
	defer conn.Close()

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	//添加kv
	pre := "client" + strconv.Itoa(index) + "_"

	for i := 0; i < num; i++ {
		val := util.GENERATEVALUE(size)
		_, err := client.Op(ctx, &pb.Request{OpType: PUT, Key: pre+strconv.Itoa(i), Value: val})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
	}
	cost := time.Now().UnixNano() / 1e6 - base
	//log.Printf("%s finish %d\n", pre, time.Now().UnixNano() / 1e6 - base )
	timeChan <- int(cost)
}

func get(host, port string, deepth, num, parallel int, timeChan chan int) {
	base := time.Now().UnixNano() / 1e6

	//添加kv
	pre := "client" + strconv.Itoa(mrand.Intn(parallel)) + "_"
	mrand.Seed(time.Now().UnixNano())
	ctx := context.Background()
	client, conn := getConn(host, port)
	defer conn.Close()
	for i := 0; i < num; i++ {
		//client.Op(ctx, &pb.Request{OpType: GET, Key: pre+strconv.Itoa(mrand.Intn(num)), Value: ""})
		for j:=0; j<deepth; j++ {
			//res, err := client.Op(ctx, &pb.Request{OpType: GET, Key: pre+strconv.Itoa(mrand.Intn(num)), Value: ""})
			//for i:=0; i<10; i++ {
			//	util.MD5V(res.GetValue())
			//}
			res, err := client.Op(ctx, &pb.Request{OpType: GET, Key: pre+strconv.Itoa(mrand.Intn(num)), Value: ""})
			util.MD5V(res.GetValue())
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			//log.Printf("value %s pi: %d", res.GetValue(), pigo.Pi(1000))
		}
	}
	cost := time.Now().UnixNano() / 1e6 - base
	log.Printf("%s finish %d\n", pre, time.Now().UnixNano() / 1e6 - base )
	timeChan <- int(cost)


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
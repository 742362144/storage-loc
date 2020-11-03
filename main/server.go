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

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/742362144/storage-loc/pb"
	"github.com/boltdb/bolt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)


// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedKVServiceServer
	lock *sync.RWMutex
	db map[string]string
	boltDB *bolt.DB
}

// SayHello implements helloworld.GreeterServer
//func (s *server) Op(ctx context.Context, in *pb.Request) (*pb.Response, error) {
//	log.Printf("Received: %v", in.GetKey())
//	if in.OpType == 1 {
//		s.lock.RLock()
//		defer s.lock.RUnlock()
//		if _, ok := s.db[in.GetKey()]; ok {
//			return &pb.Response{OpType: in.GetOpType(), Value: s.db[in.GetKey()]}, nil
//		} else {
//			return &pb.Response{OpType: in.GetOpType(), Value: "key not found"}, nil
//		}
//	} else if in.OpType == 2 {
//		s.lock.Lock()
//		defer s.lock.Unlock()
//		s.db[in.GetKey()] = in.GetValue()
//		return &pb.Response{OpType: in.GetOpType(), Value: in.GetValue()}, nil
//	}
//	return &pb.Response{OpType: 0, Value: in.GetValue()}, nil
//}

func (s *server) Op(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Printf("Received: %v", in.GetKey())



	var res pb.Response
	if in.OpType == 1 {
		s.boltDB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MyBucket"))
			v := b.Get([]byte(in.GetKey()))
			if v != nil {
				res.OpType =  in.GetOpType()
				res.Value = string(v)
			} else {
				res.OpType =  in.GetOpType()
				res.Value = string(v)
			}
			return nil
		})
	} else if in.OpType == 2 {
		s.boltDB.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
			if err != nil {
				return fmt.Errorf("create bucket: %v", err)
			}

			if err = b.Put([]byte(in.GetKey()), []byte(in.GetValue())); err != nil {
				return err
			}
			return nil
		})

		res.OpType =  in.GetOpType()
		res.Value = string(in.GetValue())
	}
	return &res, nil
}


func main() {
	HOST := *flag.String("ip", "133.133.135.22", "host ip")
	PORT := *flag.String("port", "50051", "port")

	lis, err := net.Listen("tcp", HOST + ":" + PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	kvserver := new(server)
	kvserver.lock = new(sync.RWMutex)
	kvserver.db = make(map[string]string)

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	kvserver.boltDB = db
	defer db.Close()

	pb.RegisterKVServiceServer(s, kvserver)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

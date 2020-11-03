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
	"github.com/742362144/storage-loc/pb"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedKVServiceServer
	lock *sync.RWMutex
	db map[string]string
}

// SayHello implements helloworld.GreeterServer
func (s *server) Op(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Printf("Received: %v", in.GetKey())
	if in.OpType == 1 {
		s.lock.RLock()
		defer s.lock.RUnlock()
		if _, ok := s.db[in.GetKey()]; ok {
			return &pb.Response{OpType: in.GetOpType(), Value: s.db[in.GetKey()]}, nil
		} else {
			return &pb.Response{OpType: in.GetOpType(), Value: "key not found"}, nil
		}
	} else if in.OpType == 2 {
		s.lock.Lock()
		defer s.lock.Unlock()
		s.db[in.GetKey()] = in.GetValue()
		return &pb.Response{OpType: in.GetOpType(), Value: in.GetValue()}, nil
	}
	return &pb.Response{OpType: 0, Value: in.GetValue()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	kvserver := new(server)
	kvserver.lock = new(sync.RWMutex)
	kvserver.db = make(map[string]string)
	pb.RegisterKVServiceServer(s, kvserver)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

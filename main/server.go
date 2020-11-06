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
	"github.com/742362144/storage-loc/pb"
	"github.com/742362144/storage-loc/util"
	"github.com/boltdb/bolt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)


// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedKVServiceServer
	lock *sync.RWMutex
	deepth int
	db map[string]string
	boltDB *bolt.DB
	ldb *leveldb.DB
}

//func (s *server) Op(ctx context.Context, in *pb.Request) (*pb.Response, error) {
//	log.Printf("Received: %v", in.GetKey())
//	var res pb.Response
//	if in.OpType == 1 {
//		s.lock.RLock()
//		defer s.lock.RUnlock()
//		if _, ok := s.db[in.GetKey()]; ok {
//			 res.OpType = in.GetOpType()
//			 res.Value = s.db[in.GetKey()]
//		} else {
//			 res.OpType = in.GetOpType()
//			 res.Value = "key not found"
//		}
//		for j:=0; j<6; j++ {
//			for i:=0; i<100; i++ {
//				util.MD5V(res.GetValue())
//			}
//			//log.Printf("value %s pi: %d", res.GetValue(), pigo.Pi(1000))
//		}
//	} else if in.OpType == 2 {
//		s.lock.Lock()
//		defer s.lock.Unlock()
//		s.db[in.GetKey()] = in.GetValue()
//		res.OpType = in.GetOpType()
//		res.Value = in.GetValue()
//	}
//	return &res, nil
//}

//func (s *server) Op(ctx context.Context, in *pb.Request) (*pb.Response, error) {
//	log.Printf("Received: %v", in.GetKey())
//
//	var res pb.Response
//	if in.OpType == 1 {
//		s.boltDB.View(func(tx *bolt.Tx) error {
//			b := tx.Bucket([]byte("MyBucket"))
//			v := b.Get([]byte(in.GetKey()))
//			if v != nil {
//				res.OpType =  in.GetOpType()
//				res.Value = string(v)
//			} else {
//				res.OpType =  in.GetOpType()
//				res.Value = string(v)
//			}
//			return nil
//		})
//	} else if in.OpType == 2 {
//		s.boltDB.Update(func(tx *bolt.Tx) error {
//			b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
//			if err != nil {
//				return fmt.Errorf("create bucket: %v", err)
//			}
//
//			if err = b.Put([]byte(in.GetKey()), []byte(in.GetValue())); err != nil {
//				return err
//			}
//			return nil
//		})
//
//		res.OpType =  in.GetOpType()
//		res.Value = string(in.GetValue())
//	}
//	return &res, nil
//}

func (s *server) Op(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Printf("Received: %v", in)

	var res pb.Response
	if in.OpType == 1 {
		data, err := s.ldb.Get([]byte(in.GetKey()), nil)
		if err == nil {
			res.OpType =  in.GetOpType()
			res.Value = string(data)
			//log.Printf("key: %s", data)
			for j:=0; j<s.deepth; j++ {
				util.MD5V(res.GetValue())
				//log.Printf("value %s pi: %d", res.GetValue(), pigo.Pi(1000))
			}
		}
	} else if in.OpType == 2 {
		err := s.ldb.Put([]byte(in.GetKey()), []byte(in.GetValue()), nil)
		if err == nil {
			res.OpType =  in.GetOpType()
			res.Value = string(in.GetValue())
		}
	}
	return &res, nil
}

func main() {
	HOST := *flag.String("host", "133.133.133.127", "host ip")
	PORT := *flag.String("port", "50051", "port")
	DEEPTH := *flag.Int("deepth", 0, "request value SIZE")

	lis, err := net.Listen("tcp", HOST + ":" + PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	kvserver := new(server)
	kvserver.lock = new(sync.RWMutex)
	kvserver.deepth = DEEPTH
	kvserver.db = make(map[string]string)

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	kvserver.boltDB = db
	defer db.Close()

	ldb, err := leveldb.OpenFile("lerveldb", nil)
	kvserver.ldb = ldb
	defer ldb.Close()

	pb.RegisterKVServiceServer(s, kvserver)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

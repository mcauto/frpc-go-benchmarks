/*
Copyright 2022 Loophole Labs

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"io"
	"net"
	"os"
	"time"

	"github.com/loopholelabs/frpc-go-benchmarks/config"
	"github.com/loopholelabs/frpc-go-benchmarks/grpc/benchmark"
	grpcZeroLog "github.com/philip-bui/grpc-zerolog"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type svc struct {
	benchmark.UnimplementedBenchmarkServiceServer
}

func (s *svc) Benchmark(_ context.Context, req *benchmark.Request) (*benchmark.Response, error) {
	if config.Sleep {
		time.Sleep(time.Microsecond * 50)
	}
	res := new(benchmark.Response)
	res.Message = req.Message
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		panic(err)
	}

	shouldLog := len(os.Args) > 2
	log := zerolog.New(io.Discard)
	grpcServer := grpc.NewServer(
		grpcZeroLog.UnaryInterceptorWithLogger(&log),
	)

	benchmark.RegisterBenchmarkServiceServer(grpcServer, new(svc))

	if shouldLog {
		go func() {
			err = grpcServer.Serve(lis)
			if err != nil {
				panic(err)
			}
		}()

		for {
			// log.Printf("Num goroutines: %d\n", runtime.NumGoroutine())
			time.Sleep(time.Millisecond * 500)
		}
	} else {
		err = grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}
}

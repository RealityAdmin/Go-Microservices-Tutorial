package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"movieexamplekhubaib.com/gen"
	"movieexamplekhubaib.com/pkg/discovery"
	"movieexamplekhubaib.com/pkg/discovery/consul"
	"movieexamplekhubaib.com/rating/internal/controller/rating"
	grpchandler "movieexamplekhubaib.com/rating/internal/handler/grpc"
	"movieexamplekhubaib.com/rating/internal/repository/memory"
)

const serviceName = "rating"

func main() {
	// log.Println("Starting rating service")
	// repo := memory.New()
	// ctrl := rating.New(repo)
	// h := httphandler.New(ctrl)
	// http.Handle("/rating", http.HandlerFunc(h.Handle))
	// if err := http.ListenAndServe(":8082", nil); err != nil {
	// 	panic(err)
	// }

	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()
	log.Printf("Starting metadata service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	// On separate thread, keep reporting the healthy state
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Once this program halts, deregister
	defer registry.Deregister(ctx, instanceID, serviceName)

	// repo := memory.New()
	// svc := rating.New(repo)
	// h := httphandler.New(svc)
	// http.Handle("/rating", http.HandlerFunc(h.Handle))
	// if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
	// 	panic(err)
	// }

	repo := memory.New()
	svc := rating.New(repo)
	h := grpchandler.New(svc)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	srv.Serve(lis)

}

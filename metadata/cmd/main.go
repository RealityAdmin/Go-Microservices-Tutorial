package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"gopkg.in/yaml.v3"
	"movieexamplekhubaib.com/gen"
	"movieexamplekhubaib.com/metadata/internal/controller/metadata"
	"movieexamplekhubaib.com/metadata/internal/repository/memory"
	"movieexamplekhubaib.com/pkg/discovery"
	"movieexamplekhubaib.com/pkg/discovery/consul"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	grpchandler "movieexamplekhubaib.com/metadata/internal/handler/grpc"
)

const serviceName = "metadata"

func main() {

	// var port int
	// flag.IntVar(&port, "port", 8081, "API handler port")
	// flag.Parse()
	// log.Printf("Starting metadata service on port %d", port)

	f, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var cfg serviceConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	port := cfg.APIConfig.Port

	log.Printf("Starting movie metadata service on port %v", port)

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%v", cfg.APIConfig.Port)); err != nil {
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

	repo := memory.New()
	svc := metadata.New(repo)
	// h := httphandler.New(svc)
	// http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	// if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
	// 	panic(err)
	// }

	h := grpchandler.New(svc)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", cfg.APIConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	reflection.Register(server)
	gen.RegisterMetadataServiceServer(server, h)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"movieexamplekhubaib.com/gen"
	"movieexamplekhubaib.com/movie/internal/controller/movie"
	metadatagateway "movieexamplekhubaib.com/movie/internal/gateway/metadata/http"
	ratinggateway "movieexamplekhubaib.com/movie/internal/gateway/rating/http"
	"movieexamplekhubaib.com/pkg/discovery"
	"movieexamplekhubaib.com/pkg/discovery/consul"

	grpchandler "movieexamplekhubaib.com/movie/internal/handler/grpc"
)

// sudo docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0

const serviceName = "movie"

func main() {
	// log.Println("Starting movie service")
	// metadataGateway := metadatagateway.New("localhost:8081")
	// ratingGateway := ratinggateway.New("localhost:8082")
	// ctrl := movie.New(ratingGateway, metadataGateway)
	// h := httphandler.New(ctrl)
	// http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	// if err := http.ListenAndServe(":8083", nil); err != nil {
	// 	panic(err)
	// }

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

	log.Printf("Starting movie rating service on port %v", port)

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%v", port)); err != nil {
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

	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	svc := movie.New(ratingGateway, metadataGateway)
	// h := httphandler.New(svc)
	// http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	// if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
	// 	panic(err)
	// }
	h := grpchandler.New(svc)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMovieServiceServer(srv, h)
	srv.Serve(lis)
}

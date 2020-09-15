package main

import (
	"context"
	"log"
	"os"
	"net"

	"google.golang.org/grpc"
	pb "github.com/kikeyama/grpc-sfx-demo/pb"

	grpctrace "github.com/signalfx/signalfx-go-tracing/contrib/google.golang.org/grpc"
	"github.com/signalfx/signalfx-go-tracing/tracing"
)

//var logger log.Logger
var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

const (
	port = ":50051"
	serviceName = "kikeyama_grpc_server"
)

// getMessageService implements pb.DemoServer.GetMessage
func getMessageService(ctx context.Context, in *pb.DemoRequest) (*pb.DemoReply, error) {
	logger.Printf("level=info message=\"Received: %v\"", in.GetName())
	return &pb.DemoReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// Create a listener for the server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatalf("level=fatal message=\"failed to listen: %v\"", err)
	}

	// Use signalfx tracing
	tracing.Start(tracing.WithGlobalTag("stage", "demo"), tracing.WithServiceName(serviceName))
//	tracing.Start()
	defer tracing.Stop()

	// Create the server interceptor using the grpc trace package.
	si := grpctrace.StreamServerInterceptor(grpctrace.WithServiceName(serviceName))
	ui := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName(serviceName))

	// Initialize the grpc server as normal, using the tracing interceptor.
	//s := grpc.NewServer()
	s := grpc.NewServer(grpc.StreamInterceptor(si), grpc.UnaryInterceptor(ui))

	pb.RegisterDemoService(s, &pb.DemoService{GetMessageService: getMessageService})
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("level=fatal message=\"failed to serve: %v\"", err)
	}
}

package main

import (
	"fmt"
	"net"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	pb "github.com/rexlx/portping/src/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	svc := Server{}
	// s := grpc.NewServer()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	grpcEndpoint := fmt.Sprintf("0.0.0.0:%s", port)
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	grpcServer := grpc.NewServer(
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry),
		),
	)
	pb.RegisterReacherServer(grpcServer, &svc)
	listen, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.WithField("grpcEndpoint", grpcEndpoint).Info("Starting: gRPC Listener")
	logrus.Fatal(grpcServer.Serve(listen))
}

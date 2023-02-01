package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"time"

	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/rexlx/portping/src/proto"
	grpcMetadata "google.golang.org/grpc/metadata"
)

// you can call the example functions inside main here

var addr string = "mrbyte:8080"

// var certPath string = "/Users/rxlx/bin/data"
var logPath string = "unary_client.log"

// var statsFile string = "o_ercot_rts_conditions.csv"

func main() {
	conn, err := NewConn(addr, true)
	// conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	c := pb.NewReacherClient(conn)
	pingRequest(c)
}

func pingRequest(c pb.ReacherClient) {
	req := &pb.PingRequest{
		Address: "8.8.8.8",
		Port:    53,
		Wait:    555,
		Count:   5,
	}
	stream, err := c.Ping(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%v\t%v", msg.Duration, msg.Msg)
	}
}

func pingRequestWithAuth(c pb.ReacherClient, audience string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create an identity token.
	// With a global TokenSource tokens would be reused and auto-refreshed at need.
	// A given TokenSource is specific to the audience.
	tokenSource, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		log.Printf("idtoken.NewTokenSource: %v", err)
	}
	token, err := tokenSource.Token()
	if err != nil {
		log.Printf("TokenSource.Token: %v", err)
	}

	// Add token to gRPC Request.
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)
	req := &pb.PingRequest{
		Address: "google.com",
		Port:    443,
		Wait:    555,
		Count:   6,
	}
	stream, err := c.Ping(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%v\t%v", msg.Duration, msg.Msg)
	}

}

// NewConn creates a new gRPC connection.
// host should be of the form domain:port, e.g., example.com:443
func NewConn(host string, insecure bool) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if host != "" {
		opts = append(opts, grpc.WithAuthority(host))
	}

	if insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	return grpc.Dial(host, opts...)
}

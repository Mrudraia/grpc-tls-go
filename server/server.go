package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	pb "github.com/mrudraia/grpc-tls-go/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"k8s.io/klog"
)

type ROSAServer struct {
	pb.UnimplementedRosaServiceServer
}

func (*ROSAServer) InstallAgent(ctx context.Context, req *pb.InstallRequest) (*pb.InstallResponse, error) {
	fmt.Println("Install Agent")
	name := req.GetName()
	kind := req.GetKind()
	// data := req.GetData()

	res := fmt.Sprintf("the name is %v and the data is  %v", name, kind)
	fmt.Println(res)
	return &pb.InstallResponse{Response: res}, nil

}

func (s *ROSAServer) Start(port int, opts []grpc.ServerOption) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		klog.Fatalf("failed to listen: %v", err)
	}

	certFile := "cert/server-cert.pem"
	keyFile := "cert/server-key.pem"

	creds, sslErr := credentials.NewClientTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatal(err)
	}

	opts = append(opts, grpc.Creds(creds))
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterRosaServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	err = grpcServer.Serve(lis)
	if err != nil {
		klog.Fatalf("Failed to start grpc server : %v", err)
	}

}

func main() {
	// listen port
	lis, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		log.Fatalf("list port err: %v", err)
	}

	// read ca's cert, verify to client's certificate
	caPem, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		log.Fatal(err)
	}

	// create cert pool and appen ca's cert
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		log.Fatal(err)
	}

	// read server cert and key
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	// configuration of the certificate what we want to
	conf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	tlsCredentials := credentials.NewTLS(conf)

	grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))

	pb.RegisterRosaServiceServer(grpcServer, &ROSAServer{})

	log.Printf("listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc serve err: %v", err)
	}

}

package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	pb "github.com/mrudraia/grpc-tls-go/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ROSAClient struct {
	Client     pb.RosaServiceClient
	clientConn *grpc.ClientConn
	timeout    time.Duration
}

func NewRosaClient(ctx context.Context, serverAddr string, timeout time.Duration) (*ROSAClient, error) {
	apiCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// read ca's cert
	caCert, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		log.Fatal(caCert)
	}

	// create cert pool and append ca's cert
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatal(err)
	}

	//read client cert
	clientCert, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithBlock(),
	}

	conn, err := grpc.DialContext(apiCtx, "0.0.0.0:9000", opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	return &ROSAClient{
		Client:     pb.NewRosaServiceClient(conn),
		clientConn: conn,
		timeout:    timeout,
	}, nil
}

func (cc *ROSAClient) Close() {
	if cc.clientConn != nil {
		_ = cc.clientConn.Close()
		cc.clientConn = nil
	}
	cc.Client = nil
}

func (cc *ROSAClient) AgentInstallation(ctx context.Context, name, kind string) (*pb.InstallResponse, error) {
	if cc.Client == nil || cc.clientConn == nil {
		return nil, fmt.Errorf("rosa client is closed")
	}

	req := &pb.InstallRequest{
		Name: "abcd",
		Kind: "efgh",
		Data: []byte{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := cc.Client.InstallAgent(ctx, req) // Cleint.InstallAgent(ctx2, &pb.InstallRequest{Name: "djnd", Kind: "sdsd", Data: []byte{1, 2, 3}})
	if err != nil {
		log.Fatal(err)
	}

	return resp, nil

}

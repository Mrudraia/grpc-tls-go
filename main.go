package main

import (
	"context"
	"time"

	client "github.com/mrudraia/grpc-tls-go/client"
	"k8s.io/klog"
)

func main() {
	ctx := context.Background()
	rosaClient, err := client.NewRosaClient(
		ctx,
		"localhost:3000",
		10*time.Second,
	)

	if err != nil {
		klog.Exitf("Failed to create grpc client: %v", err)
	}

	defer rosaClient.Close()

	rosaClient.AgentInstallation(ctx, "test", "test123")

}

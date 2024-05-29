package keto

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

// import (
// 	"fmt"

// 	acl "github.com/ory/keto/proto/ory/keto/acl/v1alpha1"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// type KetoClient struct {
// 	ReadClient   acl.ReadServiceClient
// 	WriteClient  acl.WriteServiceClient
// 	CheckClient  acl.CheckServiceClient
// 	ExpandClient acl.ExpandServiceClient
// }

// func NewKetoClient() (*KetoClient, error) {
// 	rConn, err := grpc.Dial("keto:4466", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to keto read grpc read server: %w", err)
// 	}

// 	wConn, err := grpc.Dial("keto:4467", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to keto read grpc write server: %w", err)
// 	}

// 	c := &KetoClient{
// 		ReadClient:   acl.NewReadServiceClient(rConn),
// 		WriteClient:  acl.NewWriteServiceClient(wConn),
// 		CheckClient:  acl.NewCheckServiceClient(rConn),
// 		ExpandClient: acl.NewExpandServiceClient(rConn),
// 	}

// 	return c, nil
// }

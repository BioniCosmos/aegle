package api

import (
    "context"
    "errors"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

func dialAPIServer(apiAddress string) (*grpc.ClientConn, context.Context, context.CancelFunc, error) {
    cred, err := credentials.NewClientTLSFromFile("./cert.pem", "")
    if err != nil {
        return nil, nil, nil, err
    }
    ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
    conn, err := grpc.DialContext(
        ctx,
        apiAddress,
        grpc.WithTransportCredentials(cred),
    )
    if err != nil {
        return nil, nil, cancelFunc, errors.New("failed to dial " + apiAddress)
    }
    return conn, ctx, cancelFunc, nil
}

package edge

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	pb "github.com/bionicosmos/aegle/edge/xray"
	"google.golang.org/grpc"
)

func Start() {
	listener, err := net.Listen("tcp", os.Getenv("LISTEN"))
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	pb.RegisterXrayServer(server, &pb.Server{})
	fmt.Println(listener.Addr())
	go func() {
		if err := server.Serve(listener); err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	server.GracefulStop()
}

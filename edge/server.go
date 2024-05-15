package edge

import (
	"fmt"
	"log"
	"net"

	pb "github.com/bionicosmos/aegle/edge/xray"
	"google.golang.org/grpc"
)

func Start() {
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	pb.RegisterXrayServer(server, &pb.Server{})
	fmt.Println(listener.Addr())
	log.Fatal(server.Serve(listener))
}

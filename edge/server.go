package edge

import (
	"fmt"
	"log"
	"net"

	"github.com/bionicosmos/aegle/config"
	pb "github.com/bionicosmos/aegle/edge/xray"
	"google.golang.org/grpc"
)

func Start() {
	listener, err := net.Listen("tcp", config.C.Listen)
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	pb.RegisterXrayServer(server, &pb.Server{})
	fmt.Println(listener.Addr())
	log.Fatal(server.Serve(listener))
}

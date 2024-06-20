package edge

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/bionicosmos/aegle/edge/xray"
	"google.golang.org/grpc"
)

func Start() {
	listener, err := net.Listen("tcp", os.Getenv("LISTEN"))
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	pb.RegisterXrayServer(server, &pb.Server{})
	fmt.Println(listener.Addr())
	log.Fatal(server.Serve(listener))
}

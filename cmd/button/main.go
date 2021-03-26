package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/coryb/poolpi/pb"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:8888", "The server address in the format of host:port")
)

func main() {
	flag.Parse()
	key := os.Args[1]
	pbKey, ok := pb.Key_value[key]
	if !ok {
		log.Fatalf("Invalid Key: %s", key)
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewPoolClient(conn)

	stream, err := client.Events(context.Background())
	if err != nil {
		log.Fatalf("%v.Events(_) = _, %v", client, err)
	}

	stream.Send(&pb.KeyEvent{
		Key: pb.Key(pbKey),
	})

	for {
		ev, err := stream.Recv()
		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}
		if state := ev.GetState(); state != nil {
			log.Printf("Active: %s", state.Summary())
			stream.CloseSend()
			return
		}
	}
}

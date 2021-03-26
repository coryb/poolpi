package main

import (
	"context"
	"flag"
	"io"
	"log"

	"github.com/coryb/poolpi/pb"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:8888", "The server address in the format of host:port")
)

func main() {
	flag.Parse()
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

	waitC := make(chan struct{})
	go func() {
		for {
			ev, err := stream.Recv()
			if err != nil {
				// read done.
				close(waitC)
				if err != io.EOF {
					log.Printf("ERROR: %s", err)
				}
				return
			}
			log.Printf("|<--- %s", ev.Summary())
		}
	}()
	// for _, note := range notes {
	// 	if err := stream.Send(note); err != nil {
	// 		log.Fatalf("Failed to send a note: %v", err)
	// 	}
	// }
	// stream.CloseSend()
	<-waitC
}

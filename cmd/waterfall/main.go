package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/coryb/poolpi/pb"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:8888", "The server address in the format of host:port")
	duration   = flag.String("duration", "20m", "How long to run the waterfall")
)

// The goal is to run the waterfall for a fixed period so that the water does not stagnate.
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

	offTime, err := time.ParseDuration(*duration)
	if err != nil {
		log.Fatalf("Failed to parse duration %q: %s", *duration, err)
	}
	end := time.After(offTime)

	stream.Send(&pb.KeyEvent{Key: pb.Key_Aux3})
	log.Printf("Waterfall started")

	// wait for the end time, meanwhile just print out the message events
	done := false
	for !done {
		select {
		case <-end:
			done = true
		default:
			ev, err := stream.Recv()
			if err != nil {
				log.Printf("ERROR: %s", err)
				done = true
			}
			if msg := ev.GetMessage(); msg != nil {
				log.Printf("Message: %s", msg.Plain())
			}
		}
	}
	stream.Send(&pb.KeyEvent{Key: pb.Key_Aux3})
	log.Printf("Waterfall stopped")
}

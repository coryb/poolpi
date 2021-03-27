package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/coryb/poolpi/events"
	"github.com/coryb/poolpi/pb"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:8888", "The server address in the format of host:port")
	duration   = flag.String("duration", "20m", "How long to run the waterfall")
)

func fatalErr(err error) {
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
}

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

	ctx := context.Background()
	client, err := events.NewClient(ctx, conn)
	fatalErr(err)

	offTime, err := time.ParseDuration(*duration)
	fatalErr(err)

	err = client.Key(pb.Key_Aux3)
	fatalErr(err)
	log.Printf("Waterfall started")

	// wait for the end time, meanwhile just print out the message events
	end, cancel := context.WithTimeout(ctx, offTime)
	defer cancel()
	err = client.Messages(end, func(m *pb.MessageEvent) {
		log.Printf("Message: %s", m.Plain())
	})

	err = client.Key(pb.Key_Aux3)
	fatalErr(err)
	log.Printf("Waterfall stopped")
}

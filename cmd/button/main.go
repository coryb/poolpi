package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/coryb/poolpi/events"
	"github.com/coryb/poolpi/pb"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:8888", "The server address in the format of host:port")
)

func fatalErr(err error) {
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
}

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

	ctx := context.Background()
	client, err := events.NewClient(ctx, conn)
	fatalErr(err)

	err = client.Key(pb.Key(pbKey))
	fatalErr(err)

	state, err := client.CurrentState()
	fatalErr(err)
	log.Print(state.Summary())
}

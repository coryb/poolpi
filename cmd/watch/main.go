package main

import (
	"context"
	"flag"
	"io"
	"log"

	"github.com/coryb/poolpi"
	"github.com/coryb/poolpi/pb"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:8888", "The server address in the format of host:port")
	raw        = flag.Bool("raw", false, "Show raw uninterpreted content")
	unknown    = flag.Bool("unknown", false, "Show only unknown events")
	binary     = flag.Bool("binary", false, "show raw events in binary instead of hex")
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

	format := poolpi.FormatHex
	if *binary {
		format = poolpi.FormatBinary
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
			if _, ok := ev.Event.(*pb.Event_Unknown); *unknown && !ok {
				continue
			}
			if *raw {
				log.Printf("<--- %s", poolpi.EventFromPB(ev).Format(format))
			} else {
				log.Printf("|<--- %s", ev.Summary())
			}
		}
	}()
	<-waitC
}

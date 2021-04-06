package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/coryb/poolpi"
	"github.com/coryb/poolpi/pb"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8888, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	sys, err := poolpi.NewSystem("/dev/ttyS0")
	if err != nil {
		log.Fatal(err)
	}

	unsubscribe := sys.Subscribe(func(e poolpi.Event) {
		log.Printf("|<--- %s", e.Summary())
	})

	defer unsubscribe()

	sys.Start()
	defer sys.Close()

	pb.RegisterPoolServer(grpcServer, newServer(sys))
	grpcServer.Serve(lis)
}

type poolServer struct {
	pool *poolpi.System
	pb.UnimplementedPoolServer
}

func newServer(s *poolpi.System) pb.PoolServer {
	return &poolServer{
		pool:                    s,
		UnimplementedPoolServer: pb.UnimplementedPoolServer{},
	}
}

func (s *poolServer) Events(stream pb.Pool_EventsServer) error {
	eg, ctx := errgroup.WithContext(stream.Context())
	eg.Go(func() error {
		for {
			pbKey, err := stream.Recv()
			if err != nil {
				return err
			}
			key := poolpi.KeyFromPB(pbKey.Key)
			s.pool.Send(key.ToEvent(poolpi.EventRemoteKey))
			log.Printf("|---> Key %s", key.String())
		}
	})

	eg.Go(func() error {
		errC := make(chan error)
		defer close(errC)
		unsubscribe := s.pool.Subscribe(func(e poolpi.Event) {
			err := stream.Send(e.ToPB())
			if err != nil {
				errC <- err
			}
		})
		defer unsubscribe()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e := <-errC:
			return e
		}
	})

	err := eg.Wait()
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

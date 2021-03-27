package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/coryb/poolpi/pb"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:8888", "The server address in the format of host:port")
)

// This will assert the pool/spa is in a good state
// Set Spa Temp ???
// Set Pool Temp ???
// Set Heater Manual Off
// Set Spa Pump Speed 90%

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

	assertState(stream)
	setSpaSpeed(stream, "90%")
}

func setSpaSpeed(stream pb.Pool_EventsClient, percent string) {
	menu(stream, "Settings Menu")
	keyUntil(stream, pb.Key_Right, "Spa Heater1")
	keyUntil(stream, pb.Key_Right, "Pool Heater1")
	keyUntil(stream, pb.Key_Right, "VSP Speed Settings")
	keyUntil(stream, pb.Key_Plus, "Filter Speed1")
	keyUntil(stream, pb.Key_Right, "Filter Speed2")
	keyUntil(stream, pb.Key_Right, "Filter Speed3")
	keyUntil(stream, pb.Key_Right, "Filter Speed4")
	msg := keyUntil(stream, pb.Key_Right, "Spa Speed")
	for !strings.Contains(msg.Plain(), percent) {
		msg = keyUntil(stream, pb.Key_Plus, "Spa Speed")
	}
	log.Printf("Asserted Spa Speed: %s", percent)
}

func menu(stream pb.Pool_EventsClient, name string) {
	for {
		msg := keyUntil(stream, pb.Key_Menu, "Menu")
		if strings.Contains(msg.Plain(), name) {
			return
		}
	}
}

func keyUntil(stream pb.Pool_EventsClient, key pb.Key, expected string) (message *pb.MessageEvent) {
	done := make(chan struct{})
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(done)
		for {
			ev, err := stream.Recv()
			if err != nil {
				log.Fatalf("ERROR: %s", err)
			}
			if msg := ev.GetMessage(); msg != nil {
				log.Printf("Message: %s", msg.Plain())
				if strings.Contains(msg.Plain(), expected) {
					message = msg
					return
				}
			}
		}
	}()

	for {
		stream.Send(&pb.KeyEvent{Key: key})
		timeout := time.After(1000 * time.Millisecond)
		select {
		case <-done:
			return
		case <-timeout:
		}
	}
}

func assertState(stream pb.Pool_EventsClient) {
	for {
		ev, err := stream.Recv()
		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}
		if state := ev.GetState(); state != nil {
			for name, ind := range state.Indicators() {
				switch name {
				case pb.Pool:
					if !ind.GetActive() {
						stream.Send(&pb.KeyEvent{Key: pb.Key_PoolSpa})
						log.Printf("Asserted Pool Mode")
						assertState(stream)
						return
					}
				case pb.Filter:
					// no op, this is managed by the Prologic automation
				default:
					toggle := map[pb.IndicatorName]pb.Key{
						pb.Heater1: pb.Key_Heater,
						pb.Valve3:  pb.Key_Valve3,
						pb.Lights:  pb.Key_Lights,
						pb.Aux1:    pb.Key_Aux1,
						pb.Aux2:    pb.Key_Aux2,
						pb.Aux3:    pb.Key_Aux3,
						pb.Aux4:    pb.Key_Aux4,
						pb.Aux5:    pb.Key_Aux5,
						pb.Aux6:    pb.Key_Aux6,
						pb.Valve4:  pb.Key_Valve4,
						pb.Aux7:    pb.Key_Aux7,
					}
					if key, ok := toggle[name]; ok && ind.GetActive() {
						stream.Send(&pb.KeyEvent{Key: key})
						log.Printf("Asserted %s Off", name)
						assertState(stream)
						return
					}
				}
			}
			return
		}
	}
}

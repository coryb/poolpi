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
	duration   = flag.String("duration", "30m", "How long to filter the spa")
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

	current := currentState(stream)
	if current.GetHeater1().GetActive() {
		log.Printf("Heater on, skipping Spa filtering")
		return
	}

	offTime, err := time.ParseDuration(*duration)
	if err != nil {
		log.Fatalf("Failed to parse duration %q: %s", *duration, err)
	}
	end := time.After(offTime)

	// ensure Heater is disabled
	menu(stream, "Default Menu")
	keyUntil(stream, pb.Key_Right, "day")
	keyUntil(stream, pb.Key_Right, "Air Temp")
	msg := keyUntil(stream, pb.Key_Right, "Heater1")
	for !strings.Contains(msg.Plain(), "Manual Off") {
		msg = keyUntil(stream, pb.Key_Heater, "Heater1")
	}

	setSpaSpeed(stream, pb.Key_Minus, "50%")

	stream.Send(&pb.KeyEvent{Key: pb.Key_PoolSpa})
	log.Printf("Set Spa Mode")

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
	stream.Send(&pb.KeyEvent{Key: pb.Key_PoolSpa})
	log.Printf("Set Pool Mode")
	setSpaSpeed(stream, pb.Key_Plus, "90%")

	// ensure Heater is back in Auto mode3
	menu(stream, "Default Menu")
	keyUntil(stream, pb.Key_Right, "day")
	keyUntil(stream, pb.Key_Right, "Air Temp")
	msg = keyUntil(stream, pb.Key_Right, "Heater1")
	for !strings.Contains(msg.Plain(), "Auto") {
		msg = keyUntil(stream, pb.Key_Heater, "Heater1")
	}
}

func setSpaSpeed(stream pb.Pool_EventsClient, direction pb.Key, percent string) {
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
		msg = keyUntil(stream, direction, "Spa Speed")
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

func currentState(stream pb.Pool_EventsClient) *pb.StateEvent {
	for {
		ev, err := stream.Recv()
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
		if state := ev.GetState(); state != nil {
			return state
		}
	}
}

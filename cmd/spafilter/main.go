package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/coryb/poolpi/events"
	"github.com/coryb/poolpi/pb"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:8888", "The server address in the format of host:port")
	duration   = flag.String("duration", "30m", "How long to filter the spa")
)

func fatalErr(err error) {
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	fatalErr(err)
	defer conn.Close()

	ctx := context.Background()
	client, err := events.NewClient(ctx, conn)
	fatalErr(err)

	current, err := client.CurrentState()
	fatalErr(err)
	currentMap := current.Indicators()
	for _, name := range []pb.IndicatorName{pb.Heater1, pb.Spa, pb.SystemOff, pb.Service} {
		if ind, ok := currentMap[name]; ok && ind.GetActive() {
			log.Printf("%s active, skipping Spa filtering", name)
			return
		}
	}

	offTime, err := time.ParseDuration(*duration)
	fatalErr(err)

	// ensure Heater is disabled
	err = client.KeyMenu(ctx, "Default Menu")
	fatalErr(err)
	msg, err := client.KeySequence(ctx,
		[]events.KeyPrompt{
			{Key: pb.Key_Right, Expected: "day"},
			{Key: pb.Key_Right, Expected: "Heater1"},
			{Key: pb.Key_Right, Expected: "Air Temp"},
			{Key: pb.Key_Right, Expected: "Heater1"},
		}...,
	)
	fatalErr(err)
	for !strings.Contains(msg.Plain(), "Manual Off") {
		msg, err = client.KeyUntil(ctx, pb.Key_Heater, "Heater1")
		fatalErr(err)
	}

	setSpaSpeed(ctx, client, pb.Key_Minus, "50%")

	err = client.Key(pb.Key_PoolSpa)
	fatalErr(err)
	log.Printf("Set Spa Mode")

	// wait for the end time, meanwhile just print out the message events
	end, cancel := context.WithTimeout(ctx, offTime)
	defer cancel()
	err = client.Messages(end, func(m *pb.MessageEvent) {
		log.Printf("Message: %s", m.Plain())
	})

	err = client.Key(pb.Key_PoolSpa)
	fatalErr(err)
	log.Printf("Set Pool Mode")

	setSpaSpeed(ctx, client, pb.Key_Plus, "90%")

	// ensure Heater is back in Auto mode3
	err = client.KeyMenu(ctx, "Default Menu")
	fatalErr(err)

	msg, err = client.KeySequence(ctx, []events.KeyPrompt{
		{Key: pb.Key_Right, Expected: "day"},
		{Key: pb.Key_Right, Expected: "Air Temp"},
		{Key: pb.Key_Right, Expected: "Heater1"},
	}...)
	for !strings.Contains(msg.Plain(), "Auto") {
		msg, err = client.KeyUntil(ctx, pb.Key_Heater, "Heater1")
		fatalErr(err)
	}
}

func setSpaSpeed(ctx context.Context, c *events.Client, direction pb.Key, percent string) {
	err := c.KeyMenu(ctx, "Settings Menu")
	fatalErr(err)
	msg, err := c.KeySequence(ctx, []events.KeyPrompt{
		{Key: pb.Key_Right, Expected: "Spa Heater1"},
		{Key: pb.Key_Right, Expected: "Pool Heater1"},
		{Key: pb.Key_Right, Expected: "VSP Speed Settings"},
		{Key: pb.Key_Plus, Expected: "Filter Speed1"},
		{Key: pb.Key_Right, Expected: "Filter Speed2"},
		{Key: pb.Key_Right, Expected: "Filter Speed3"},
		{Key: pb.Key_Right, Expected: "Filter Speed4"},
		{Key: pb.Key_Right, Expected: "Spa Speed"},
	}...)
	for !strings.Contains(msg.Plain(), percent) {
		msg, err = c.KeyUntil(ctx, direction, "Spa Speed")
		fatalErr(err)
	}
	log.Printf("Asserted Spa Speed: %s", percent)
}

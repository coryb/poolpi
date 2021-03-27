package main

import (
	"context"
	"flag"
	"log"
	"strings"

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

// This will assert the pool/spa is in a good state
// TODO Set Spa Temp ???
// TODO Set Pool Temp ???
// TODO Set Clock ???

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

	assertState(client)
	setSpaSpeed(ctx, client, "90%")
}

func setSpaSpeed(ctx context.Context, c *events.Client, percent string) {
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
		msg, err = c.KeyUntil(ctx, pb.Key_Plus, "Spa Speed")
		fatalErr(err)
	}
	log.Printf("Asserted Spa Speed: %s", percent)
}

func assertState(c *events.Client) {
	for {
		state, err := c.CurrentState()
		fatalErr(err)
		for name, ind := range state.Indicators() {
			switch name {
			case pb.Pool:
				if !ind.GetActive() {
					err = c.Key(pb.Key_PoolSpa)
					fatalErr(err)
					log.Printf("Asserted Pool Mode")
					assertState(c)
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
					err = c.Key(key)
					fatalErr(err)
					log.Printf("Asserted %s Off", name)
					assertState(c)
					return
				}
			}
			return
		}
	}
}

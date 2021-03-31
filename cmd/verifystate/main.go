package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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
	setClock(ctx, client)
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

func setClock(ctx context.Context, c *events.Client) {
	n := time.Now().Local()
	if n.Hour() > 22 || n.Hour() < 4 {
		log.Printf("Skipping time check to avoid any timezone clock adjustments")
		return
	}

	err := c.KeyMenu(ctx, "Default Menu")
	fatalErr(err)
	msg, err := c.KeyUntil(ctx, pb.Key_Right, "day")
	fatalErr(err)

	// display only shows "A" not "AM"
	display := msg.Plain() + "M"

	parsed, err := time.ParseInLocation("Monday: 3:04PM", display, n.Location())
	fatalErr(err)

	// Valid times are "now +/- 1m"
	for _, offset := range []time.Duration{-time.Minute, 0, time.Minute} {
		nowish := n.Add(offset).Format("Monday: 3:04PM")
		if display == nowish {
			log.Printf("Time Validated: %s", nowish)
			return
		}
	}

	err = c.KeyMenu(ctx, "Settings Menu")
	fatalErr(err)
	msg, err = c.KeySequence(ctx, []events.KeyPrompt{
		{Key: pb.Key_Right, Expected: "Spa Heater1"},
		{Key: pb.Key_Right, Expected: "Pool Heater1"},
		{Key: pb.Key_Right, Expected: "VSP Speed Settings"},
		{Key: pb.Key_Right, Expected: "Set Day and Time"},
	}...)
	fatalErr(err)

	line2 := msg.Message[len(msg.Message)/2:]
	line2 = bytes.TrimSpace(line2)

	weekday := n.Weekday()

	if parsed.Weekday() != weekday {
		log.Printf("Clock wrong, weekday %s expected %s", parsed.Weekday(), n.Weekday())
		for !isFocused(FocusDay, msg) {
			msg, err = c.KeyUntil(ctx, pb.Key_Right, "Set Day and Time")
			fatalErr(err)
		}
		key := pb.Key_Plus
		if parsed.Weekday() > weekday {
			key = pb.Key_Minus
		}
		for !strings.Contains(msg.Plain(), weekday.String()) {
			msg, err = c.KeyUntil(ctx, key, "Set Day and Time")
			fatalErr(err)
		}
	}

	hour := n.Hour()
	if parsed.Hour() != hour {
		suffix := "A"
		if hour > 11 {
			suffix = "P"
		}
		h := hour % 12
		if h == 0 {
			h = 12
		}
		log.Printf("Setting Hour to %d%s", h, suffix)
		for !isFocused(FocusHour, msg) {
			msg, err = c.KeyUntil(ctx, pb.Key_Right, "Set Day and Time")
			fatalErr(err)
		}
		key := pb.Key_Plus
		if parsed.Hour()%12 > hour%12 {
			key = pb.Key_Minus
		}
		for !isHour(h, suffix, msg) {
			msg, err = c.KeyUntil(ctx, key, "Set Day and Time")
			fatalErr(err)
		}
	}

	minute := n.Minute()
	if parsed.Minute() != minute {
		for !isFocused(FocusMinute, msg) {
			msg, err = c.KeyUntil(ctx, pb.Key_Right, "Set Day and Time")
			fatalErr(err)
		}
		key := pb.Key_Plus
		if parsed.Minute() > minute {
			key = pb.Key_Minus
		}
		strMinute := fmt.Sprintf(":%02d", minute)
		log.Printf("Setting Minute to %s", strMinute)
		for !strings.Contains(msg.Plain(), strMinute) {
			msg, err = c.KeyUntil(ctx, key, "Set Day and Time")
			fatalErr(err)
		}
	}
}

func isHour(h int, suffix string, msg *pb.MessageEvent) bool {
	plain := msg.Plain()
	if !strings.HasPrefix(plain, "Set Day and Time") {
		return false
	}
	if strings.Contains(plain, " "+strconv.Itoa(h)+":") && strings.HasSuffix(plain, suffix) {
		return true
	}
	return false
}

type Focus int8

const (
	FocusDay    Focus = iota
	FocusHour         = iota
	FocusMinute       = iota
)

func isFocused(f Focus, msg *pb.MessageEvent) bool {
	if !strings.HasPrefix(msg.Plain(), "Set Day and Time") {
		return false
	}
	line2 := bytes.TrimSpace(msg.Message[len(msg.Message)/2:])
	words := bytes.Fields(line2)
	switch f {
	case FocusDay:
		focused := words[0][0]&0x80 > 0
		log.Printf("Day Focused: %t", focused)
		return focused
	case FocusHour:
		focused := words[1][0]&0x80 > 0
		log.Printf("Hour Focused: %t", focused)
		return focused
	case FocusMinute:
		split := bytes.Split(words[1], []byte{':'})
		focused := split[1][0]&0x80 > 0
		log.Printf("Minute Focused: %t", focused)
		return focused
	}
	return false
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

package events

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/coryb/poolpi/pb"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Client struct {
	stream pb.Pool_EventsClient
}

func NewClient(ctx context.Context, conn *grpc.ClientConn) (*Client, error) {
	client := pb.NewPoolClient(conn)
	stream, err := client.Events(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to establish events stream: %w", err)
	}
	return &Client{
		stream: stream,
	}, nil
}

func (c *Client) CurrentState() (*pb.StateEvent, error) {
	for {
		ev, err := c.stream.Recv()
		if err != nil {
			return nil, fmt.Errorf("failed to read from event stream: %w", err)
		}
		if state := ev.GetState(); state != nil {
			return state, nil
		}
	}
}

func (c *Client) Key(key pb.Key) error {
	return c.stream.Send(&pb.KeyEvent{Key: key})
}

func (c *Client) KeyMenu(ctx context.Context, name string) error {
	for {
		msg, err := c.KeyUntil(ctx, pb.Key_Menu, "Menu")
		if err != nil {
			return err
		}
		if strings.Contains(msg.Plain(), name) {
			return nil
		}
	}
}

func (c *Client) KeyUntil(ctx context.Context, key pb.Key, expected string) (message *pb.MessageEvent, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	defer func() {
		egErr := eg.Wait()
		if err == nil {
			err = egErr
		}
	}()

	eg.Go(func() error {
		for {
			ev, err := c.stream.Recv()
			if err != nil {
				return fmt.Errorf("failed to read from event stream: %w", err)
			}
			if msg := ev.GetMessageUpdate(); msg != nil {
				log.Printf("Message: %s", msg.Plain())
				if strings.Contains(msg.Plain(), expected) {
					message = msg.GetMessage()
					cancel()
					return nil
				}
			}
		}
	})

	count := 0
	for {
		count++
		if count > 10 {
			cancel()
			return nil, fmt.Errorf(
				"After sending %s %d times, expected message %q not found",
				pb.Key_name[int32(key)], count, expected,
			)
		}
		err = c.Key(key)
		if err != nil {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(1000 * time.Millisecond):
		}
	}
}

type KeyPrompt struct {
	Key      pb.Key
	Expected string
}

func (c *Client) KeySequence(ctx context.Context, seq ...KeyPrompt) (message *pb.MessageEvent, err error) {
	for _, keyPrompt := range seq {
		message, err = c.KeyUntil(ctx, keyPrompt.Key, keyPrompt.Expected)
		if err != nil {
			return nil, err
		}
	}
	return
}

func (c *Client) Messages(ctx context.Context, f func(*pb.MessageEvent)) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			ev, err := c.stream.Recv()
			if err != nil {
				return err
			}
			if msg := ev.GetMessage(); msg != nil {
				f(msg)
			}
		}
	}
}

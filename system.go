package poolpi

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/tarm/serial"
	"golang.org/x/sync/errgroup"
)

const (
	FrameDLE = 0x10 // "data link escape"
	FrameESC = 0x00 // escape
	FrameSTX = 0x02 // start
	FrameETX = 0x03 // end
)

type System struct {
	serialPort  *serial.Port
	outgoing    chan []byte
	subscribers sync.Map
	eg          *errgroup.Group
}

func NewSystem(serialFile string) (*System, error) {
	c := &serial.Config{Name: serialFile, Baud: 19200, StopBits: serial.Stop2}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}

	return &System{
		serialPort: s,
		outgoing:   make(chan []byte),
	}, nil
}

func (s *System) Start() {
	s.eg = &errgroup.Group{}

	s.eg.Go(func() error {
		data := []byte{}
		collectData := false
		buf := make([]byte, 1)
		for {
			if !collectData {
				data = data[:0]
			}
			err := s.read(buf)
			if err != nil {
				return err
			}
			switch b := buf[0]; b {
			case FrameDLE:
				err = s.read(buf)
				if err != nil {
					return err
				}
				switch b := buf[0]; b {
				case FrameESC:
					// escape sequence 0x10 0x00 => 0x10
					data = append(data, 0x10)
				case FrameSTX:
					// start
					collectData = true
					continue
				case FrameETX:
					// end
					collectData = false
					if len(data) < 2 {
						continue
					}
					crc := NewCRC(data[len(data)-2:])
					data = data[:len(data)-2]
					if chk := ComputeCRC(data); crc != chk {
						log.Printf("WARNING: Invalid CRC %d/%d: %v", crc, chk, data)
						continue
					}
					s.handleEvent(NewEvent(data))
				default:
					log.Printf("Unknown Sequence: 0x10%x", b)
					collectData = false
				}
			default:
				if collectData {
					data = append(data, b)
				}
			}
		}
	})
}

func (s *System) read(buf []byte) error {
	for {
		n, err := s.serialPort.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if n <= 0 {
			// got no bytes, keep going
			continue
		}
		return nil
	}
}

func (s *System) Close() error {
	if s.serialPort != nil {
		err := s.serialPort.Flush()
		if err != nil {
			return fmt.Errorf("failed to flush serial port connection: %w", err)
		}
		err = s.serialPort.Close()
		if err != nil {
			return fmt.Errorf("failed to close serial port connection: %w", err)
		}
	}
	return s.eg.Wait()
}

func (s *System) Send(e Event) {
	s.outgoing <- e.ToBytes()
}

func (s *System) Subscribe(f func(e Event)) (unsubscribe func()) {
	id := rand.Uint64()
	s.subscribers.Store(id, f)
	return func() {
		s.subscribers.Delete(id)
	}
}

func (s *System) handleEvent(e Event) {
	if e.Type == EventReady {
		select {
		case m := <-s.outgoing:
			// this sleep is a fuzz, if we respond "too quickly" it seems
			// the automation system is not ready to read key presses
			time.Sleep(2 * time.Millisecond)
			n, err := s.serialPort.Write(m)
			if err != nil {
				log.Printf("ERROR: %s", err)
				return
			}
			if n != len(m) {
				log.Printf("WARNING: Wrote %d/%d", n, len(m))
			}
		default:
		}
		return
	}
	s.subscribers.Range(func(_, value interface{}) bool {
		if f, ok := value.(func(Event)); ok {
			f(e)
		}
		return true
	})
}

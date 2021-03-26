package main

import (
	"context"
	"log"
	"os"

	"github.com/coryb/poolpi"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	sys, err := poolpi.NewSystem("/dev/ttyS0")
	if err != nil {
		log.Fatal(err)
	}

	unsubscribe := sys.Subscribe(func(e poolpi.Event) {
		log.Printf("|<-- %s", e.Summary())
	})
	defer unsubscribe()

	sys.Start()
	<-context.Background().Done()
}

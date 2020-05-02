package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/teivah/tourniquet"
	"google.golang.org/grpc"
)

func main() {
	t, err := tourniquet.NewTourniquet(newConnection, 3, time.Second)
	if err != nil {
		log.Fatal(err)
	}

	// Spin up 10 consumer goroutines that will race for an available connection
	nbConsumers := 10
	// The wait group is just here for the sake of the example to wait for all the consumers before
	// to stop the example
	wg := sync.WaitGroup{}
	wg.Add(nbConsumers)
	for i := 0; i < nbConsumers; i++ {
		i := i
		go func() {
			// We allow passing a context while retrieving a connection
			// In this example we pass a timeout that will represent the maximum time we want to wait
			// before to get an available connection from the pool
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			conn, err := t.Get(ctx)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("I'm consumer %d and I received connection %v\n", i, conn)

			// Once we used the connection, we put it back to the pool
			t.Free(conn)

			// If we get an error while using the connection, we may want to recreate one
			// conn.ClientConn.Close()
			// conn, err = t.Recreate()

			wg.Done()
		}()
	}
	wg.Wait()
}

func newConnection() (*grpc.ClientConn, error) {
	// Create a connection using grpc.Dial
	// In this example we create a dummy connection
	return &grpc.ClientConn{}, nil
}

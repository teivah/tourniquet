package tourniquet

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// Pool is the main pool structure.
type Pool struct {
	connFactory func() (*grpc.ClientConn, error)
	pool        chan Connection
	ttl         time.Duration
}

// Connection wraps a gRPC connection.
type Connection struct {
	ClientConn *grpc.ClientConn
	t          time.Time
}

// NewPool creates a pool of gRPC connections.
func NewPool(connFactory func() (*grpc.ClientConn, error), desiredPoolSize int, ttl time.Duration) (*Pool, error) {
	pool := make(chan Connection, desiredPoolSize)
	for i := 0; i < desiredPoolSize; i++ {
		conn, err := connFactory()
		if err != nil {
			return nil, err
		}
		pool <- Connection{
			ClientConn: conn,
			t:          time.Now(),
		}
	}

	return &Pool{
		connFactory: connFactory,
		pool:        pool,
		ttl:         ttl,
	}, nil
}

// Get returns a connection from the pool or recreates one.
func (t *Pool) Get(ctx context.Context) (Connection, error) {
	select {
	case <-ctx.Done():
		return Connection{}, ctx.Err()
	case conn := <-t.pool:
		if time.Since(conn.t) > t.ttl {
			return t.Recreate()
		}
		return conn, nil
	}
}

// Free frees a connection.
func (t *Pool) Free(conn Connection) {
	t.pool <- conn
}

// Recreate recreates a connection.
func (t *Pool) Recreate() (Connection, error) {
	conn, err := t.connFactory()
	if err != nil {
		return Connection{}, err
	}
	return Connection{
		ClientConn: conn,
		t:          time.Now(),
	}, err
}

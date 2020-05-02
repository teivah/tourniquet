package tourniquet

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type Tourniquet struct {
	connFactory func() (*grpc.ClientConn, error)
	pool        chan Connection
	ttl         time.Duration
}

type Connection struct {
	ClientConn *grpc.ClientConn
	t          time.Time
}

func NewTourniquet(connFactory func() (*grpc.ClientConn, error), desiredPoolSize int, ttl time.Duration) (*Tourniquet, error) {
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

	return &Tourniquet{
		connFactory: connFactory,
		pool:        pool,
		ttl:         ttl,
	}, nil
}

func (t *Tourniquet) Get(ctx context.Context) (Connection, error) {
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

func (t *Tourniquet) Free(conn Connection) {
	t.pool <- conn
}

func (t *Tourniquet) Recreate() (Connection, error) {
	conn, err := t.connFactory()
	if err != nil {
		return Connection{}, err
	}
	return Connection{
		ClientConn: conn,
		t:          time.Now(),
	}, err
}

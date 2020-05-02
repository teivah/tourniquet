package tourniquet

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func Test_Tourniquet(t *testing.T) {
	size := 3
	tourniquet, err := NewPool(func() (*grpc.ClientConn, error) {
		return &grpc.ClientConn{}, nil
	}, size, time.Second)
	require.NoError(t, err)

	local := make([]Connection, 0)
	for i := 0; i < size; i++ {
		conn, err := tourniquet.Get(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, conn.ClientConn)
		local = append(local, conn)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	_, err = tourniquet.Get(ctx)
	require.Error(t, err)

	tourniquet.Free(local[0])
	conn, err := tourniquet.Get(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, conn.ClientConn)
}

func Test_Tourniquet_Recreate(t *testing.T) {
	tourniquet, err := NewPool(func() (*grpc.ClientConn, error) {
		return &grpc.ClientConn{}, nil
	}, 1, time.Nanosecond)
	require.NoError(t, err)
	time.Sleep(time.Millisecond)
	conn, err := tourniquet.Get(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, conn.ClientConn)
}

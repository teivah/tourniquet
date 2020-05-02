# Tourniquet

Tourniquet manages a pool of gRPC connections to handle gRPC client-side load-balancing in Go.

It does not rely on DNS as it may cause many problems in Kubernetes. Especially in the event of scaling the number of replicas.

The principle is to set a desired pool size a TTL for each connection. Once the TTL is reached, it will force to recreate a connection.

It allows handling gRPC load-balancing and discoverability.

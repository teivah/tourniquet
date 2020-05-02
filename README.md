# Tourniquet

Tourniquet manages a thread-safe pool of gRPC connections to handle gRPC client-side load-balancing in Go.

It does not rely on DNS as it may cause many problems in Kubernetes. Especially in the event of scaling the number of replicas.

The principle is to set a desired pool size a TTL for each connection. Once the TTL is reached, it will force to recreate a connection.

It allows handling gRPC load-balancing and discoverability [without tears and without a service mesh](https://kubernetes.io/blog/2018/11/07/grpc-load-balancing-on-kubernetes-without-tears/).
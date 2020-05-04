# Tourniquet

gRPC load-balancing on Kubernetes can be painful: https://kubernetes.io/blog/2018/11/07/grpc-load-balancing-on-kubernetes-without-tears/

Instead of having to rely on a service-mesh, a cheap alternative is to manage **client-side load balancing**. 

In Kubernetes, we have to configure a Service on top of a Pod. Then, with Tourniquet we can create a pool of gRPC connections and a max TTL. Each connection created is associated to the TTL specified. Once the TTL is reached, the stale connection will be closed and recreated. It allows querying again the service and potentially _discover_ new replica instances.   

Tourniquet is a cheap and race-safe solution to handle gRPC load balancing on Kubernetes.

More info:
* [Documentation](https://pkg.go.dev/github.com/teivah/tourniquet)
* [Example](examples/examples.go)
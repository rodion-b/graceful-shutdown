

//Requirment:
//Implement a TCP Server with Graceful Shutdown

//Requirements:
// Implement a graceful shutdown mechanism, allowing active requests to complete before shutting down.
// This is to support running the service in a cloud environment where it will have a short amount of time to complete in-flight requests before the server is terminated.
// The server should stop accepting new connections, but can continue accepting requests.
// The allowed grace period for active requests to complete should be configurable, for example 3 seconds.
// Requests that have been accepted, but not completed after that grace period - should be rejected with: RESPONSE|REJECTED|Cancelled
// Requests that have not been accepted, can be discarded without a response. (ex. slow clients)
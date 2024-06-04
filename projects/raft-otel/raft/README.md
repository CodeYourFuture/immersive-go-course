# Sample Raft implementation

This is based on Eli Bendersky's [https://eli.thegreenplace.net] RAFT demo code.

I've modified it in a few ways:
 * Adds a main.go so you can run the RAFT code as docker containers (or Kube) - peers are found via DNS lookup
 * Changed from integer based peer IDs to use of IP addresses, so that the instances can come up without coordination with each other
 * Changed from standard RPC to gRPC (as we've been using throughout this course)
 * Adds Dockerfile and docker-compose.yml
 * Added structure to Command (to simplify gRPCing)
 * Removed one test (TestCrashAfterSubmit) as could not make the timing work to reliably crash leader before it had a chance to commit a change (would be easier to do this if code were restructured to inject time)
 * Added endpoint for doing some sets/gets of data, and a simple client that calls this - to demo what's usually done with RAFT, also added a client that exercises it

 ## Building and running this project

 If you change the raft.proto protocol buffer definitions, you must regenerate the bindings by:

 ```
protoc --proto_path=. --go_out=. --go-grpc_out=. raft.proto
 ```

 To run this under docker-compose, use `docker-compose up --build -d` or your preferred variant. 
 Docker must be installed and running.
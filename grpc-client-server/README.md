# GRPC Client and Server Communication 

The goals of this project are: 
 * to understand the differences between gRPC and REST
 * to get experience using deadlines with gRPCs
 * to build a simple prober

Start by building a simple gRPC ‘hello world’ server and cli client in golang. See [gRPC Quickstart](https://grpc.io/docs/languages/go/quickstart/) for detailed instructions.

Now, modify your client and server. Rename from greeter_client and greeter_server to prober_client and prober_server.   

Add a client timeout after 0.5 second - https://pkg.go.dev/context#WithTimeout
Check the timeout 
if ctx.Err() == context.Canceled {
	return status.New(codes.Canceled, "Client cancelled, deadline exceeded.")
}



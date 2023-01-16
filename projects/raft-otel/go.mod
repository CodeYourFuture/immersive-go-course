module github.com/Jille/raft-grpc-example

go 1.13

require (
	github.com/Jille/grpc-multi-resolver v1.1.0
	github.com/Jille/raft-grpc-leader-rpc v1.1.0
	github.com/Jille/raft-grpc-transport v1.3.0
	github.com/Jille/raftadmin v1.2.0
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-hclog v1.3.1 // indirect
	github.com/hashicorp/raft v1.3.11
	github.com/hashicorp/raft-boltdb v0.0.0-20220329195025-15018e9b97e0
	github.com/honeycombio/honeycomb-opentelemetry-go v0.4.0
	github.com/honeycombio/opentelemetry-go-contrib/launcher v0.0.0-20230104152713-a01f612b1b01
	github.com/mattn/go-colorable v0.1.13 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.37.0
	go.opentelemetry.io/otel v1.11.2
	go.opentelemetry.io/otel/trace v1.11.2
	golang.org/x/net v0.2.0 // indirect
	google.golang.org/genproto v0.0.0-20221114212237-e4508ebdbee1 // indirect
	google.golang.org/grpc v1.51.0
	google.golang.org/protobuf v1.28.1
	moul.io/number-to-words v0.6.0
)

replace github.com/hashicorp/raft => /Users/margaritaglushkova/Desktop/GitHub/CYF_WORK/raft

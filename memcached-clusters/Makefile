replicated_up:
	docker-compose -f replicated/docker-compose.yml up -d

replicated_down:
	docker-compose -f replicated/docker-compose.yml down

sharded_up:
	docker-compose -f sharded/docker-compose.yml up -d

sharded_down:
	docker-compose -f sharded/docker-compose.yml down

test:
	go test ./... -v -cover
.PHONY: run
run:
	docker-compose up -d

destroy:
	docker-compose down
	docker rmi $(shell docker images | grep 'kafka-cron')

test:
	go test -v ./...

# publisher:
# 	go run ./producer/ --brokerList=localhost:9092 --topicPrefix=jobs --cronPath=cron.yaml

# consumer-a:
# 	go run ./consumer/ --brokerList=localhost:9092 --topic=jobs-cluster-a

# consumer-b:
# 	go run ./consumer/ --brokerList=localhost:9092 --topic=jobs-cluster-b
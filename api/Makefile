test:
	docker compose -f docker-compose-test.yaml up -d

	sleep 1

	go test -v ./... | tc

	docker compose -f docker-compose-test.yaml down


test-raw:
	go test -v ./...
# test a whole astragalaxy api
test-all:
	docker compose -f docker-compose-test.yaml up -d
	sleep 10
	go test -v ./... -json | tparse --all || true
	docker compose -f docker-compose-test.yaml down

test TEST:
    docker compose -f docker-compose-test.yaml up -d
    sleep 10
    go test -v ./... -json -run Test{{TEST}} | tparse --all || true
    docker compose -f docker-compose-test.yaml down


# test without pretty print
test-raw:
	go test -v ./...

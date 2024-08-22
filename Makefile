GEN=./internal/gen
PROTOPATH=./api/proto
SSOPROTOS=

protoc:
	protoc --go_out=${GEN} --go_opt=paths=source_relative \
    --go-grpc_out=${GEN} --go-grpc_opt=paths=source_relative \
	--proto_path=${PROTOPATH} ${SSOPROTOS}

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

unit_tests:
	go test -timeout 5m -v ./internal/...

int_tests:
	go test -timeout 10m -v ./tests/...

tests: unit_tests int_tests